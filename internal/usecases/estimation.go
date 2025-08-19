package usecases

import (
	"fmt"
	"lysk-battle-record/internal/models"
	"strconv"
)

type CombatPowerEstimator interface {
	EstimateCombatPower(record models.Record) models.CombatPower
}

type combatPowerEstimator struct{}

func NewCombatPowerEstimator() CombatPowerEstimator {
	return &combatPowerEstimator{}
}

func (e *combatPowerEstimator) EstimateCombatPower(record models.Record) models.CombatPower {
	var matchingBuff, championshipsBuff float64 = 1, 1
	matchingBuff, err := strconv.ParseFloat(record.MatchingBuff, 64)
	if err == nil {
		matchingBuff = 1 + matchingBuff/100
	}
	championshipsBuff, err = strconv.ParseFloat(record.Buff, 64)
	if err == nil {
		championshipsBuff = 1 + championshipsBuff/100
	}

	score := 0

	//stats := record.ToStats()
	//partnerFlow := buildPartnerFlow(record.Partner, record.SetCard, record.Stage, record.Weapon)

	bufferedScore := int(float64(score) * matchingBuff * championshipsBuff)
	return models.CombatPower{
		Score:         fmt.Sprintf("%d", score),
		BufferedScore: fmt.Sprintf("%d", bufferedScore),
	}
}

func buildPartnerFlow(partner string, setCard string, stage string, weapon string) models.PartnerFlow {
	// use partner to first identify the partner's skill set
	// check set card, if matches partner, then stage buffer is enhanced buffer
	return models.PartnerFlow{}
}
