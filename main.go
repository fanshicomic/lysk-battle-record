package main

import (
	"lysk-battle-record/internal/datastores"
	"lysk-battle-record/internal/sheet_clients"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"lysk-battle-record/internal"
)

const (
	spreadsheetID     = "1-ORnXBnav4JVtP673Oio5sNdVpk0taUSzG3kWZqhIuY"
	orbitSheetName    = "轨道测试"
	championSheetName = "锦标赛"
	userSheetName     = "用户"
)

func main() {
	orbitGoogleSheetClient := sheet_clients.NewRecordSheetClient(spreadsheetID, orbitSheetName)
	orbitRecordStore := datastores.NewInMemoryRecordStore(orbitGoogleSheetClient)

	championshipsGoogleSheetClient := sheet_clients.NewRecordSheetClient(spreadsheetID, championSheetName)
	championshipsRecordStore := datastores.NewInMemoryRecordStore(championshipsGoogleSheetClient)

	userGoogleSheetClient := sheet_clients.NewUserSheetClient(spreadsheetID, userSheetName)
	userStore := datastores.NewInMemoryUserStore(userGoogleSheetClient)

	server := internal.InitLyskServer(
		orbitRecordStore,
		orbitGoogleSheetClient,
		championshipsRecordStore,
		championshipsGoogleSheetClient,
		userStore,
		userGoogleSheetClient,
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
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
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
		authRequired.PUT("/orbit-record/:id", server.UpdateOrbitRecord)
		authRequired.DELETE("/orbit-record/:id", server.DeleteOrbitRecord)
		authRequired.POST("/championships-record", server.ProcessChampionshipsRecord)
		authRequired.PUT("/championships-record/:id", server.UpdateChampionshipsRecord)
		authRequired.DELETE("/championships-record/:id", server.DeleteChampionshipsRecord)
		authRequired.GET("/my-orbit-record", server.GetMyOrbitRecords)
		authRequired.GET("/all-my-orbit-records", server.GetAllMyOrbitRecords)

		authRequired.POST("/user", server.CreateUser)
		authRequired.GET("/user", server.GetUser)
		authRequired.PUT("/user", server.UpdateUser)
	}

	r.GET("/orbit-records", server.GetOrbitRecords)
	r.GET("/championships-records", server.GetChampionshipsRecords)

	r.GET("/latest-orbit-records", server.GetLatestOrbitRecords)
	r.GET("/latest-championships-records", server.GetLatestChampionshipsRecords)

	r.GET("/ranking", server.GetRanking)

	r.Run(":8080")
}

func isLocal() bool {
	if _, err := os.ReadFile("credentials.json"); err == nil {
		return true
	}

	return false
}
