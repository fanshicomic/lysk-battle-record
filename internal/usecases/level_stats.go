package usecases

import (
	"net/http"
	"sort"
	"strconv"
	"strings"

	"lysk-battle-record/internal/models"

	"github.com/gin-gonic/gin"
)

type LevelMinCP struct {
	LevelType   string `json:"level_type"`
	LevelNumber string `json:"level_number"`
	LevelMode   string `json:"level_mode"`
	MinCP       int    `json:"min_cp"`
}

func (s *LyskServer) GetMinCombatPower(c *gin.Context) {
	levelRecords := s.orbitRecordStore.GetAllLevelRecords()
	var response []LevelMinCP

	for _, records := range levelRecords {
		if len(records) == 0 {
			continue
		}

		// Use the first record to identify the level
		firstRecord := records[0]

		minCP := -1

		for _, record := range records {
			if record.CombatPower.BuffedScore == "" || record.CombatPower.BuffedScore == models.NoData {
				continue
			}
			score, err := strconv.Atoi(record.CombatPower.BuffedScore)
			if err != nil {
				continue
			}

			if score == 0 {
				continue
			}

			if minCP == -1 || score < minCP {
				minCP = score
			}
		}

		if minCP != -1 {
			response = append(response, LevelMinCP{
				LevelType:   firstRecord.LevelType,
				LevelNumber: firstRecord.LevelNumber,
				LevelMode:   firstRecord.LevelMode,
				MinCP:       minCP,
			})
		}
	}

	sort.Slice(response, func(i, j int) bool {
		if response[i].LevelType != response[j].LevelType {
			return response[i].LevelType < response[j].LevelType
		}

		parseLevel := func(s string) (int, string) {
			parts := strings.Split(s, "_")
			val, _ := strconv.Atoi(parts[0])
			suffix := ""
			if len(parts) > 1 {
				suffix = parts[1]
			}
			return val, suffix
		}

		valI, suffI := parseLevel(response[i].LevelNumber)
		valJ, suffJ := parseLevel(response[j].LevelNumber)

		if valI != valJ {
			return valI < valJ
		}
		return suffI < suffJ
	})

	c.JSON(http.StatusOK, response)
}
