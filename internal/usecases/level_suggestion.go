package usecases

import (
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"lysk-battle-record/internal/datastores"
	"lysk-battle-record/internal/models"
	"lysk-battle-record/internal/utils"
)

type LevelCPResponse struct {
	LevelType   string `json:"level_type"`
	LevelNumber string `json:"level_number"`
	LevelMode   string `json:"level_mode"`
	CPs         []int  `json:"cps"`
}

func (s *LyskServer) GetLevelCPs(records []models.Record) []int {
	// Extract buffed combat power from all records
	cps := make([]int, 0, len(records))
	for _, record := range records {
		if record.CombatPower.BuffedScore != "" && record.CombatPower.BuffedScore != "无数据" {
			if score, err := strconv.Atoi(record.CombatPower.BuffedScore); err == nil && score > 0 {
				cps = append(cps, score)
			}
		}
	}

	// Sort CPs in descending order
	sort.Sort(sort.Reverse(sort.IntSlice(cps)))

	return cps
}

func (s *LyskServer) GetLevelSuggestion(c *gin.Context) {
	levelType, levelNumber, levelMode := c.Query("type"), c.Query("level"), c.Query("mode")
	if levelNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "关卡参数不能为空"})
		return
	}

	// Fetch records once
	isChampionships := strings.Contains(levelNumber, "A4") || strings.Contains(levelNumber, "B4") || strings.Contains(levelNumber, "C4")

	var queryOptions datastores.QueryOptions
	var store datastores.RecordStore

	if isChampionships {
		store = s.championshipsRecordStore
		start, end := utils.GetCurrentChampionshipsRound()
		queryOptions = datastores.QueryOptions{
			Filters: map[string]string{
				"关卡": levelNumber,
			},
			TimeStart: start,
			TimeEnd:   end,
		}
	} else {
		if levelMode == "" {
			levelMode = "稳定"
		}

		store = s.orbitRecordStore
		queryOptions = datastores.QueryOptions{
			Filters: map[string]string{
				"关卡": levelType,
				"关数": levelNumber,
				"模式": levelMode,
			},
		}
	}

	result := store.Query(queryOptions)

	// Get CPs for the level using the fetched records
	cps := s.GetLevelCPs(result.Records)

	// TODO: Add other suggestion methods here that can reuse result.Records

	response := LevelCPResponse{
		LevelType:   levelType,
		LevelNumber: levelNumber,
		LevelMode:   levelMode,
		CPs:         cps,
	}

	c.JSON(http.StatusOK, response)
}
