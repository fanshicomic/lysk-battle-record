package estimator

import (
	"fmt"
	"testing"

	"lysk-battle-record/internal/models"
)

func TestEstimateCombatPower(t *testing.T) {
	record := models.Record{
		Attack: "11200",
		//Attack:       "7900",
		HP:           "210000",
		Defense:      "5045",
		Matching:     "顺",
		MatchingBuff: "20",
		CritRate:     "65",
		CritDmg:      "300",
		//EnergyRegen:  "48",
		WeakenBoost: "90",
		OathBoost:   "40",
		OathRegen:   "",
		TotalLevel:  "300",
		//Companion:    "逐光骑士",
		//SetCard:      "逐光",
		//Companion: "永恒先知",
		//SetCard:   "永恒",
		//Companion: "深海潜行者",
		//SetCard:   "深海",
		//Companion: "潮汐之神",
		//SetCard: "神殿",
		//Companion: "光猎",
		//SetCard:   "末夜",
		//Companion: "九黎司命",
		//SetCard:   "拥雪",
		//Companion: "无尽掠夺者",
		//SetCard: "掠心",
		//Companion: "深渊主宰",
		//SetCard:   "深渊",
		//Companion: "终极兵器X-02",
		//SetCard:   "寂路",
		//Companion: "利莫里亚海神",
		//SetCard:   "雾海",
		Companion: "暗蚀国王",
		SetCard:   "夜誓",
		//Companion: "测试",
		Stage: "IV",
		//Weapon: "重剑",
		Weapon: "专武",
		Buff:   ""}

	estimator := NewCombatPowerEstimator()
	combatPower := estimator.EstimateCombatPower(record)
	fmt.Println("combatPower:", combatPower)
}

func TestEstimateCombatPowerReal(t *testing.T) {
	record := models.Record{
		Attack:       "6483",
		HP:           "148694",
		Defense:      "2830",
		Matching:     "顺",
		MatchingBuff: "25",
		CritRate:     "24",
		CritDmg:      "200.9",
		EnergyRegen:  "24",
		WeakenBoost:  "51.1",
		//OathBoost:    "1.8",
		OathRegen:  "",
		TotalLevel: "381",
		//Companion:      "逐光骑士",
		//SetCard:      "逐光",
		//Companion: "永恒先知",
		//SetCard: "永恒",
		//Companion: "潮汐之神",
		//SetCard: "神殿",
		//Companion: "光猎",
		//SetCard: "末夜",
		//Companion: "九黎司命",
		//SetCard: "拥雪",
		Companion: "无尽掠夺者",
		SetCard:   "掠心",
		Stage:     "IV",
		Weapon:    "重剑",
		//Weapon: "专武",
		Buff: "40"}

	estimator := NewCombatPowerEstimator()
	combatPower := estimator.EstimateCombatPower(record)
	fmt.Println("combatPower:", combatPower)
}
