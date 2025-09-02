package usecases

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"lysk-battle-record/internal/datastores"
	"lysk-battle-record/internal/models"
	"lysk-battle-record/internal/pkg"
)

// All Orbit Records

func (s *LyskServer) GetOrbitRecords(c *gin.Context) {
	levelType := c.Query("type")
	levelNum := c.Query("level")
	levelMode := c.Query("mode")
	offsetStr := c.DefaultQuery("offset", "0")
	offset, _ := strconv.Atoi(offsetStr)

	if levelMode == "" {
		levelMode = "稳定"
	}

	record := s.orbitRecordStore.Query(datastores.QueryOptions{
		Filters: map[string]string{
			"关卡": levelType,
			"关数": levelNum,
			"模式": levelMode,
		},
		Offset: offset,
	})
	s.populateNicknameForRecords(record.Records)
	c.JSON(http.StatusOK, record)
}

func (s *LyskServer) GetLatestOrbitRecords(c *gin.Context) {
	record := s.orbitRecordStore.Query(datastores.QueryOptions{
		Limit: 5,
	})
	s.populateNicknameForRecords(record.Records)
	c.JSON(http.StatusOK, record)
}

// My Orbit Records

func (s *LyskServer) GetMyOrbitRecords(c *gin.Context) {
	levelType := c.Query("type")
	levelNum := c.Query("level")
	levelMode := c.Query("mode")
	offsetStr := c.DefaultQuery("offset", "0")
	offset, _ := strconv.Atoi(offsetStr)
	userId, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录或无效的用户"})
		return
	}

	if levelMode == "" {
		levelMode = "稳定"
	}

	record := s.orbitRecordStore.Query(datastores.QueryOptions{
		Filters: map[string]string{
			"关卡":   levelType,
			"关数":   levelNum,
			"模式":   levelMode,
			"用户ID": userId.(string),
		},
		Offset: offset,
	})
	s.populateNicknameForRecords(record.Records)
	c.JSON(http.StatusOK, record)
}

func (s *LyskServer) GetAllMyOrbitRecords(c *gin.Context) {
	userId, exists := c.Get("userID")
	offsetStr := c.DefaultQuery("offset", "0")
	offset, _ := strconv.Atoi(offsetStr)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录或无效的用户"})
		return
	}

	record := s.orbitRecordStore.Query(datastores.QueryOptions{
		Filters: map[string]string{
			"用户ID": userId.(string),
		},
		Offset: offset,
	})
	s.populateNicknameForRecords(record.Records)
	c.JSON(http.StatusOK, record)
}

// CRUD methods

func (s *LyskServer) ProcessOrbitRecord(c *gin.Context) {
	var input map[string]interface{}
	if err := c.BindJSON(&input); err != nil {
		logrus.Errorf("[Orbit] Failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求格式错误", "detail": err.Error()})
		return
	}

	record := models.Record{}
	if userID, exists := c.Get("userID"); exists {
		record.UserID = userID.(string)
	}

	record.LevelType = pkg.GetValue(input, "关卡")
	record.LevelNumber = pkg.GetValue(input, "关数")
	record.LevelMode = pkg.GetValue(input, "模式")
	record.Attack = pkg.GetValue(input, "攻击")
	record.HP = pkg.GetValue(input, "生命")
	record.Defense = pkg.GetValue(input, "防御")
	record.Matching = pkg.GetValue(input, "对谱")
	record.MatchingBuff = pkg.GetValue(input, "对谱加成")
	record.CritRate = pkg.GetValue(input, "暴击")
	record.CritDmg = pkg.GetValue(input, "暴伤")
	record.EnergyRegen = pkg.GetValue(input, "加速回能")
	record.WeakenBoost = pkg.GetValue(input, "虚弱增伤")
	record.OathBoost = pkg.GetValue(input, "誓约增伤")
	record.OathRegen = pkg.GetValue(input, "誓约回能")
	record.Companion = pkg.GetValue(input, "搭档身份")
	record.SetCard = pkg.GetValue(input, "日卡")
	record.Stage = pkg.GetValue(input, "阶数")
	record.Weapon = pkg.GetValue(input, "武器")
	record.TotalLevel = pkg.GetValue(input, "卡总等级")
	record.Note = pkg.GetValue(input, "备注")
	record.StarRank = cleanUpStarRankValue(pkg.GetValue(input, "星级"))

	if _, err := record.ValidateOrbit(); err != nil {
		logrus.Errorf("[Orbit] Record validation failed: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if s.orbitRecordStore.IsDuplicate(record) {
		logrus.Errorf("[Orbit] Record is duplicated")
		c.JSON(http.StatusBadRequest, gin.H{"error": "记录已存在"})
		return
	}

	err := s.orbitRecordStore.PrepareInsert(record)
	if err != nil {
		logrus.Errorf("[Orbit] Failed to prepare record for insertion: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "记录上传失败", "detail": err.Error()})
		return
	}

	// 时间字段处理为 ISO 格式
	if t, ok := input["时间"].(string); ok {
		if parsedTime, err := time.Parse(time.RFC3339, t); err == nil {
			record.Time = parsedTime.Format(time.RFC3339)
		} else {
			logrus.Errorf("[Orbit] Failed to parse time: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "时间格式错误"})
			return
		}
	} else {
		logrus.Error("[Orbit] Time field is missing or in wrong format")
		c.JSON(http.StatusBadRequest, gin.H{"error": "时间字段缺失或格式错误"})
		return
	}

	ingestedRecord, err := s.orbitSheetClient.ProcessRecord(record)
	if err != nil {
		logrus.Errorf("[Orbit] Failed to write record to Google Sheet: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "写入失败", "detail": err.Error()})
		return
	}

	s.orbitRecordStore.Insert(*ingestedRecord)

	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

func (s *LyskServer) UpdateOrbitRecord(c *gin.Context) {
	var input map[string]interface{}
	if err := c.BindJSON(&input); err != nil {
		logrus.Errorf("[Orbit] Failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求格式错误", "detail": err.Error()})
		return
	}

	recordId := c.Param("id")
	existingRecord, ok := s.orbitRecordStore.Get(recordId)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "记录不存在"})
		return
	}

	if existingRecord.Deleted {
		c.JSON(http.StatusNotFound, gin.H{"error": "记录已被删除"})
		return
	}

	userId, exists := c.Get("userID")
	if !exists || userId.(string) != existingRecord.UserID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无权修改此记录"})
		return
	}

	record := models.Record{}
	record.Id = recordId
	record.RowNumber = existingRecord.RowNumber
	record.UserID = existingRecord.UserID
	record.Time = existingRecord.Time

	record.LevelType = existingRecord.LevelType
	record.LevelNumber = existingRecord.LevelNumber
	record.LevelMode = existingRecord.LevelMode
	record.Attack = pkg.GetValue(input, "攻击")
	record.HP = pkg.GetValue(input, "生命")
	record.Defense = pkg.GetValue(input, "防御")
	record.Matching = pkg.GetValue(input, "对谱")
	record.MatchingBuff = pkg.GetValue(input, "对谱加成")
	record.CritRate = pkg.GetValue(input, "暴击")
	record.CritDmg = pkg.GetValue(input, "暴伤")
	record.EnergyRegen = pkg.GetValue(input, "加速回能")
	record.WeakenBoost = pkg.GetValue(input, "虚弱增伤")
	record.OathBoost = pkg.GetValue(input, "誓约增伤")
	record.OathRegen = pkg.GetValue(input, "誓约回能")
	record.Companion = pkg.GetValue(input, "搭档身份")
	record.SetCard = pkg.GetValue(input, "日卡")
	record.Stage = pkg.GetValue(input, "阶数")
	record.Weapon = pkg.GetValue(input, "武器")
	record.TotalLevel = pkg.GetValue(input, "卡总等级")
	record.Note = pkg.GetValue(input, "备注")
	record.StarRank = cleanUpStarRankValue(pkg.GetValue(input, "星级"))

	if _, err := record.ValidateOrbit(); err != nil {
		logrus.Errorf("[Orbit] Record validation failed: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.orbitSheetClient.UpdateRecord(record); err != nil {
		logrus.Errorf("[Orbit] Failed to update record in Google Sheet: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败", "detail": err.Error()})
		return
	}

	if err := s.orbitRecordStore.Update(record); err != nil {
		logrus.Errorf("[Orbit] Failed to update record in memory: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

func (s *LyskServer) DeleteOrbitRecord(c *gin.Context) {
	recordId := c.Param("id")
	existingRecord, ok := s.orbitRecordStore.Get(recordId)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "记录不存在"})
		return
	}

	if existingRecord.Deleted {
		c.JSON(http.StatusNotFound, gin.H{"error": "记录已被删除"})
		return
	}

	userId, exists := c.Get("userID")
	if !exists || userId.(string) != existingRecord.UserID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无权删除此记录"})
		return
	}

	if err := s.orbitSheetClient.DeleteRecord(existingRecord); err != nil {
		logrus.Errorf("[Orbit] Failed to delete record from Google Sheet: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败", "detail": err.Error()})
		return
	}

	if err := s.orbitRecordStore.Delete(existingRecord); err != nil {
		logrus.Errorf("[Orbit] Failed to delete record from memory: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

func cleanUpStarRankValue(value string) string {
	validValue := map[string]bool{
		"零星": true,
		"一星": true,
		"二星": true,
		"三星": true,
	}
	if _, exists := validValue[value]; exists {
		return value
	}

	return ""
}
