package usecases

import (
	"fmt"
	"testing"

	"lysk-battle-record/internal/models"
)

func TestEstimateCombatPower(t *testing.T) {
	record := models.Record{
		Attack:       "11208",
		HP:           "87192",
		Defense:      "4944",
		Matching:     "顺",
		MatchingBuff: "20",
		CritRate:     "13.3",
		CritDmg:      "203.7",
		EnergyRegen:  "48",
		WeakenBoost:  "39.7",
		OathBoost:    "1.8",
		OathRegen:    "",
		TotalLevel:   "300",
		Partner:      "逐光骑士",
		SetCard:      "逐光",
		//Partner: "永恒先知",
		//SetCard: "永恒",
		//Partner: "深海潜行者",
		//SetCard: "深海",
		Stage: "IV",
		//Weapon: "重剑",
		Weapon: "专武",
		Buff:   ""}

	estimator := NewCombatPowerEstimator()
	combatPower := estimator.EstimateCombatPower(record)
	fmt.Println("combatPower:", combatPower)
}
