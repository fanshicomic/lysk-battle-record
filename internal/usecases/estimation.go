package usecases

import (
	"fmt"
	"lysk-battle-record/internal/models"
	"lysk-battle-record/internal/partners"
	"lysk-battle-record/internal/set_cards"
)

type CombatPowerEstimator interface {
	EstimateCombatPower(record models.Record) models.CombatPower
}

type combatPowerEstimator struct{}

func NewCombatPowerEstimator() CombatPowerEstimator {
	return &combatPowerEstimator{}
}

func (e *combatPowerEstimator) EstimateCombatPower(record models.Record) models.CombatPower {
	stats := record.ToStats()
	partnerFlow := getPartnerFlow(stats)
	setCard := getSetCard(stats)
	setCardBuff := setCard.GetSetCardBuff()[stats.Stage]
	applySetCardBuff(&partnerFlow, setCardBuff)

	score := estimate(stats, partnerFlow)
	fmt.Println(partnerFlow)
	fmt.Println(score)
	return score
}

func getPartnerFlow(stats models.Stats) models.PartnerFlow {
	var partner partners.Partner
	switch stats.Partner {
	case "逐光骑士":
		partner = partners.PartnerLightSeeker{}
	default:
		partner = partners.PartnerLightSeeker{}
	}
	return partner.GetPartnerFlow(stats)
}

func getSetCard(stats models.Stats) set_cards.SetCard {
	var setCard set_cards.SetCard
	switch stats.SetCard {
	case "逐光":
		setCard = set_cards.LightSeeking{}
	default:
		setCard = set_cards.NoSet{}
	}
	return setCard
}

func applySetCardBuff(flow *models.PartnerFlow, setCardBuff models.StageBuff) {
	for periodIdx := range flow.Periods {
		period := &flow.Periods[periodIdx]
		for i, skill := range period.SkillSet.Skills {
			if allBuff, exists := setCardBuff.Buffs["所有"]; exists {
				skill.CritRate += allBuff.CritRate
				skill.CritDmg += allBuff.CritDmg
				skill.WeakenBoost += allBuff.WeakenBoost
				skill.DamageBoost += allBuff.DamageBoost
				if allBuff.CountBonus > 1 {
					skill.Count = int(float64(skill.Count) * allBuff.CountBonus)
				}
			}

			if skillBuff, exists := setCardBuff.Buffs[skill.Name]; exists {
				skill.CritRate += skillBuff.CritRate
				skill.CritDmg += skillBuff.CritDmg
				skill.WeakenBoost += skillBuff.WeakenBoost
				skill.DamageBoost += skillBuff.DamageBoost
				if skillBuff.CountBonus > 1 {
					skill.Count = int(float64(skill.Count) * skillBuff.CountBonus)
				}
			}

			period.SkillSet.Skills[i] = skill
		}
	}
}

func estimate(stats models.Stats, partnerFlow models.PartnerFlow) models.CombatPower {
	var total, weakenScore, nonWeakenScore float64 = 0, 0, 0
	for _, period := range partnerFlow.Periods {
		var score float64 = 0

		for _, skill := range period.SkillSet.Skills {
			rawSkillScore := skill.Base +
				(skill.HpRate/100)*float64(stats.HP) +
				(skill.AttackRate/100)*float64(stats.Attack) +
				(skill.DefenseRate/100)*float64(stats.Defense)

			// apply damage boost
			rawSkillScore *= 1 + skill.DamageBoost/100

			// consider level - defence relationship
			levelDefenseRatio := 1 + float64(stats.TotalLevel)/(float64(stats.TotalLevel)+300+(80*3+100)*(1-skill.DefenseRate/100))
			rawSkillScore *= levelDefenseRatio

			// consider non-weaken period
			nonWeakenSkillCount := (1 - partnerFlow.WeakenRate) * float64(skill.Count)
			if skill.Name == "誓约" { // 誓约只在虚弱期放，X-02誓约技会算作特殊技能
				nonWeakenSkillCount = 0
			}
			nonWeakenPeriodScore := rawSkillScore * nonWeakenSkillCount
			critRate := (stats.CritRate + skill.CritRate) / 100
			critDmg := (stats.CritDmg + skill.CritDmg) / 100
			nonWeakenPeriodScore = nonWeakenPeriodScore*(1-critRate) +
				nonWeakenPeriodScore*critRate*critDmg

			// consider weaken period
			weakenSkillCount := partnerFlow.WeakenRate * float64(skill.Count)
			if skill.Name == "共鸣" { // 共鸣技不会进入虚弱期
				weakenSkillCount = 0
			}
			weakenPeriodScore := rawSkillScore * weakenSkillCount
			weakenBoost := stats.WeakenBoost + skill.WeakenBoost
			if stats.Matching == "顺" {
				weakenBoost += 250
			} else {
				weakenBoost += 150
			}
			weakenPeriodScore = weakenPeriodScore * (1 + weakenBoost/100)

			score += nonWeakenPeriodScore + weakenPeriodScore
			nonWeakenScore += nonWeakenPeriodScore
			weakenScore += weakenPeriodScore
		}

		total += score
	}

	matchingBuff := 1 + stats.MatchingBuff/100
	championshipsBuff := 1 + stats.Buff/100

	total *= 1 + partnerFlow.Boost/100
	weakenScore *= 1 + partnerFlow.Boost/100
	nonWeakenScore *= 1 + partnerFlow.Boost/100
	buffedTotal := int(total * matchingBuff * championshipsBuff)

	return models.CombatPower{
		Score:          fmt.Sprintf("%d", int(total)),
		BuffedScore:    fmt.Sprintf("%d", buffedTotal),
		WeakenScore:    fmt.Sprintf("%d", int(weakenScore)),
		NonWeakenScore: fmt.Sprintf("%d", int(nonWeakenScore)),
	}
}
