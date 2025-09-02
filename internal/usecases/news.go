package usecases

import (
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"lysk-battle-record/internal/models"
	"lysk-battle-record/internal/utils"
)

type LevelInfo struct {
	Level string `json:"level"`
	Count int    `json:"count"`
}

type TopRecord struct {
	Level string `json:"level"`
	CP    string `json:"cp"`
}

type News struct {
	OrbitRecordCompanionCounts         map[string]int `json:"orbit_record_companion_counts,omitempty"`
	ChampionshipsRecordCompanionCounts map[string]int `json:"championships_record_companion_counts,omitempty"`
	OrbitRecordPartnerCounts           map[string]int `json:"orbit_record_partner_counts,omitempty"`
	ChampionshipsRecordPartnerCounts   map[string]int `json:"championships_record_partner_counts,omitempty"`
	OrbitPartnerLevelCounts            map[string]int `json:"orbit_partner_level_counts,omitempty"`
	ChampionshipsPartnerLevelCounts    map[string]int `json:"championships_partner_level_counts,omitempty"`
	OrbitTopMostRecordsLevels          []LevelInfo    `json:"top_most_records_levels,omitempty"`
	OrbitLevelCounts                   int            `json:"orbit_level_counts,omitempty"`
	ChampionshipsLevelCounts           int            `json:"championships_level_counts,omitempty"`
	OrbitTopCPRecords                  []TopRecord    `json:"orbit_top_cp_records,omitempty"`
	ChampionshipsTopCPRecords          []TopRecord    `json:"championships_top_cp_records,omitempty"`
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

func (s *LyskServer) getUserOrbitRecord(userId string) []models.Record {
	allRecords := s.orbitRecordStore.GetAll()
	var userRecords []models.Record
	for _, record := range allRecords {
		if record.UserID == userId {
			userRecords = append(userRecords, record)
		}
	}

	return userRecords
}

func (s *LyskServer) getUserChampionshipsRecord(userId string) []models.Record {
	allRecords := s.championshipsRecordStore.GetAll()
	var userRecords []models.Record
	for _, record := range allRecords {
		if record.UserID == userId {
			userRecords = append(userRecords, record)
		}
	}

	return userRecords
}

func (s *LyskServer) getUserCompanionCounts(records []models.Record) map[string]int {
	companionCounts := make(map[string]int)

	for _, record := range records {
		if record.Companion != "" {
			companionCounts[record.Companion]++
		}
	}

	return companionCounts
}

func (s *LyskServer) getTopCPRecords(records []models.Record) []TopRecord {
	var topRecords []TopRecord

	for _, record := range records {
		if record.CombatPower.BuffedScore != "" && record.CombatPower.BuffedScore != models.NoData {
			level := formatLevelName(record.GenerateLevelKey())
			topRecords = append(topRecords, TopRecord{
				Level: level,
				CP:    record.CombatPower.BuffedScore,
			})
		}
	}

	sort.Slice(topRecords, func(i, j int) bool {
		cpI, errI := strconv.Atoi(topRecords[i].CP)
		cpJ, errJ := strconv.Atoi(topRecords[j].CP)

		if errI != nil || errJ != nil {
			return topRecords[i].CP > topRecords[j].CP
		}

		return cpI > cpJ
	})

	if len(topRecords) > 3 {
		return topRecords[:3]
	}

	return topRecords
}

func (s *LyskServer) GetNews(c *gin.Context) {
	news := News{
		OrbitRecordCompanionCounts:         s.GetOrbitRecordCompanionCounts(),
		ChampionshipsRecordCompanionCounts: s.GetChampionshipsRecordCompanionCounts(),
		OrbitRecordPartnerCounts:           s.GetOrbitRecordPartnerCounts(),
		ChampionshipsRecordPartnerCounts:   s.GetChampionshipsRecordPartnerCounts(),
		//OrbitPartnerLevelCounts:            s.GetOrbitPartnerLevelCounts(),
		//ChampionshipsPartnerLevelCounts:    s.GetChampionshipsPartnerLevelCounts(),
		OrbitTopMostRecordsLevels: s.GetOrbitTopMostRecordsLevels(),
		OrbitLevelCounts:          s.GetOrbitLevelCounts(),
	}

	c.JSON(http.StatusOK, news)
}

func (s *LyskServer) GetUserNews(c *gin.Context) {
	userId, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录或无效的用户"})
		return
	}

	news := News{}
	orbitRecords := s.getUserOrbitRecord(userId.(string))
	championshipsRecords := s.getUserChampionshipsRecord(userId.(string))

	news.OrbitLevelCounts = len(orbitRecords)
	news.ChampionshipsLevelCounts = len(championshipsRecords)

	news.OrbitRecordCompanionCounts = s.getUserCompanionCounts(orbitRecords)
	news.ChampionshipsRecordCompanionCounts = s.getUserCompanionCounts(championshipsRecords)
	news.OrbitRecordPartnerCounts = s.convertCompanionCountsToPartnerCounts(news.OrbitRecordCompanionCounts)
	news.ChampionshipsRecordPartnerCounts = s.convertCompanionCountsToPartnerCounts(news.ChampionshipsRecordCompanionCounts)

	news.OrbitTopCPRecords = s.getTopCPRecords(orbitRecords)
	news.ChampionshipsTopCPRecords = s.getTopCPRecords(championshipsRecords)

	c.JSON(http.StatusOK, news)
}

func formatLevelName(level string) string {
	// orbit
	if strings.Contains(level, "稳定") {
		level = strings.ReplaceAll(level, "-", " ")
		level = strings.ReplaceAll(level, "_", " ")
		if strings.Contains(level, "开放") {
			return level
		}
		level = strings.ReplaceAll(level, " 稳定", "")
		return level
	}

	// championships
	level = strings.ReplaceAll(level, "-A4", " A4")
	level = strings.ReplaceAll(level, "-B4", " B4")
	level = strings.ReplaceAll(level, "-C4", " C4")

	return level
}
