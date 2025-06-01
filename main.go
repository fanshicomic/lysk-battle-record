package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"lysk-battle-record/internal"
)

type Server struct {
	orbitRecordStore internal.RecordStore
	orbitSheetClient internal.GoogleSheetClient

	championshipsRecordStore internal.RecordStore
	championshipsSheetClient internal.GoogleSheetClient
}

const (
	spreadsheetID     = "1-ORnXBnav4JVtP673Oio5sNdVpk0taUSzG3kWZqhIuY"
	orbitSheetName    = "轨道"
	championSheetName = "锦标赛"
)

func main() {
	orbitGoogleSheetClient := internal.NewGoogleSheetClient(spreadsheetID, orbitSheetName)
	orbitRecordStore := internal.NewInMemoryRecordStore(orbitGoogleSheetClient)

	championshipsGoogleSheetClient := internal.NewGoogleSheetClient(spreadsheetID, championSheetName)
	championshipsRecordStore := internal.NewInMemoryRecordStore(championshipsGoogleSheetClient)

	server := &Server{
		orbitRecordStore:         orbitRecordStore,
		orbitSheetClient:         orbitGoogleSheetClient,
		championshipsRecordStore: championshipsRecordStore,
		championshipsSheetClient: championshipsGoogleSheetClient,
	}

	r := gin.Default()

	allowOrigins := []string{"https://uygnim.com"}
	if isLocal() {
		allowOrigins = []string{"http://localhost:63343"}
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins:     allowOrigins,
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           1 * time.Hour,
	}))

	r.GET("/ping", server.ping)

	r.POST("/orbit-record", server.processOrbitRecord)
	r.GET("/orbit-records", server.getOrbitRecords)

	r.POST("/championships-record", server.processChampionshipsRecord)
	r.GET("/championships-records", server.getChampionshipsRecords)

	r.GET("/latest-orbit-records", server.getLatestOrbitRecords)
	r.GET("/latest-championships-records", server.getLatestChampionshipsRecords)

	r.Run(":8080")
}

func (s *Server) ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

func (s *Server) processOrbitRecord(c *gin.Context) {
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

	if !record.Validate() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "记录数据错误"})
		return
	}

	if s.orbitRecordStore.IsDuplicate(record) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "记录已存在"})
		return
	}

	err := s.orbitRecordStore.PrepareInsert(record)
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

	if err := s.orbitSheetClient.ProcessRecord(record); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "写入失败", "detail": err.Error()})
		return
	}

	s.orbitRecordStore.Insert(record)

	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

func (s *Server) getOrbitRecords(c *gin.Context) {
	levelType := c.Query("type")
	levelNum := c.Query("level")
	offsetStr := c.DefaultQuery("offset", "0")
	offset, _ := strconv.Atoi(offsetStr)

	record := s.orbitRecordStore.Query(internal.QueryOptions{
		Filters: map[string]string{
			"关卡": levelType,
			"关数": levelNum,
		},
		Offset: offset,
	})
	c.JSON(http.StatusOK, record)
}

func (s *Server) processChampionshipsRecord(c *gin.Context) {
	var input map[string]interface{}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求格式错误", "detail": err.Error()})
		return
	}

	record := internal.Record{}
	record.LevelType = fmt.Sprintf("%v", input["关卡"])
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

	if s.championshipsRecordStore.IsDuplicate(record) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "记录已存在"})
		return
	}

	err := s.championshipsRecordStore.PrepareInsert(record)
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

	if err := s.championshipsSheetClient.ProcessRecord(record); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "写入失败", "detail": err.Error()})
		return
	}

	s.championshipsRecordStore.Insert(record)

	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

func (s *Server) getChampionshipsRecords(c *gin.Context) {
	level := c.Query("level")
	offsetStr := c.DefaultQuery("offset", "0")
	offset, _ := strconv.Atoi(offsetStr)

	record := s.championshipsRecordStore.Query(internal.QueryOptions{
		Filters: map[string]string{
			"关卡": level,
		},
		Offset: offset,
	})
	c.JSON(http.StatusOK, record)
}

func (s *Server) getLatestOrbitRecords(c *gin.Context) {
	record := s.orbitRecordStore.Query(internal.QueryOptions{
		Limit: 5,
	})
	c.JSON(http.StatusOK, record)
}

func (s *Server) getLatestChampionshipsRecords(c *gin.Context) {
	record := s.championshipsRecordStore.Query(internal.QueryOptions{
		Limit: 5,
	})
	c.JSON(http.StatusOK, record)
}

func isLocal() bool {
	if _, err := os.ReadFile("credentials.json"); err == nil {
		return true
	}

	return false
}
