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

// All Championships Records

func getCurrentChampionshipsRound() (time.Time, time.Time) {
	// The first round started on 2025 June 02
	firstRoundStartDate := time.Date(2025, time.June, 2, 0, 0, 0, 0, time.UTC)
	// One round of championships last for 2 weeks
	roundDuration := 14 * 24 * time.Hour

	// Calculate the time elapsed since the first round
	elapsed := time.Now().UTC().Sub(firstRoundStartDate)
	// Calculate the number of rounds that have passed
	roundsPassed := int(elapsed / roundDuration)

	// Calculate the start date of the current round
	currentRoundStartDate := firstRoundStartDate.Add(time.Duration(roundsPassed) * roundDuration)
	// Calculate the end date of the current round
	currentRoundEndDate := currentRoundStartDate.Add(roundDuration)

	return currentRoundStartDate, currentRoundEndDate
}

func (s *LyskServer) GetChampionshipsRecords(c *gin.Context) {
	level := c.Query("level")
	offsetStr := c.DefaultQuery("offset", "0")
	offset, _ := strconv.Atoi(offsetStr)
	start, end := getCurrentChampionshipsRound()

	records := s.championshipsRecordStore.Query(datastores.QueryOptions{
		Filters: map[string]string{
			"关卡": level,
		},
		Offset:    offset,
		TimeStart: start,
		TimeEnd:   end,
	})
	c.JSON(http.StatusOK, records)
}

func (s *LyskServer) GetLatestChampionshipsRecords(c *gin.Context) {
	start, end := getCurrentChampionshipsRound()
	records := s.championshipsRecordStore.Query(datastores.QueryOptions{
		Limit:     5,
		TimeStart: start,
		TimeEnd:   end,
	})
	c.JSON(http.StatusOK, records)
}

// CRUD methods

func (s *LyskServer) ProcessChampionshipsRecord(c *gin.Context) {
	var input map[string]interface{}
	if err := c.BindJSON(&input); err != nil {
		logrus.Errorf("[Championships] Failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求格式错误", "detail": err.Error()})
		return
	}

	record := models.Record{}
	if userID, exists := c.Get("userID"); exists {
		record.UserID = userID.(string)
	}

	record.LevelType = pkg.GetValue(input, "关卡")
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
	record.Partner = pkg.GetValue(input, "搭档身份")
	record.SetCard = pkg.GetValue(input, "日卡")
	record.Stage = pkg.GetValue(input, "阶数")
	record.Weapon = pkg.GetValue(input, "武器")
	record.Buff = pkg.GetValue(input, "加成")
	record.TotalLevel = pkg.GetValue(input, "卡总等级")
	record.Note = pkg.GetValue(input, "备注")

	if _, err := record.ValidateChampionships(); err != nil {
		logrus.Errorf("[Championships] Record validation failed: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if s.championshipsRecordStore.IsDuplicate(record) {
		logrus.Errorf("[Championships] Record is duplicated")
		c.JSON(http.StatusBadRequest, gin.H{"error": "记录已存在"})
		return
	}

	err := s.championshipsRecordStore.PrepareInsert(record)
	if err != nil {
		logrus.Errorf("[Championships] Failed to prepare record for insertion: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "记录上传失败", "detail": err.Error()})
		return
	}

	// 时间字段处理为 ISO 格式
	if t, ok := input["时间"].(string); ok {
		if parsedTime, err := time.Parse(time.RFC3339, t); err == nil {
			record.Time = parsedTime.Format(time.RFC3339)
		} else {
			logrus.Errorf("[Championships] Failed to parse time: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "时间格式错误"})
			return
		}
	} else {
		logrus.Error("[Championships] Time field is missing or in wrong format")
		c.JSON(http.StatusBadRequest, gin.H{"error": "时间字段缺失或格式错误"})
		return
	}

	ingestedRecord, err := s.championshipsSheetClient.ProcessRecord(record)
	if err != nil {
		logrus.Errorf("[Championships] Failed to write record to Google Sheet: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "写入失败", "detail": err.Error()})
		return
	}

	s.championshipsRecordStore.Insert(*ingestedRecord)

	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

func (s *LyskServer) UpdateChampionshipsRecord(c *gin.Context) {
	var input map[string]interface{}
	if err := c.BindJSON(&input); err != nil {
		logrus.Errorf("[Championships] Failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求格式错误", "detail": err.Error()})
		return
	}

	recordId := c.Param("id")
	existingRecord, ok := s.championshipsRecordStore.Get(recordId)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "记录不存在"})
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

	record.LevelType = pkg.GetValue(input, "关卡")
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
	record.Partner = pkg.GetValue(input, "搭档身份")
	record.SetCard = pkg.GetValue(input, "日卡")
	record.Stage = pkg.GetValue(input, "阶数")
	record.Weapon = pkg.GetValue(input, "武器")
	record.Buff = pkg.GetValue(input, "加成")
	record.TotalLevel = pkg.GetValue(input, "卡总等级")
	record.Note = pkg.GetValue(input, "备注")

	if _, err := record.ValidateChampionships(); err != nil {
		logrus.Errorf("[Championships] Record validation failed: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.championshipsSheetClient.UpdateRecord(record); err != nil {
		logrus.Errorf("[Championships] Failed to update record in Google Sheet: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败", "detail": err.Error()})
		return
	}

	if err := s.championshipsRecordStore.Update(record); err != nil {
		logrus.Errorf("[Championships] Failed to update record in memory: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

func (s *LyskServer) DeleteChampionshipsRecord(c *gin.Context) {
	recordId := c.Param("id")
	existingRecord, ok := s.championshipsRecordStore.Get(recordId)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "记录不存在"})
		return
	}

	userId, exists := c.Get("userID")
	if !exists || userId.(string) != existingRecord.UserID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无权删除此记录"})
		return
	}

	if err := s.championshipsSheetClient.DeleteRecord(existingRecord); err != nil {
		logrus.Errorf("[Championships] Failed to delete record from Google Sheet: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败", "detail": err.Error()})
		return
	}

	if err := s.championshipsRecordStore.Delete(existingRecord); err != nil {
		logrus.Errorf("[Championships] Failed to delete record from memory: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}
