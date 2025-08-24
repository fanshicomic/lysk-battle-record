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
	printPartnerFlow(partnerFlow)
	fmt.Printf("Final Score: %+v\n", score)
	return score
}

func getPartnerFlow(stats models.Stats) models.PartnerFlow {
	var partner partners.Partner
	switch stats.Partner {
	case "逐光骑士":
		partner = partners.LightSeeker{}
	case "永恒先知":
		partner = partners.Foreseer{}
	case "深海潜行者":
		partner = partners.AbyssWalker{}
	case "潮汐之神":
		partner = partners.GodOfTheTides{}
	case "光猎":
		partner = partners.Lumiere{}
	case "九黎司命":
		partner = partners.MasterOfFate{}
	case "无尽掠夺者":
		partner = partners.RelentLessConqueror{}
	default:
		partner = partners.DefaultPartner{}
	}
	return partner.GetPartnerFlow(stats)
}

func getSetCard(stats models.Stats) set_cards.SetCard {
	var setCard set_cards.SetCard
	switch stats.SetCard {
	case "逐光":
		setCard = set_cards.LightSeeking{}
	case "永恒":
		setCard = set_cards.Forever{}
	case "深海":
		setCard = set_cards.DeepSea{}
	case "神殿":
		setCard = set_cards.Temple{}
	case "末夜":
		setCard = set_cards.Midnight{}
	case "拥雪":
		setCard = set_cards.SnowFall{}
	case "掠心":
		setCard = set_cards.Captivating{}
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
			levelDefenseRatio := 1 + float64(stats.TotalLevel)/(float64(stats.TotalLevel)+300+(80*3+100)*(1-skill.EnemyDefenceReduction/100))
			rawSkillScore *= levelDefenseRatio

			// consider non-weaken period
			nonWeakenSkillCount := (1 - partnerFlow.WeakenRate) * float64(skill.Count)
			if skill.Name == "誓约" { // 誓约只在虚弱期放，X-02誓约技会算作特殊技能
				nonWeakenSkillCount = 0
			}
			nonWeakenPeriodScore := rawSkillScore * nonWeakenSkillCount
			critRate := (stats.CritRate + skill.CritRate) / 100
			if !skill.CanBeCrit {
				critRate = 0
			}
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

	total *= 1 + partnerFlow.Boost/100.0
	weakenScore *= 1 + partnerFlow.Boost/100.0
	nonWeakenScore *= 1 + partnerFlow.Boost/100.0
	buffedTotal := int(total * matchingBuff * championshipsBuff)

	return models.CombatPower{
		Score:          fmt.Sprintf("%d", int(total)),
		BuffedScore:    fmt.Sprintf("%d", buffedTotal),
		WeakenScore:    fmt.Sprintf("%d", int(weakenScore)),
		NonWeakenScore: fmt.Sprintf("%d", int(nonWeakenScore)),
	}
}

func printPartnerFlow(flow models.PartnerFlow) {
	fmt.Println("=== Partner Flow Debug ===")
	fmt.Printf("Boost: %.1f%% | WeakenRate: %.1f%%\n", flow.Boost, flow.WeakenRate*100)
	fmt.Println()

	for periodIdx, period := range flow.Periods {
		fmt.Printf("Period %d:\n", periodIdx+1)
		fmt.Println("Skills:")

		for _, skill := range period.SkillSet.Skills {
			fmt.Printf("  [%s]\n", skill.Name)
			fmt.Printf("    Base: %.0f | AttackRate: %.1f%% | Count: %d\n",
				skill.Base, skill.AttackRate, skill.Count)

			if skill.HpRate > 0 {
				fmt.Printf("    HpRate: %.1f%%", skill.HpRate)
			}
			if skill.DefenseRate > 0 {
				fmt.Printf("    DefenseRate: %.1f%%", skill.DefenseRate)
			}
			if skill.DamageBoost > 0 {
				fmt.Printf("    DamageBoost: %.1f%%", skill.DamageBoost)
			}
			if skill.CritRate > 0 {
				fmt.Printf("    CritRate: %.1f%%", skill.CritRate)
			}
			if skill.CritDmg > 0 {
				fmt.Printf("    CritDmg: %.1f%%", skill.CritDmg)
			}
			if skill.WeakenBoost > 0 {
				fmt.Printf("    WeakenBoost: %.1f%%", skill.WeakenBoost)
			}
			if skill.EnemyDefenceReduction > 0 {
				fmt.Printf("    DefenceReduction: %.1f%%", skill.EnemyDefenceReduction)
			}

			fmt.Printf("    CanBeCrit: %t\n", skill.CanBeCrit)
			fmt.Println()
		}
	}
	fmt.Println("========================")
}
