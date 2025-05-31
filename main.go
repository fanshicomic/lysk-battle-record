package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"lysk-battle-record/internal"
)

type Server struct {
	recordStore internal.RecordStore
	sheetClient internal.GoogleSheetClient
}

func main() {
	googleSheetClient := internal.NewGoogleSheetClient()
	recordStore := internal.NewInMemoryRecordStore(googleSheetClient)

	server := &Server{
		recordStore: recordStore,
		sheetClient: googleSheetClient,
	}

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // 可改为指定域名，如 "https://yourdomain.com"
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           1 * time.Hour,
	}))

	r.GET("/ping", server.ping)

	r.POST("/record", server.processRecord)
	r.GET("/records", server.getRecords)
	r.GET("/last-records", server.getLastRecords)

	r.Run(":8080")
}

func (s *Server) ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

func (s *Server) processRecord(c *gin.Context) {
	var input map[string]interface{}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求格式错误", "detail": err.Error()})
		return
	}

	record := internal.Record{}
	record.LevelType = fmt.Sprintf("%v", input["关卡"])
	record.LevelNumber = fmt.Sprintf("%v", input["关数"])
	record.Attack = fmt.Sprintf("%v", input["攻击"])
	record.HP = fmt.Sprintf("%v", input["生命"])
	record.Defense = fmt.Sprintf("%v", input["防御"])
	record.Matching = fmt.Sprintf("%v", input["对谱"])
	record.CritRate = fmt.Sprintf("%v", input["暴击"])
	record.CritDmg = fmt.Sprintf("%v", input["暴伤"])
	record.EnergyRegen = fmt.Sprintf("%v", input["加速回能"])
	record.WeakenBoost = fmt.Sprintf("%v", input["虚弱增伤"])
	record.OathBoost = fmt.Sprintf("%v", input["誓约增伤"])
	record.OathRegen = fmt.Sprintf("%v", input["誓约回能"])
	record.Partner = fmt.Sprintf("%v", input["搭档身份"])
	record.SetCard = fmt.Sprintf("%v", input["日卡"])
	record.Stage = fmt.Sprintf("%v", input["阶数"])
	record.Weapon = fmt.Sprintf("%v", input["武器"])

	if s.recordStore.IsDuplicate(record) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "记录已存在"})
		return
	}

	err := s.recordStore.PrepareInsert(record)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "记录上传失败", "detail": err.Error()})
		return
	}

	// 时间字段处理为 ISO 格式
	if t, ok := input["时间"].(string); ok {
		if parsedTime, err := time.Parse(time.RFC3339, t); err == nil {
			record.Time = parsedTime.Format(time.RFC3339)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "时间格式错误"})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "时间字段缺失或格式错误"})
		return
	}

	if err := s.sheetClient.ProcessRecord(record); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "写入失败", "detail": err.Error()})
		return
	}

	s.recordStore.Insert(record)

	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

func (s *Server) getRecords(c *gin.Context) {
	levelType := c.Query("type")
	levelNum := c.Query("level")
	offsetStr := c.DefaultQuery("offset", "0")
	offset, _ := strconv.Atoi(offsetStr)

	record := s.recordStore.Query(internal.QueryOptions{
		Filters: map[string]string{
			"关卡": levelType,
			"关数": levelNum,
		},
		Offset: offset,
	})
	c.JSON(http.StatusOK, record)
}

func (s *Server) getLastRecords(c *gin.Context) {
	record := s.recordStore.Query(internal.QueryOptions{
		Limit: 5,
	})
	c.JSON(http.StatusOK, record)
}
