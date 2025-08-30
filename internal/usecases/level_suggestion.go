package usecases

import (
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"lysk-battle-record/internal/datastores"
	"lysk-battle-record/internal/models"
)

type CompanionSetCardPair struct {
	Companion string `json:"companion"`
	SetCard   string `json:"set_card"`
	Count     int    `json:"count"`
}

type LevelSuggestionResponse struct {
	LevelType             string                 `json:"level_type"`
	LevelNumber           string                 `json:"level_number"`
	LevelMode             string                 `json:"level_mode"`
	CPs                   []int                  `json:"cps"`
	SuggestedCP           int                    `json:"suggested_cp"`
	CompanionSetCardPairs []CompanionSetCardPair `json:"companion_setcard_pairs"`
	Crit                  int                    `json:"crit"`
	Weak                  int                    `json:"weak"`
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
	sort.Sort(sort.IntSlice(cps))

	return cps
}

func (s *LyskServer) GetCompanionSetCardPairs(records []models.Record) []CompanionSetCardPair {
	pairCounts := make(map[string]int)

	for _, record := range records {
		companion := record.Companion
		setCard := record.SetCard

		pairKey := companion + "<>" + setCard
		pairCounts[pairKey]++
	}

	var pairs []CompanionSetCardPair
	for pairKey, count := range pairCounts {
		parts := strings.Split(pairKey, "<>")
		pairs = append(pairs, CompanionSetCardPair{
			Companion: parts[0],
			SetCard:   parts[1],
			Count:     count,
		})
	}

	return pairs
}

func (s *LyskServer) GetSuggestedCP(cps []int) int {
	if len(cps) == 0 {
		return 0
	}

	// Sort CPs in ascending order for percentile calculation
	sortedCPs := make([]int, len(cps))
	copy(sortedCPs, cps)
	sort.Ints(sortedCPs)

	// Calculate 25th percentile index
	index := int(float64(len(sortedCPs)) * 0.25)
	if index >= len(sortedCPs) {
		index = len(sortedCPs) - 1
	}

	return sortedCPs[index]
}

func (s *LyskServer) GetLevelSuggestion(c *gin.Context) {
	levelType, levelNumber, levelMode := c.Query("type"), c.Query("level"), c.Query("mode")
	if levelNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "关卡参数不能为空"})
		return
	}

	// Create a temporary record to use with GetLevelRecords
	tempRecord := models.Record{
		LevelType:   levelType,
		LevelNumber: levelNumber,
		LevelMode:   levelMode,
		Time:        time.Now().Format(time.RFC3339),
	}

	// Determine which store to use
	isChampionships := strings.Contains(levelNumber, "A4") || strings.Contains(levelNumber, "B4") || strings.Contains(levelNumber, "C4")
	var store datastores.RecordStore
	if isChampionships {
		store = s.championshipsRecordStore
	} else {
		if levelMode == "" {
			levelMode = "稳定"
			tempRecord.LevelMode = levelMode
		}
		store = s.orbitRecordStore
	}

	// Get records for the level using the optimized method
	records := store.GetLevelRecords(tempRecord)

	// Get CPs for the level using the fetched records
	cps := s.GetLevelCPs(records)

	// Get companion and set card pairs
	companionSetCardPairs := s.GetCompanionSetCardPairs(records)

	// Calculate suggested CP (25th percentile)
	suggestedCP := s.GetSuggestedCP(cps)

	critCount := 0
	weakCount := 0

	for _, record := range records {
		if record.CombatPower.CritScore != "" && record.CombatPower.WeakenScore != "" &&
			record.CombatPower.CritScore != "无数据" && record.CombatPower.WeakenScore != "无数据" {

			critScore, critErr := strconv.Atoi(record.CombatPower.CritScore)
			weakenScore, weakenErr := strconv.Atoi(record.CombatPower.WeakenScore)

			if critErr == nil && weakenErr == nil {
				if critScore > weakenScore {
					critCount++
				} else if weakenScore > critScore {
					weakCount++
				}
			}
		}
	}

	response := LevelSuggestionResponse{
		LevelType:             levelType,
		LevelNumber:           levelNumber,
		LevelMode:             levelMode,
		CPs:                   cps,
		SuggestedCP:           suggestedCP,
		CompanionSetCardPairs: companionSetCardPairs,
		Crit:                  critCount,
		Weak:                  weakCount,
	}

	c.JSON(http.StatusOK, response)
}
