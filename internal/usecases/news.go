package usecases

import (
	"github.com/gin-gonic/gin"
	"lysk-battle-record/internal/utils"
	"net/http"
	"sort"
	"strings"
)

type LevelInfo struct {
	Level string `json:"level"`
	Count int    `json:"count"`
}

type News struct {
	OrbitRecordCompanionCounts         map[string]int `json:"orbit_record_companion_counts"`
	ChampionshipsRecordCompanionCounts map[string]int `json:"championships_record_companion_counts"`
	OrbitRecordPartnerCounts           map[string]int `json:"orbit_record_partner_counts"`
	ChampionshipsRecordPartnerCounts   map[string]int `json:"championships_record_partner_counts"`
	OrbitPartnerLevelCounts            map[string]int `json:"orbit_partner_level_counts"`
	ChampionshipsPartnerLevelCounts    map[string]int `json:"championships_partner_level_counts"`
	OrbitTopMostRecordsLevels          []LevelInfo    `json:"top_most_records_levels"`
	OrbitLevelCounts                   int            `json:"orbit_level_counts"`
}

func (s *LyskServer) GetOrbitRecordCompanionCounts() map[string]int {
	return s.orbitRecordStore.GetCompanionCounts()
}

func (s *LyskServer) GetChampionshipsRecordCompanionCounts() map[string]int {
	return s.championshipsRecordStore.GetCompanionCounts()
}

func (s *LyskServer) GetOrbitRecordPartnerCounts() map[string]int {
	companionCounts := s.orbitRecordStore.GetCompanionCounts()
	return s.convertCompanionCountsToPartnerCounts(companionCounts)
}

func (s *LyskServer) GetChampionshipsRecordPartnerCounts() map[string]int {
	companionCounts := s.championshipsRecordStore.GetCompanionCounts()
	return s.convertCompanionCountsToPartnerCounts(companionCounts)
}

func (s *LyskServer) GetOrbitPartnerLevelCounts() map[string]int {
	return s.orbitRecordStore.GetPartnerLevelCounts()
}

func (s *LyskServer) GetChampionshipsPartnerLevelCounts() map[string]int {
	return s.championshipsRecordStore.GetPartnerLevelCounts()
}

func (s *LyskServer) convertCompanionCountsToPartnerCounts(companionCounts map[string]int) map[string]int {
	partnerCounts := make(map[string]int)
	partnerCompanionMap := utils.GetPartnerCompanionMap()

	for partner, companions := range partnerCompanionMap {
		totalCount := 0
		for _, companion := range companions {
			if count, exists := companionCounts[companion]; exists {
				totalCount += count
			}
		}
		if totalCount > 0 {
			partnerCounts[partner] = totalCount
		}
	}

	return partnerCounts
}

func (s *LyskServer) GetOrbitTopMostRecordsLevels() []LevelInfo {
	var levels []LevelInfo

	levelRecordsMap := s.orbitRecordStore.GetAllLevelRecords()
	for level, records := range levelRecordsMap {
		levels = append(levels, LevelInfo{
			Level: level,
			Count: len(records),
		})
	}

	sort.Slice(levels, func(i, j int) bool {
		return levels[i].Count > levels[j].Count
	})

	for i := range levels {
		levels[i].Level = formatLevelName(levels[i].Level)
	}

	if len(levels) > 3 {
		return levels[:3]
	}

	return levels
}

func (s *LyskServer) GetOrbitLevelCounts() int {
	levelRecordsMap := s.orbitRecordStore.GetAllLevelRecords()
	return len(levelRecordsMap)
}

func (s *LyskServer) GetNews(c *gin.Context) {
	news := News{
		OrbitRecordCompanionCounts:         s.GetOrbitRecordCompanionCounts(),
		ChampionshipsRecordCompanionCounts: s.GetChampionshipsRecordCompanionCounts(),
		OrbitRecordPartnerCounts:           s.GetOrbitRecordPartnerCounts(),
		ChampionshipsRecordPartnerCounts:   s.GetChampionshipsRecordPartnerCounts(),
		OrbitPartnerLevelCounts:            s.GetOrbitPartnerLevelCounts(),
		ChampionshipsPartnerLevelCounts:    s.GetChampionshipsPartnerLevelCounts(),
		OrbitTopMostRecordsLevels:          s.GetOrbitTopMostRecordsLevels(),
		OrbitLevelCounts:                   s.GetOrbitLevelCounts(),
	}

	c.JSON(http.StatusOK, news)
}

func formatLevelName(level string) string {
	level = strings.ReplaceAll(level, "-", " ")
	level = strings.ReplaceAll(level, "_", " ")
	return level
}
