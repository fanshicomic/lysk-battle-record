package usecases

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"lysk-battle-record/internal/utils"
)

type News struct {
	OrbitRecordCompanionCounts         map[string]int `json:"orbit_record_companion_counts"`
	ChampionshipsRecordCompanionCounts map[string]int `json:"championships_record_companion_counts"`
	OrbitRecordPartnerCounts           map[string]int `json:"orbit_record_partner_counts"`
	ChampionshipsRecordPartnerCounts   map[string]int `json:"championships_record_partner_counts"`
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

func (s *LyskServer) GetNews(c *gin.Context) {
	news := News{
		OrbitRecordCompanionCounts:         s.GetOrbitRecordCompanionCounts(),
		ChampionshipsRecordCompanionCounts: s.GetChampionshipsRecordCompanionCounts(),
		OrbitRecordPartnerCounts:           s.GetOrbitRecordPartnerCounts(),
		ChampionshipsRecordPartnerCounts:   s.GetChampionshipsRecordPartnerCounts(),
	}

	c.JSON(http.StatusOK, news)
}
