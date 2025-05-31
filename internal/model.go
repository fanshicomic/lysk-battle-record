package internal

import (
	"sort"
	"time"
)

type Record struct {
	LevelType   string `json:"关卡"`
	LevelNumber string `json:"关数"`
	Attack      string `json:"攻击"`
	HP          string `json:"生命"`
	Defense     string `json:"防御"`
	Matching    string `json:"对谱"`
	CritRate    string `json:"暴击"`
	CritDmg     string `json:"暴伤"`
	EnergyRegen string `json:"加速回能"`
	WeakenBoost string `json:"虚弱增伤"`
	OathBoost   string `json:"誓约增伤"`
	OathRegen   string `json:"誓约回能"`
	Partner     string `json:"搭档身份"`
	SetCard     string `json:"日卡"`
	Stage       string `json:"阶数"`
	Weapon      string `json:"武器"`
	Time        string `json:"时间"` // 可额外解析为 time.Time
}
type Records []Record

func (r Records) filter(filterFunc func(Record) bool) []Record {
	result := Records{}
	for _, r := range r {
		if filterFunc(r) {
			result = append(result, r)
		}
	}
	return result
}

func (r Records) sortByTimeDesc() []Record {
	sort.Slice(r, func(i, j int) bool {
		ti, _ := time.Parse("2006-01-02T15:04:05Z", r[i].Time)
		tj, _ := time.Parse("2006-01-02T15:04:05Z", r[j].Time)
		return tj.Before(ti)
	})
	return r
}

func (r Records) pagination(offset, size int) []Record {
	total := len(r)
	if offset >= total {
		return []Record{}
	}
	end := offset + size
	if end > total {
		end = total
	}
	return r[offset:end]
}
