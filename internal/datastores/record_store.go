package datastores

import (
	"errors"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"lysk-battle-record/internal/estimator"
	"lysk-battle-record/internal/models"
	"lysk-battle-record/internal/sheet_clients"
	"lysk-battle-record/internal/utils"
)

type RecordStore interface {
	GetAll() []models.Record
	Get(id string) (models.Record, bool)
	Query(opt QueryOptions) QueryResult
	Insert(record models.Record)
	Update(record models.Record) error
	Delete(record models.Record) error
	PrepareInsert(record models.Record) error
	IsDuplicate(record models.Record) bool
	GetRanking(userId string) []models.RankingItem
	EvaluateRecord(record models.Record) string

	GetAllLevelRecords() map[string][]models.Record
	GetLevelRecords(record models.Record) []models.Record
	GetCompanionCounts() map[string]int
	GetPartnerLevelCounts() map[string]int
}

type InMemoryRecordStore struct {
	mu              sync.RWMutex
	records         []models.Record
	recordsHash     map[string]bool
	ingestPoolHash  map[string]bool
	sheetClient     sheet_clients.RecordSheetClient
	cpEstimator     estimator.CombatPowerEstimator
	ranking         []models.RankingItem
	levelRecords    map[string][]models.Record // New field to store records by level key
	companionCounts map[string]int             // Cache companion counts
}

type QueryOptions struct {
	Filters   map[string]string
	SortBy    string // 排序字段
	Desc      bool   // 是否降序
	Offset    int
	Limit     int
	TimeStart time.Time
	TimeEnd   time.Time
}

type QueryResult struct {
	Total   int             `json:"total"`
	Records []models.Record `json:"records"`
}

func NewInMemoryRecordStore(sheetClient sheet_clients.RecordSheetClient, cpEstimator estimator.CombatPowerEstimator) *InMemoryRecordStore {
	store := &InMemoryRecordStore{
		sheetClient:     sheetClient,
		cpEstimator:     cpEstimator,
		ingestPoolHash:  make(map[string]bool),
		levelRecords:    make(map[string][]models.Record), // Initialize the levelRecords map
		companionCounts: make(map[string]int),             // Initialize the companionCounts map
	}
	go store.autoRefresh()
	return store
}

func (s *InMemoryRecordStore) autoRefresh() {
	for {
		s.refresh()
		time.Sleep(5 * time.Minute)
	}
}

func (s *InMemoryRecordStore) refresh() {
	data, err := s.sheetClient.FetchAllSheetData()
	if err != nil {
		logrus.Errorf("failed to refresh cache for sheet: %s with error %v", s.sheetClient.GetType(), err)
		return
	}

	contribution := map[string]int32{}
	for i, record := range data {
		data[i].CombatPower = s.cpEstimator.EstimateCombatPower(record)
		if len(record.UserID) > 0 && record.UserID != "<nil>" {
			contribution[record.UserID] += 1
		}
	}

	ranking := []models.RankingItem{}

	for userId, count := range contribution {
		ranking = append(ranking, models.RankingItem{
			OpenID:       userId,
			Contribution: count,
		})
	}

	sort.Slice(ranking, func(i, j int) bool {
		return ranking[i].Contribution > ranking[j].Contribution
	})

	// Organize records by level key for quick retrieval
	levelRecords := make(map[string][]models.Record)
	companionCounts := make(map[string]int)

	for _, record := range data {
		if !record.Deleted {
			levelKey := record.GenerateLevelKey()
			levelRecords[levelKey] = append(levelRecords[levelKey], record)

			// Count companions
			if record.Companion != "" {
				companionCounts[record.Companion]++
			}
		}
	}

	s.mu.Lock()
	s.records = data
	s.ranking = ranking
	s.levelRecords = levelRecords
	s.companionCounts = companionCounts
	s.mu.Unlock()

	s.recordsHash = map[string]bool{}
	for _, record := range data {
		s.ingestHash(record)
	}
	logrus.Infof("sheet %s refreshed %d records", s.sheetClient.GetType(), len(s.records))
}

func (s *InMemoryRecordStore) ingestHash(record models.Record) {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := record.GetHash()
	s.recordsHash[key] = true
}

func (s *InMemoryRecordStore) GetAll() []models.Record {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]models.Record(nil), s.records...)
}

func (s *InMemoryRecordStore) Query(opt QueryOptions) QueryResult {
	s.mu.RLock()
	res := append(models.Records{}, s.records...)
	s.mu.RUnlock()

	if opt.Limit <= 0 {
		opt.Limit = 10
	}

	for k, v := range opt.Filters {
		filterFunc := getFilters(k, v)
		res = res.Filter(filterFunc)
	}
	res = res.Filter(filterOutDeleted())

	if !opt.TimeStart.IsZero() && !opt.TimeEnd.IsZero() {
		res = res.Filter(func(r models.Record) bool {
			recordTime, err := time.Parse(time.RFC3339, r.Time)
			if err != nil {
				logrus.Errorf("failed to parse record time %s: %v", r.Time, err)
				return false
			}
			return recordTime.After(opt.TimeStart) && recordTime.Before(opt.TimeEnd)
		})
	}
	count := len(res)

	res = res.SortByTimeDesc()
	res = res.Pagination(opt.Offset, opt.Limit)
	res = s.populateEvaluation(res)

	return QueryResult{
		Total:   count,
		Records: res,
	}
}

func (s *InMemoryRecordStore) Insert(record models.Record) {
	s.ingestHash(record)
	record.CombatPower = s.cpEstimator.EstimateCombatPower(record)

	s.mu.Lock()
	defer s.mu.Unlock()
	s.records = append(s.records, record)
	delete(s.ingestPoolHash, record.GetHash())

	// Update the levelRecords map and companion counts
	if !record.Deleted {
		levelKey := record.GenerateLevelKey()
		s.levelRecords[levelKey] = append(s.levelRecords[levelKey], record)

		// Update companion counts
		if record.Companion != "" {
			s.companionCounts[record.Companion]++
		}
	}
}

func (s *InMemoryRecordStore) PrepareInsert(record models.Record) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := record.GetHash()
	if s.ingestPoolHash[key] {
		return errors.New("记录已在上传准备中")
	}
	s.ingestPoolHash[key] = true
	return nil
}

func (s *InMemoryRecordStore) IsDuplicate(record models.Record) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	hash := record.GetHash()
	return s.recordsHash[hash] || s.ingestPoolHash[hash]
}

func (s *InMemoryRecordStore) Get(id string) (models.Record, bool) {
	s.mu.RLock()
	records := s.records
	s.mu.RUnlock()

	for _, r := range records {
		if r.Id == id {
			return r, true
		}
	}
	return models.Record{}, false
}

func (s *InMemoryRecordStore) Update(record models.Record) error {
	record.CombatPower = s.cpEstimator.EstimateCombatPower(record)
	levelKey := record.GenerateLevelKey()

	s.mu.Lock()
	defer s.mu.Unlock()
	for i, r := range s.records {
		if r.Id == record.Id {
			if r.Deleted {
				return errors.New("cannot update a deleted record")
			}

			oldCompanion := r.Companion
			s.records[i] = record

			// Update in level bucket - find and replace by ID
			if levelRecords, exists := s.levelRecords[levelKey]; exists {
				for j, lr := range levelRecords {
					if lr.Id == record.Id {
						s.levelRecords[levelKey][j] = record
						break
					}
				}
			}

			// Update companion counts
			if oldCompanion != "" {
				s.companionCounts[oldCompanion]--
				if s.companionCounts[oldCompanion] <= 0 {
					delete(s.companionCounts, oldCompanion)
				}
			}
			if record.Companion != "" {
				s.companionCounts[record.Companion]++
			}
			break
		}
	}
	return nil
}

func (s *InMemoryRecordStore) Delete(record models.Record) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	hash := record.GetHash()
	s.recordsHash[hash] = false
	for i, r := range s.records {
		if r.Id == record.Id {
			s.records[i].Deleted = true

			// Remove from the levelRecords map by finding the correct record by ID
			levelKey := r.GenerateLevelKey()
			if levelRecords, exists := s.levelRecords[levelKey]; exists {
				for j, lr := range levelRecords {
					if lr.Id == record.Id {
						s.levelRecords[levelKey] = append(levelRecords[:j], levelRecords[j+1:]...)
						break
					}
				}
			}

			// Update companion counts
			if r.Companion != "" {
				s.companionCounts[r.Companion]--
				if s.companionCounts[r.Companion] <= 0 {
					delete(s.companionCounts, r.Companion)
				}
			}
			break
		}
	}

	return nil
}

func (s *InMemoryRecordStore) GetRanking(userId string) []models.RankingItem {
	s.mu.RLock()
	defer s.mu.RUnlock()

	limit := 10

	processedRanking := []models.RankingItem{}
	rank := 1
	for i, item := range s.ranking {
		if i > 0 && item.Contribution < s.ranking[i-1].Contribution {
			rank = i + 1
		}
		processedRanking = append(processedRanking, models.RankingItem{
			OpenID:       item.OpenID,
			Contribution: item.Contribution,
			Rank:         rank,
		})
	}

	result := []models.RankingItem{}
	userInResult := false
	var userRank *models.RankingItem

	for _, item := range processedRanking {
		if item.Rank <= limit {
			result = append(result, item)
			if item.OpenID == userId {
				userInResult = true
			}
		}
		if item.OpenID == userId {
			userRank = &item
		}
	}

	if !userInResult && userRank != nil {
		result = append(result, *userRank)
	}

	return result
}

func (s *InMemoryRecordStore) GetAllLevelRecords() map[string][]models.Record {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.levelRecords
}

func (s *InMemoryRecordStore) GetLevelRecords(record models.Record) []models.Record {
	// Create a temporary record to generate the level key
	tempRecord := models.Record{
		LevelType:   record.LevelType,
		LevelNumber: record.LevelNumber,
		LevelMode:   record.LevelMode,
		Time:        record.Time,
	}

	levelKey := s.generateLevelKey(tempRecord)

	s.mu.RLock()
	levelRecords, exists := s.levelRecords[levelKey]
	s.mu.RUnlock()

	if !exists {
		return []models.Record{}
	}

	// Filter out deleted records and apply championships time filtering if needed
	result := []models.Record{}
	isChampionships := record.LevelType == "A4" || record.LevelType == "B4" || record.LevelType == "C4"
	var start, end time.Time
	if isChampionships {
		recordTime, err := time.Parse(time.RFC3339, record.Time)
		if err != nil {
			recordTime = time.Now()
		}
		start, end = utils.GetChampionshipsRoundByTime(recordTime)
	}

	for _, r := range levelRecords {
		if r.Deleted {
			continue
		}

		if isChampionships {
			rTime, err := time.Parse(time.RFC3339, r.Time)
			if err != nil || rTime.Before(start) || rTime.After(end) {
				continue
			}
		}

		if record.Companion != "" && record.Companion != utils.AllCompanion && record.Companion != r.Companion {
			continue
		}

		if record.SetCard != "" && record.SetCard != utils.AllSetCard && record.SetCard != r.SetCard {
			continue
		}

		result = append(result, r)
	}

	return result
}

func (s *InMemoryRecordStore) EvaluateRecord(record models.Record) string {
	sameLevelRecords := s.GetLevelRecords(models.Record{
		LevelType:   record.LevelType,
		LevelNumber: record.LevelNumber,
		LevelMode:   record.LevelMode,
	})

	// Filter records with valid buffed scores
	validRecords := []models.Record{}
	for _, r := range sameLevelRecords {
		if r.CombatPower.BuffedScore != "" && r.CombatPower.BuffedScore != models.NoData {
			validRecords = append(validRecords, r)
		}
	}

	if len(validRecords) < 5 {
		return "标准"
	}

	buffedScores := []int{}
	for _, r := range validRecords {
		if score, err := strconv.Atoi(r.CombatPower.BuffedScore); err == nil {
			buffedScores = append(buffedScores, score)
		}
	}

	if len(buffedScores) < 5 {
		return "标准"
	}

	sort.Ints(buffedScores)

	recordBuffedScore, err := strconv.Atoi(record.CombatPower.BuffedScore)
	if err != nil {
		return "标准"
	}

	q1Index := len(buffedScores) / 4
	q3Index := (3 * len(buffedScores)) / 4
	q1 := buffedScores[q1Index]
	q3 := buffedScores[q3Index]

	if recordBuffedScore >= q3 {
		return "溢出"
	} else if recordBuffedScore <= q1 {
		return "极限"
	} else {
		return "标准"
	}
}

func (s *InMemoryRecordStore) populateEvaluation(records []models.Record) []models.Record {
	for i, r := range records {
		records[i].CombatPower.Evaluation = s.EvaluateRecord(r)
	}
	return records
}

func (s *InMemoryRecordStore) generateLevelKey(record models.Record) string {
	// For championships: round-leveltype
	if record.LevelType == "A4" || record.LevelType == "B4" || record.LevelType == "C4" {
		recordTime, err := time.Parse(time.RFC3339, record.Time)
		if err != nil {
			recordTime = time.Now()
		}
		start, _ := utils.GetChampionshipsRoundByTime(recordTime)
		roundKey := start.Format("2006-01-02") // Use start date as round identifier
		return roundKey + "-" + record.LevelType
	}
	// For orbit: leveltype-levelnumber-levelmode
	return record.LevelType + "-" + record.LevelNumber + "-" + record.LevelMode
}

func (s *InMemoryRecordStore) GetCompanionCounts() map[string]int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.companionCounts
}

func (s *InMemoryRecordStore) GetPartnerLevelCounts() map[string]int {
	s.mu.RLock()
	levelRecords := s.levelRecords
	s.mu.RUnlock()

	partnerCompanionMap := utils.GetPartnerCompanionMap()
	partnerLevelSets := make(map[string]map[string]bool)

	// Initialize partner level sets
	for partner := range partnerCompanionMap {
		partnerLevelSets[partner] = make(map[string]bool)
	}

	// Iterate through all level records to count unique levels per partner
	for levelKey, records := range levelRecords {
		if len(records) == 0 {
			continue
		}

		// Track which partners have records in this level
		partnersInLevel := make(map[string]bool)

		for _, record := range records {
			if record.Deleted {
				continue
			}

			// Find which partner this companion belongs to
			for partner, companions := range partnerCompanionMap {
				for _, companion := range companions {
					if record.Companion == companion {
						partnersInLevel[partner] = true
						break
					}
				}
			}
		}

		// Add this level to each partner that has records in it
		for partner := range partnersInLevel {
			partnerLevelSets[partner][levelKey] = true
		}
	}

	// Convert sets to counts
	result := make(map[string]int)
	for partner, levelSet := range partnerLevelSets {
		result[partner] = len(levelSet)
	}

	return result
}

func getFilters(key, value string) func(models.Record) bool {
	return func(r models.Record) bool {
		switch key {
		case "关卡":
			return r.LevelType == value
		case "关数":
			return r.LevelNumber == value
		case "模式":
			return r.LevelMode == value
		case "搭档身份":
			return r.Companion == value
		case "日卡":
			return r.SetCard == value
		case "阶数":
			return r.Stage == value
		case "武器":
			return r.Weapon == value
		case "用户ID":
			return r.UserID == value
		default:
			return true
		}
	}
}

func filterOutDeleted() func(models.Record) bool {
	return func(r models.Record) bool {
		return r.Deleted != true
	}
}
