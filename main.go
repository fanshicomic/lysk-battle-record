package main

import (
	"context"
	"github.com/gin-contrib/cors"
	"sort"
	"strconv"
	"time"

	//"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

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

	r.POST("/record", processRecord)
	r.GET("/records", getRecords)
	r.GET("/record-count", getRecordCount)

	r.Run(":8080")
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
	number := c.Query("level")
	offsetStr := c.DefaultQuery("offset", "0")
	offset, _ := strconv.Atoi(offsetStr)

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

	var result []map[string]string
	for _, row := range resp.Values[1:] {
		if fmt.Sprintf("%v", row[idxMap["关卡"]]) == levelType && fmt.Sprintf("%v", row[idxMap["关数"]]) == number {
			item := make(map[string]string)
			for key, idx := range idxMap {
				if idx < len(row) {
					item[key] = fmt.Sprintf("%v", row[idx])
				}
			}
			result = append(result, item)
		}
	}

	sort.Slice(result, func(i, j int) bool {
		t1, _ := time.Parse(time.RFC3339, result[i]["时间"])
		t2, _ := time.Parse(time.RFC3339, result[j]["时间"])
		return t2.Before(t1)
	})

	pageSize := 10
	start := offset
	end := start + pageSize
	if start >= len(result) {
		c.JSON(http.StatusOK, []map[string]string{})
		return
	}
	if end > len(result) {
		end = len(result)
	}
	c.JSON(http.StatusOK, result[start:end])
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
