package internal

import (
	"log"
	"sync"
	"time"
)

type RecordStore interface {
	GetAll() []Record
	Query(opt QueryOptions) QueryResult
}

type InMemoryRecordStore struct {
	mu          sync.RWMutex
	records     []Record
	sheetClient GoogleSheetClient
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
		sheetClient: sheetClient,
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
	// 从 Google Sheets 获取数据
	data, err := s.sheetClient.FetchAllSheetData()
	if err != nil {
		log.Printf("failed to refresh cache: %v", err)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.records = data
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
