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
}

type InMemoryRecordStore struct {
	mu             sync.RWMutex
	records        []models.Record
	recordsHash    map[string]bool
	ingestPoolHash map[string]bool
	sheetClient    sheet_clients.RecordSheetClient
	cpEstimator    estimator.CombatPowerEstimator
	ranking        []models.RankingItem
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
		sheetClient:    sheetClient,
		cpEstimator:    cpEstimator,
		ingestPoolHash: make(map[string]bool),
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

	s.mu.Lock()
	s.records = data
	s.ranking = ranking
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
	defer s.mu.RUnlock()
	if opt.Limit <= 0 {
		opt.Limit = 10
	}

	res := append(models.Records{}, s.records...)
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

	s.mu.Lock()
	defer s.mu.Unlock()
	record.CombatPower = s.cpEstimator.EstimateCombatPower(record)
	s.records = append(s.records, record)
	delete(s.ingestPoolHash, record.GetHash())
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
	defer s.mu.RUnlock()

	for _, r := range s.records {
		if r.Id == id {
			return r, true
		}
	}
	return models.Record{}, false
}

func (s *InMemoryRecordStore) Update(record models.Record) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, r := range s.records {
		if r.Id == record.Id {
			if r.Deleted {
				return errors.New("cannot update a deleted record")
			}
			record.CombatPower = s.cpEstimator.EstimateCombatPower(record)
			s.records[i] = record
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

func (s *InMemoryRecordStore) EvaluateRecord(record models.Record) string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	sameLevelRecords := []models.Record{}

	// Check if this is a championships record
	isChampionships := record.LevelType == "A4" || record.LevelType == "B4" || record.LevelType == "C4"
	var start, end time.Time
	if isChampionships {
		start, end = utils.GetCurrentChampionshipsRound()
	}

	for _, r := range s.records {
		if !r.Deleted &&
			r.LevelType == record.LevelType &&
			r.LevelNumber == record.LevelNumber &&
			r.LevelMode == record.LevelMode &&
			r.CombatPower.BuffedScore != "" &&
			r.CombatPower.BuffedScore != models.NoData {

			// For championships records, filter by current round time
			if isChampionships {
				recordTime, err := time.Parse(time.RFC3339, r.Time)
				if err != nil || recordTime.Before(start) || recordTime.After(end) {
					continue
				}
			}

			sameLevelRecords = append(sameLevelRecords, r)
		}
	}

	if len(sameLevelRecords) < 5 {
		return "标准"
	}

	buffedScores := []int{}
	for _, r := range sameLevelRecords {
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
