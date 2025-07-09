package internal

import (
	"errors"
	"sort"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type RecordStore interface {
	GetAll() []Record
	Query(opt QueryOptions) QueryResult
	Insert(record Record)
	PrepareInsert(record Record) error
	IsDuplicate(record Record) bool
	GetRanking(userId string) []RankingItem
}

type InMemoryRecordStore struct {
	mu             sync.RWMutex
	records        []Record
	recordsHash    map[string]bool
	ingestPoolHash map[string]bool
	sheetClient    GoogleSheetClient
	ranking        []RankingItem
}

type QueryOptions struct {
	Filters map[string]string
	SortBy  string // 排序字段
	Desc    bool   // 是否降序
	Offset  int
	Limit   int
}

type QueryResult struct {
	Total   int      `json:"total"`
	Records []Record `json:"records"`
}

func NewInMemoryRecordStore(sheetClient GoogleSheetClient) *InMemoryRecordStore {
	store := &InMemoryRecordStore{
		sheetClient:    sheetClient,
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
	for _, record := range data {
		if len(record.UserID) > 0 && record.UserID != "<nil>" {
			contribution[record.UserID] += 1
		}
	}

	ranking := []RankingItem{}

	for userId, count := range contribution {
		ranking = append(ranking, RankingItem{
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

func (s *InMemoryRecordStore) ingestHash(record Record) {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := record.getHash()
	s.recordsHash[key] = true
}

func (s *InMemoryRecordStore) GetAll() []Record {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]Record(nil), s.records...)
}

func (s *InMemoryRecordStore) Query(opt QueryOptions) QueryResult {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if opt.Limit <= 0 {
		opt.Limit = 10
	}

	res := append(Records{}, s.records...)
	for k, v := range opt.Filters {
		filterFunc := getFilters(k, v)
		res = res.filter(filterFunc)
	}
	count := len(res)

	res = res.sortByTimeDesc()
	res = res.pagination(opt.Offset, opt.Limit)

	return QueryResult{
		Total:   count,
		Records: res,
	}
}

func (s *InMemoryRecordStore) Insert(record Record) {
	s.ingestHash(record)

	s.mu.Lock()
	defer s.mu.Unlock()
	s.records = append(s.records, record)
	delete(s.ingestPoolHash, record.getHash())
}

func (s *InMemoryRecordStore) PrepareInsert(record Record) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := record.getHash()
	if s.ingestPoolHash[key] {
		return errors.New("记录已在上传准备中")
	}
	s.ingestPoolHash[key] = true
	return nil
}

func (s *InMemoryRecordStore) IsDuplicate(record Record) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.recordsHash[record.getHash()] || s.ingestPoolHash[record.getHash()]
}

func (s *InMemoryRecordStore) GetRanking(userId string) []RankingItem {
	s.mu.RLock()
	defer s.mu.RUnlock()

	limit := 10

	processedRanking := []RankingItem{}
	rank := 1
	for i, item := range s.ranking {
		if i > 0 && item.Contribution < s.ranking[i-1].Contribution {
			rank = i + 1
		}
		processedRanking = append(processedRanking, RankingItem{
			OpenID:       item.OpenID,
			Contribution: item.Contribution,
			Rank:         rank,
		})
	}

	result := []RankingItem{}
	userInResult := false
	var userRank *RankingItem

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

func getFilters(key, value string) func(Record) bool {
	return func(r Record) bool {
		switch key {
		case "关卡":
			return r.LevelType == value
		case "关数":
			return r.LevelNumber == value
		case "模式":
			return r.LevelMode == value
		case "搭档身份":
			return r.Partner == value
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
