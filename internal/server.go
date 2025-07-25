package internal

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"lysk-battle-record/internal/datastores"
	"lysk-battle-record/internal/models"
	"lysk-battle-record/internal/sheet_clients"
)

type Server interface {
	Ping(c *gin.Context)
	Login(c *gin.Context)
	AuthMiddleware() gin.HandlerFunc

	// Orbit Records
	ProcessOrbitRecord(c *gin.Context)
	UpdateOrbitRecord(c *gin.Context)
	GetOrbitRecords(c *gin.Context)
	DeleteOrbitRecord(c *gin.Context)

	// Championships Records
	ProcessChampionshipsRecord(c *gin.Context)
	UpdateChampionshipsRecord(c *gin.Context)
	GetChampionshipsRecords(c *gin.Context)
	DeleteChampionshipsRecord(c *gin.Context)

	// Latest Records
	GetLatestOrbitRecords(c *gin.Context)
	GetLatestChampionshipsRecords(c *gin.Context)

	// My Records
	GetAllMyOrbitRecords(c *gin.Context)
	GetMyOrbitRecords(c *gin.Context)

	// Users
	CreateUser(c *gin.Context)
	GetUser(c *gin.Context)
	UpdateUser(c *gin.Context)

	GetRanking(c *gin.Context)
}

func InitLyskServer(orbitRecordStore datastores.RecordStore, orbitSheetClient sheet_clients.RecordSheetClient,
	championshipsRecordStore datastores.RecordStore, championshipsSheetClient sheet_clients.RecordSheetClient,
	userStore datastores.UserStore, userSheetClient sheet_clients.UserSheetClient, auth *Authenticator) Server {

	return &LyskServer{
		orbitRecordStore:         orbitRecordStore,
		orbitSheetClient:         orbitSheetClient,
		championshipsRecordStore: championshipsRecordStore,
		championshipsSheetClient: championshipsSheetClient,
		userStore:                userStore,
		userSheetClient:          userSheetClient,
		auth:                     auth,
	}
}

type LyskServer struct {
	orbitRecordStore         datastores.RecordStore
	orbitSheetClient         sheet_clients.RecordSheetClient
	championshipsRecordStore datastores.RecordStore
	championshipsSheetClient sheet_clients.RecordSheetClient
	userStore                datastores.UserStore
	userSheetClient          sheet_clients.UserSheetClient
	auth                     *Authenticator

	Lottery *Lottery
}

func (s *LyskServer) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

func (s *LyskServer) Login(c *gin.Context) {
	var req struct {
		Code string `json:"code"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	token, err := s.auth.Login(req.Code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (s *LyskServer) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Next()
			return
		}

		tokenString := parts[1]
		userID, err := s.auth.ValidateJWT(tokenString)
		if err != nil {
			c.Next()
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}

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

	record.LevelType = getValue(input, "关卡")
	record.LevelNumber = getValue(input, "关数")
	record.LevelMode = getValue(input, "模式")
	record.Attack = getValue(input, "攻击")
	record.HP = getValue(input, "生命")
	record.Defense = getValue(input, "防御")
	record.Matching = getValue(input, "对谱")
	record.MatchingBuff = getValue(input, "对谱加成")
	record.CritRate = getValue(input, "暴击")
	record.CritDmg = getValue(input, "暴伤")
	record.EnergyRegen = getValue(input, "加速回能")
	record.WeakenBoost = getValue(input, "虚弱增伤")
	record.OathBoost = getValue(input, "誓约增伤")
	record.OathRegen = getValue(input, "誓约回能")
	record.Partner = getValue(input, "搭档身份")
	record.SetCard = getValue(input, "日卡")
	record.Stage = getValue(input, "阶数")
	record.Weapon = getValue(input, "武器")
	record.TotalLevel = getValue(input, "卡总等级")
	record.Note = getValue(input, "备注")

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
	c.JSON(http.StatusOK, record)
}

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
	c.JSON(http.StatusOK, record)
}

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

	record.LevelType = getValue(input, "关卡")
	record.Attack = getValue(input, "攻击")
	record.HP = getValue(input, "生命")
	record.Defense = getValue(input, "防御")
	record.Matching = getValue(input, "对谱")
	record.MatchingBuff = getValue(input, "对谱加成")
	record.CritRate = getValue(input, "暴击")
	record.CritDmg = getValue(input, "暴伤")
	record.EnergyRegen = getValue(input, "加速回能")
	record.WeakenBoost = getValue(input, "虚弱增伤")
	record.OathBoost = getValue(input, "誓约增伤")
	record.OathRegen = getValue(input, "誓约回能")
	record.Partner = getValue(input, "搭档身份")
	record.SetCard = getValue(input, "日卡")
	record.Stage = getValue(input, "阶数")
	record.Weapon = getValue(input, "武器")
	record.Buff = getValue(input, "加成")
	record.TotalLevel = getValue(input, "卡总等级")
	record.Note = getValue(input, "备注")

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

func (s *LyskServer) GetChampionshipsRecords(c *gin.Context) {
	level := c.Query("level")
	offsetStr := c.DefaultQuery("offset", "0")
	offset, _ := strconv.Atoi(offsetStr)

	record := s.championshipsRecordStore.Query(datastores.QueryOptions{
		Filters: map[string]string{
			"关卡": level,
		},
		Offset: offset,
	})
	c.JSON(http.StatusOK, record)
}

func (s *LyskServer) GetLatestOrbitRecords(c *gin.Context) {
	record := s.orbitRecordStore.Query(datastores.QueryOptions{
		Limit: 5,
	})
	c.JSON(http.StatusOK, record)
}

func (s *LyskServer) GetLatestChampionshipsRecords(c *gin.Context) {
	record := s.championshipsRecordStore.Query(datastores.QueryOptions{
		Limit: 5,
	})
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
	c.JSON(http.StatusOK, record)
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

	record.LevelType = getValue(input, "关卡")
	record.LevelNumber = getValue(input, "关数")
	record.LevelMode = getValue(input, "模式")
	record.Attack = getValue(input, "攻击")
	record.HP = getValue(input, "生命")
	record.Defense = getValue(input, "防御")
	record.Matching = getValue(input, "对谱")
	record.MatchingBuff = getValue(input, "对谱加成")
	record.CritRate = getValue(input, "暴击")
	record.CritDmg = getValue(input, "暴伤")
	record.EnergyRegen = getValue(input, "加速回能")
	record.WeakenBoost = getValue(input, "虚弱增伤")
	record.OathBoost = getValue(input, "誓约增伤")
	record.OathRegen = getValue(input, "誓约回能")
	record.Partner = getValue(input, "搭档身份")
	record.SetCard = getValue(input, "日卡")
	record.Stage = getValue(input, "阶数")
	record.Weapon = getValue(input, "武器")
	record.TotalLevel = getValue(input, "卡总等级")
	record.Note = getValue(input, "备注")

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

	record.LevelType = getValue(input, "关卡")
	record.Attack = getValue(input, "攻击")
	record.HP = getValue(input, "生命")
	record.Defense = getValue(input, "防御")
	record.Matching = getValue(input, "对谱")
	record.MatchingBuff = getValue(input, "对谱加成")
	record.CritRate = getValue(input, "暴击")
	record.CritDmg = getValue(input, "暴伤")
	record.EnergyRegen = getValue(input, "加速回能")
	record.WeakenBoost = getValue(input, "虚弱增伤")
	record.OathBoost = getValue(input, "誓约增伤")
	record.OathRegen = getValue(input, "誓约回能")
	record.Partner = getValue(input, "搭档身份")
	record.SetCard = getValue(input, "日卡")
	record.Stage = getValue(input, "阶数")
	record.Weapon = getValue(input, "武器")
	record.Buff = getValue(input, "加成")
	record.TotalLevel = getValue(input, "卡总等级")
	record.Note = getValue(input, "备注")

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

func (s *LyskServer) GetRanking(c *gin.Context) {
	userId, exists := c.Get("userID")
	if !exists {
		userId = ""
	}
	ranking := s.orbitRecordStore.GetRanking(userId.(string))

	c.JSON(http.StatusOK, ranking)
}

func (s *LyskServer) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	userId, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录或无效的用户"})
		return
	}
	user.ID = userId.(string)

	user, ok := s.userStore.Get(userId.(string))
	if ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户已存在"})
		return
	}

	createdUser, err := s.userSheetClient.ProcessUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	s.userStore.Insert(*createdUser)
	c.JSON(http.StatusOK, createdUser)
}

func (s *LyskServer) GetUser(c *gin.Context) {
	userId, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录或无效的用户"})
		return
	}
	user, ok := s.userStore.Get(userId.(string))
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (s *LyskServer) UpdateUser(c *gin.Context) {
	userId, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录或无效的用户"})
		return
	}
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	user.ID = userId.(string)

	currentUser, ok := s.userStore.Get(userId.(string))
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	user.RowNumber = currentUser.RowNumber

	if err := s.userSheetClient.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := s.userStore.Update(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func getValue(input map[string]interface{}, field string) string {
	if val, ok := input[field]; ok {
		return fmt.Sprintf("%v", val)
	}

	return ""
}
