package internal

import (
	"errors"
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
}

type InMemoryRecordStore struct {
	mu             sync.RWMutex
	records        []Record
	recordsHash    map[string]bool
	ingestPoolHash map[string]bool
	sheetClient    GoogleSheetClient
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
		logrus.Errorf("failed to refresh cache: %v", err)
		return
	}

	s.mu.Lock()
	s.records = data
	s.mu.Unlock()

	s.recordsHash = map[string]bool{}
	for _, record := range data {
		s.ingestHash(record)
	}
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

func getFilters(key, value string) func(Record) bool {
	return func(r Record) bool {
		switch key {
		case "关卡":
			return r.LevelType == value
		case "关数":
			return r.LevelNumber == value
		case "搭档身份":
			return r.Partner == value
		case "日卡":
			return r.SetCard == value
		case "阶数":
			return r.Stage == value
		case "武器":
			return r.Weapon == value
		default:
			return true
		}
	}
}
