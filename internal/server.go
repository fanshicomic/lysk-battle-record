package internal

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Server interface {
	Ping(c *gin.Context)

	ProcessOrbitRecord(c *gin.Context)
	GetOrbitRecords(c *gin.Context)
	GetLatestOrbitRecords(c *gin.Context)

	ProcessChampionshipsRecord(c *gin.Context)
	GetChampionshipsRecords(c *gin.Context)
	GetLatestChampionshipsRecords(c *gin.Context)
}

func InitLyskServer(orbitRecordStore RecordStore, orbitSheetClient GoogleSheetClient,
	championshipsRecordStore RecordStore, championshipsSheetClient GoogleSheetClient) Server {

	return &LyskServer{
		orbitRecordStore:         orbitRecordStore,
		orbitSheetClient:         orbitSheetClient,
		championshipsRecordStore: championshipsRecordStore,
		championshipsSheetClient: championshipsSheetClient,
	}
}

type LyskServer struct {
	orbitRecordStore RecordStore
	orbitSheetClient GoogleSheetClient

	championshipsRecordStore RecordStore
	championshipsSheetClient GoogleSheetClient

	Lottery *Lottery
}

func (s *LyskServer) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

func (s *LyskServer) ProcessOrbitRecord(c *gin.Context) {
	var input map[string]interface{}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求格式错误", "detail": err.Error()})
		return
	}

	record := Record{}
	record.LevelType = fmt.Sprintf("%v", input["关卡"])
	record.LevelNumber = fmt.Sprintf("%v", input["关数"])
	record.Attack = fmt.Sprintf("%v", input["攻击"])
	record.HP = fmt.Sprintf("%v", input["生命"])
	record.Defense = fmt.Sprintf("%v", input["防御"])
	record.Matching = fmt.Sprintf("%v", input["对谱"])
	record.MatchingBuffer = fmt.Sprintf("%v", input["对谱加成"])
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

	if _, err := record.ValidateOrbit(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

func (s *LyskServer) GetOrbitRecords(c *gin.Context) {
	levelType := c.Query("type")
	levelNum := c.Query("level")
	offsetStr := c.DefaultQuery("offset", "0")
	offset, _ := strconv.Atoi(offsetStr)

	record := s.orbitRecordStore.Query(QueryOptions{
		Filters: map[string]string{
			"关卡": levelType,
			"关数": levelNum,
		},
		Offset: offset,
	})
	c.JSON(http.StatusOK, record)
}

func (s *LyskServer) ProcessChampionshipsRecord(c *gin.Context) {
	var input map[string]interface{}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求格式错误", "detail": err.Error()})
		return
	}

	record := Record{}
	record.LevelType = fmt.Sprintf("%v", input["关卡"])
	record.Attack = fmt.Sprintf("%v", input["攻击"])
	record.HP = fmt.Sprintf("%v", input["生命"])
	record.Defense = fmt.Sprintf("%v", input["防御"])
	record.Matching = fmt.Sprintf("%v", input["对谱"])
	record.MatchingBuffer = fmt.Sprintf("%v", input["对谱加成"])
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
	record.Buffer = fmt.Sprintf("%v", input["加成"])

	if _, err := record.ValidateChampionships(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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

func (s *LyskServer) GetChampionshipsRecords(c *gin.Context) {
	level := c.Query("level")
	offsetStr := c.DefaultQuery("offset", "0")
	offset, _ := strconv.Atoi(offsetStr)

	record := s.championshipsRecordStore.Query(QueryOptions{
		Filters: map[string]string{
			"关卡": level,
		},
		Offset: offset,
	})
	c.JSON(http.StatusOK, record)
}

func (s *LyskServer) GetLatestOrbitRecords(c *gin.Context) {
	record := s.orbitRecordStore.Query(QueryOptions{
		Limit: 5,
	})
	c.JSON(http.StatusOK, record)
}

func (s *LyskServer) GetLatestChampionshipsRecords(c *gin.Context) {
	record := s.championshipsRecordStore.Query(QueryOptions{
		Limit: 5,
	})
	c.JSON(http.StatusOK, record)
}
