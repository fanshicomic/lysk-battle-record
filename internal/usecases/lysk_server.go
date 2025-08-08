package usecases

import (
	"lysk-battle-record/internal/datastores"
	"lysk-battle-record/internal/models"
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

	Lottery *pkg.Lottery
}

func (s *LyskServer) populateNicknameForRecords(records []models.Record) {
	for i, record := range records {
		if record.UserID != "" {
			if user, ok := s.userStore.Get(record.UserID); ok {
				records[i].Nickname = user.Nickname
			}
		}
	}
}
