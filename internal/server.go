package internal

import (
	"github.com/gin-gonic/gin"
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
