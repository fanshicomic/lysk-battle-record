package usecases

import (
	"fmt"
	"testing"

	"lysk-battle-record/internal/models"
)

func TestEstimateCombatPower(t *testing.T) {
	record := models.Record{
		Attack:       "5208",
		HP:           "87192",
		Defense:      "1944",
		Matching:     "逆",
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
		Stage:        "IV",
		Weapon:       "专武",
		Buff:         ""}

	estimator := NewCombatPowerEstimator()
	combatPower := estimator.EstimateCombatPower(record)
	fmt.Println("combatPower:", combatPower)
}
