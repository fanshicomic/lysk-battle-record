package internal

import (
	"math/rand"
	"time"
)

type Prize struct {
	Name   string
	Chance float64 // 概率，0~1
	Count  int     // 剩余次数，-1表示无限
	Total  int
	IsBig  bool // 是否为大奖
}

type Lottery struct {
	BigPrizes   []*Prize
	SmallPrizes []*Prize
}

func NewLottery(bigPrizes, smallPrizes []*Prize) *Lottery {
	return &Lottery{
		BigPrizes:   bigPrizes,
		SmallPrizes: smallPrizes,
	}
}

// Draw 返回中奖奖项名称，未中奖返回空字符串
func (l *Lottery) Draw() string {
	rand.Seed(time.Now().UnixNano())

	// 先抽大奖
	r := rand.Float64()
	for _, p := range l.BigPrizes {
		if p.Count != 0 {
			if r < p.Chance {
				if p.Count > 0 {
					p.Count--
				}
				return p.Name
			}
		}
	}

	// 抽小奖
	for _, p := range l.SmallPrizes {
		if p.Count != 0 {
			if r < p.Chance {
				return p.Name
			}
		}
	}

	return ""
}

func (l *Lottery) GetAllPrizes() []Prize {
	prizes := make([]Prize, 0, len(l.BigPrizes)+len(l.SmallPrizes))
	for _, p := range l.BigPrizes {
		prizes = append(prizes, *p)
	}
	for _, p := range l.SmallPrizes {
		prizes = append(prizes, *p)
	}
	return prizes
}
