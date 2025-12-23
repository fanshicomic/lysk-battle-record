package usecases

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"lysk-battle-record/internal/datastores"
	"lysk-battle-record/internal/models"
	"github.com/gin-gonic/gin"
)

type MockRecordStore struct {}

func (m *MockRecordStore) GetAll() []models.Record { return nil }
func (m *MockRecordStore) Get(id string) (models.Record, bool) { return models.Record{}, false }
func (m *MockRecordStore) Query(opt datastores.QueryOptions) datastores.QueryResult { return datastores.QueryResult{} }
func (m *MockRecordStore) Insert(record models.Record) {}
func (m *MockRecordStore) Update(record models.Record) error { return nil }
func (m *MockRecordStore) Delete(record models.Record) error { return nil }
func (m *MockRecordStore) PrepareInsert(record models.Record) error { return nil }
func (m *MockRecordStore) IsDuplicate(record models.Record) bool { return false }
func (m *MockRecordStore) GetRanking(userId string) []models.RankingItem { return nil }
func (m *MockRecordStore) EvaluateRecord(record models.Record) string { return "" }
func (m *MockRecordStore) GetLevelRecords(record models.Record) []models.Record { return nil }
func (m *MockRecordStore) GetCompanionCounts() map[string]int { return nil }
func (m *MockRecordStore) GetPartnerLevelCounts() map[string]int { return nil }

func (m *MockRecordStore) GetAllLevelRecords() map[string][]models.Record {
	return map[string][]models.Record{
		"Light-10-Stable": {
			{LevelType: "Light", LevelNumber: "10", LevelMode: "Stable", CombatPower: models.CombatPower{BuffedScore: "1000"}},
		},
		"Light-2-Stable": {
			{LevelType: "Light", LevelNumber: "2", LevelMode: "Stable", CombatPower: models.CombatPower{BuffedScore: "2000"}},
		},
		"Fire-20-Stable": {
			{LevelType: "Fire", LevelNumber: "20", LevelMode: "Stable", CombatPower: models.CombatPower{BuffedScore: "3000"}},
		},
	}
}

func TestGetMinCombatPower(t *testing.T) {
	gin.SetMode(gin.TestMode)
	store := &MockRecordStore{}
	server := &LyskServer{
		orbitRecordStore: store,
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	server.GetMinCombatPower(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response []LevelMinCP
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(response) != 3 {
		t.Errorf("Expected 3 response items, got %d", len(response))
	}

	// Verify Sorting: Fire comes before Light (alphabetical), then Light 2 comes before Light 10 (numerical)
	
	// Index 0: Fire 20
	if response[0].LevelType != "Fire" || response[0].LevelNumber != "20" {
		t.Errorf("Expected first item to be Fire 20, got %s %s", response[0].LevelType, response[0].LevelNumber)
	}

	// Index 1: Light 2
	if response[1].LevelType != "Light" || response[1].LevelNumber != "2" {
		t.Errorf("Expected second item to be Light 2, got %s %s", response[1].LevelType, response[1].LevelNumber)
	}

	// Index 2: Light 10
	if response[2].LevelType != "Light" || response[2].LevelNumber != "10" {
		t.Errorf("Expected third item to be Light 10, got %s %s", response[2].LevelType, response[2].LevelNumber)
	}
}
