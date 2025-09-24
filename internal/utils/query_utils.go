package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	AllCompanion = "所有搭档"
	AllSetCard   = "所有日卡"
)

func GetCompanion(c *gin.Context) string {
	filteredCompanion := c.DefaultQuery("filteredCompanion", AllCompanion)
	if filteredCompanion == "" || filteredCompanion != "null" {
		return AllCompanion
	}
	return filteredCompanion
}

func GetSetCard(c *gin.Context) string {
	filteredSetCard := c.DefaultQuery("filteredSetCard", AllSetCard)
	if filteredSetCard == "" || filteredSetCard == "null" {
		return AllSetCard
	}
	return filteredSetCard
}

func BuildBasicFilters(c *gin.Context) map[string]string {
	filters := map[string]string{}

	filteredCompanion := GetCompanion(c)
	if filteredCompanion != AllCompanion {
		filters["搭档身份"] = filteredCompanion
	}

	filteredSetCard := GetSetCard(c)
	if filteredSetCard != AllSetCard {
		filters["日卡"] = filteredSetCard
	}

	return filters
}

func BuildOrbitFilters(c *gin.Context, shouldFilteredByCompanionOrSetCard bool) map[string]string {
	filters := map[string]string{}
	if shouldFilteredByCompanionOrSetCard {
		filters = BuildBasicFilters(c)
	}

	filters["关卡"] = c.Query("type")
	filters["关数"] = c.Query("level")

	mode := c.Query("mode")
	if mode == "" {
		mode = "稳定"
	}
	filters["模式"] = mode

	return filters
}

func BuildChampionshipFilters(c *gin.Context, shouldFilteredByCompanionOrSetCard bool) map[string]string {
	filters := map[string]string{}
	if shouldFilteredByCompanionOrSetCard {
		filters = BuildBasicFilters(c)
	}

	filters["关卡"] = c.Query("level")

	return filters
}

func GetOffset(c *gin.Context) int {
	offsetStr := c.DefaultQuery("offset", "0")
	offset, _ := strconv.Atoi(offsetStr)
	return offset
}
