package usecases

import (
	"lysk-battle-record/internal/datastores"
	"lysk-battle-record/internal/pkg"
	"lysk-battle-record/internal/sheet_clients"
)

func InitLyskServer(orbitRecordStore datastores.RecordStore, orbitSheetClient sheet_clients.RecordSheetClient,
	championshipsRecordStore datastores.RecordStore, championshipsSheetClient sheet_clients.RecordSheetClient,
	userStore datastores.UserStore, userSheetClient sheet_clients.UserSheetClient, auth *pkg.Authenticator) *LyskServer {

	return &LyskServer{
		orbitRecordStore:         orbitRecordStore,
		orbitSheetClient:         orbitSheetClient,
		championshipsRecordStore: championshipsRecordStore,
		championshipsSheetClient: championshipsSheetClient,
		userStore:                userStore,
		userSheetClient:          userSheetClient,
		auth:                     auth,
	}
}

type LyskServer struct {
	orbitRecordStore         datastores.RecordStore
	orbitSheetClient         sheet_clients.RecordSheetClient
	championshipsRecordStore datastores.RecordStore
	championshipsSheetClient sheet_clients.RecordSheetClient
	userStore                datastores.UserStore
	userSheetClient          sheet_clients.UserSheetClient
	auth                     *pkg.Authenticator
}
