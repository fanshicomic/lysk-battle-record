package main

import (
	"context"
	"sort"
	"strconv"
	"time"

	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

const (
	spreadsheetID = "1-ORnXBnav4JVtP673Oio5sNdVpk0taUSzG3kWZqhIuY" // 替换为你的 Google Sheets 文件ID
	sheetName     = "面板"
)

var srv *sheets.Service

func initGoogleSheets() {
	ctx := context.Background()

	b, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("无法读取凭证文件: %v", err)
	}

	config, err := google.JWTConfigFromJSON(b, sheets.SpreadsheetsScope)
	if err != nil {
		log.Fatalf("解析凭证失败: %v", err)
	}

	client := config.Client(ctx)
	srv, err = sheets.New(client)
	if err != nil {
		log.Fatalf("创建 Sheets 客户端失败: %v", err)
	}
}

func main() {
	initGoogleSheets()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // 可改为指定域名，如 "https://yourdomain.com"
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           1 * time.Hour,
	}))

	r.GET("/ping", ping)

	r.POST("/record", processRecord)
	r.GET("/records", getRecords)
	r.GET("/last-records", getLastRecords)
	r.GET("/record-count", getRecordCount)

	r.Run(":8080")
}

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

func processRecord(c *gin.Context) {
	var input map[string]interface{}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求格式错误", "detail": err.Error()})
		return
	}

	// 将所有字段转为 string
	record := map[string]string{}
	for k, v := range input {
		if s, ok := v.(string); ok {
			record[k] = s
		} else {
			record[k] = fmt.Sprintf("%v", v)
		}
	}

	// 按固定顺序写入
	keys := []string{"关卡", "关数", "攻击", "生命", "防御", "暴击", "暴伤", "加速回能", "虚弱增伤", "誓约增伤", "誓约回能", "搭档", "日卡", "阶数", "武器", "对谱", "时间"}
	row := make([]interface{}, len(keys))
	for i, key := range keys {
		row[i] = record[key]
	}

	_, err := srv.Spreadsheets.Values.Append(spreadsheetID, sheetName+"!A1", &sheets.ValueRange{
		Values: [][]interface{}{row},
	}).ValueInputOption("RAW").Do()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "写入失败", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

func getRecords(c *gin.Context) {
	levelType := c.Query("type")
	levelNum := c.Query("level")
	offsetStr := c.DefaultQuery("offset", "0")
	offset, _ := strconv.Atoi(offsetStr)

	// Step 1: 快速读取关卡和关数列，缩小目标范围
	values, err := getCachedABData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "读取失败", "detail": err.Error()})
		return
	}

	// Step 2: 找出匹配的行号
	var matchedRows []int
	for i, row := range values {
		if len(row) >= 2 && fmt.Sprintf("%v", row[0]) == levelType && fmt.Sprintf("%v", row[1]) == levelNum {
			matchedRows = append(matchedRows, i+2) // +2 因为 A2 是第二行
		}
	}

	if offset >= len(matchedRows) {
		c.JSON(http.StatusOK, []map[string]string{}) // 超出数据返回空
		return
	}

	// Step 3: 取出分页范围
	end := offset + 10
	if end > len(matchedRows) {
		end = len(matchedRows)
	}
	targetRows := matchedRows[offset:end]

	// Step 4: 批量读取目标行的完整内容
	var ranges []string
	for _, r := range targetRows {
		ranges = append(ranges, fmt.Sprintf("%s!A%d:Z%d", sheetName, r, r))
	}

	batchResp, err := srv.Spreadsheets.Values.BatchGet(spreadsheetID).Ranges(ranges...).Do()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "读取记录失败", "detail": err.Error()})
		return
	}

	keys := []string{"关卡", "关数", "攻击", "生命", "防御", "对谱", "暴击", "暴伤", "加速回能", "虚弱增伤", "誓约增伤", "誓约回能", "搭档", "日卡", "阶数", "武器", "时间"}

	var result []map[string]string
	for _, valueRange := range batchResp.ValueRanges {
		if len(valueRange.Values) == 0 {
			continue
		}
		row := valueRange.Values[0]
		record := map[string]string{}
		for i, key := range keys {
			if i < len(row) {
				record[key] = fmt.Sprintf("%v", row[i])
			} else {
				record[key] = ""
			}
		}
		result = append(result, record)
	}

	// Step 5: 排序（可选）
	sort.Slice(result, func(i, j int) bool {
		return result[i]["时间"] > result[j]["时间"] // 时间字段为 ISO 格式可直接比较
	})

	c.JSON(http.StatusOK, result)
}

func getLastRecords(c *gin.Context) {
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetID, sheetName+"!A2:Z").Do()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "读取失败", "detail": err.Error()})
		return
	}

	keys := []string{"关卡", "关数", "攻击", "生命", "防御", "对谱", "暴击", "暴伤", "加速回能", "虚弱增伤", "誓约增伤", "誓约回能", "搭档", "日卡", "阶数", "武器", "时间"}

	records := resp.Values
	n := len(records)
	latest := [][]interface{}{}
	if n >= 5 {
		latest = records[n-5:]
	} else {
		latest = records
	}

	// 拼接为 JSON 返回
	var result []map[string]string
	for _, row := range latest {
		record := map[string]string{}
		for i, key := range keys {
			if i < len(row) {
				record[key] = fmt.Sprintf("%v", row[i])
			} else {
				record[key] = ""
			}
		}
		result = append(result, record)
	}

	sort.Slice(result, func(i, j int) bool {
		t1, _ := time.Parse(time.RFC3339, result[i]["时间"])
		t2, _ := time.Parse(time.RFC3339, result[j]["时间"])
		return t2.Before(t1)
	})

	c.JSON(http.StatusOK, result)
}

func getRecordCount(c *gin.Context) {
	levelType := c.Query("type")
	number := c.Query("level")

	readRange := sheetName + "!A1:Z"
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetID, readRange).Do()
	if err != nil || len(resp.Values) < 2 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "读取失败", "detail": err.Error()})
		return
	}

	headers := resp.Values[0]
	idxMap := make(map[string]int)
	for i, h := range headers {
		idxMap[fmt.Sprintf("%v", h)] = i
	}

	count := 0
	for _, row := range resp.Values[1:] {
		if fmt.Sprintf("%v", row[idxMap["关卡"]]) == levelType && fmt.Sprintf("%v", row[idxMap["关数"]]) == number {
			count++
		}
	}

	c.JSON(http.StatusOK, gin.H{"count": count})
}
