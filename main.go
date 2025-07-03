package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"lysk-battle-record/internal"
)

const (
	spreadsheetID     = "1-ORnXBnav4JVtP673Oio5sNdVpk0taUSzG3kWZqhIuY"
	orbitSheetName    = "轨道"
	championSheetName = "锦标赛"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	orbitGoogleSheetClient := internal.NewGoogleSheetClient(spreadsheetID, orbitSheetName)
	orbitRecordStore := internal.NewInMemoryRecordStore(orbitGoogleSheetClient)

	championshipsGoogleSheetClient := internal.NewGoogleSheetClient(spreadsheetID, championSheetName)
	championshipsRecordStore := internal.NewInMemoryRecordStore(championshipsGoogleSheetClient)

	server := internal.InitLyskServer(
		orbitRecordStore,
		orbitGoogleSheetClient,
		championshipsRecordStore,
		championshipsGoogleSheetClient,
		internal.NewAuthenticator(),
	)

	r := gin.Default()

	//allowOrigins := []string{"https://uygnim.com"}
	//if isLocal() {
	//	allowOrigins = []string{"*"}
	//}
	allowOrigins := []string{"*"}

	r.Use(cors.New(cors.Config{
		AllowOrigins:     allowOrigins,
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           1 * time.Hour,
	}))

	r.GET("/ping", server.Ping)
	r.POST("/login", server.Login)

	authRequired := r.Group("/")
	authRequired.Use(server.AuthMiddleware())
	{
		authRequired.POST("/orbit-record", server.ProcessOrbitRecord)
		authRequired.POST("/championships-record", server.ProcessChampionshipsRecord)
	}

	r.GET("/orbit-records", server.GetOrbitRecords)
	r.GET("/championships-records", server.GetChampionshipsRecords)

	r.GET("/latest-orbit-records", server.GetLatestOrbitRecords)
	r.GET("/latest-championships-records", server.GetLatestChampionshipsRecords)

	r.Run(":8080")
}

func isLocal() bool {
	if _, err := os.ReadFile("credentials.json"); err == nil {
		return true
	}

	return false
}
