package estimator

import (
	"fmt"

	"lysk-battle-record/internal/companions"
	"lysk-battle-record/internal/models"
	"lysk-battle-record/internal/set_cards"
)

type CombatPowerEstimator interface {
	EstimateCombatPower(record models.Record) models.CombatPower
}

type LyskCPEstimator struct{}

func NewCombatPowerEstimator() CombatPowerEstimator {
	return &LyskCPEstimator{}
}

func (e *LyskCPEstimator) EstimateCombatPower(record models.Record) models.CombatPower {
	stats := record.ToStats()
	companionFlow := getCompanionFlow(stats)
	setCard := getSetCard(stats)
	setCardBuff := setCard.GetSetCardBuff()[stats.Stage]
	applySetCardBuff(&companionFlow, setCardBuff)

	score := estimate(stats, companionFlow)
	printCompanionFlow(companionFlow)
	fmt.Printf("Final Score: %+v\n", score)
	return score
}

func getCompanionFlow(stats models.Stats) models.CompanionFlow {
	var companion companions.Companion
	switch stats.Companion {
	case "逐光骑士":
		companion = companions.LightSeeker{}
	case "永恒先知":
		companion = companions.Foreseer{}
	case "深海潜行者":
		companion = companions.AbyssWalker{}
	case "潮汐之神":
		companion = companions.GodOfTheTides{}
	case "光猎":
		companion = companions.Lumiere{}
	case "九黎司命":
		companion = companions.MasterOfFate{}
	case "无尽掠夺者":
		companion = companions.RelentLessConqueror{}
	case "深渊主宰":
		companion = companions.AbysmSovereign{}
	case "远空执舰官":
		companion = companions.FarspaceColonel{}
	case "终极兵器X-02":
		companion = companions.UltimateWeaponX02{}
	default:
		companion = companions.AbysmSovereign{}
	}
	return companion.GetCompanionFlow(stats)
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
	case "深渊":
		setCard = set_cards.Abyssal{}
	case "远空":
		setCard = set_cards.Farspace{}
	case "寂路":
		setCard = set_cards.LoneRoad{}
	default:
		setCard = set_cards.NoSet{}
	}
	return setCard
}

func applySetCardBuff(flow *models.CompanionFlow, setCardBuff models.StageBuff) {
	for periodIdx := range flow.Periods {
		period := &flow.Periods[periodIdx]
		for i, skill := range period.SkillSet.Skills {
			if allBuff, exists := setCardBuff.Buffs["所有"]; exists {
				skill.CritRate += allBuff.CritRate
				skill.CritDmg += allBuff.CritDmg
				skill.WeakenBoost += allBuff.WeakenBoost
				skill.DamageBoost += allBuff.DamageBoost
				skill.OathBoost += allBuff.OathBoost
				if allBuff.CountBonus > 1 {
					skill.Count = int(float64(skill.Count) * allBuff.CountBonus)
				}
			}

			if skillBuff, exists := setCardBuff.Buffs[skill.Name]; exists {
				skill.CritRate += skillBuff.CritRate
				skill.CritDmg += skillBuff.CritDmg
				skill.WeakenBoost += skillBuff.WeakenBoost
				skill.DamageBoost += skillBuff.DamageBoost
				skill.OathBoost += skillBuff.OathBoost
				if skillBuff.CountBonus > 1 {
					skill.Count = int(float64(skill.Count) * skillBuff.CountBonus)
				}
			}

			period.SkillSet.Skills[i] = skill
		}
	}
}

func estimate(stats models.Stats, companionFlow models.CompanionFlow) models.CombatPower {
	var total, weakenScore, nonWeakenScore float64 = 0, 0, 0
	for _, period := range companionFlow.Periods {
		var score float64 = 0

		for _, skill := range period.SkillSet.Skills {
			rawSkillScore := skill.Base +
				(skill.HpRate/100)*float64(stats.HP) +
				(skill.AttackRate/100)*float64(stats.Attack) +
				(skill.DefenseRate/100)*float64(stats.Defense)

			// apply damage boost and period boost
			rawSkillScore *= 1 + (skill.DamageBoost+period.Boost)/100

			// apply oath boost
			if skill.Name == "誓约" || skill.Name == "誓约-同频觉醒" || skill.Name == "誓约-同频攻击" {
				rawSkillScore *= 1 + skill.OathBoost/100
			}

			// consider level - defence relationship
			levelDefenseRatio := 1 + float64(stats.TotalLevel)/(float64(stats.TotalLevel)+300+(80*3+100)*(1-skill.EnemyDefenceReduction/100))
			rawSkillScore *= levelDefenseRatio

			// consider non-weaken period
			weakenRate := period.WeakenRate
			if skill.NoWeakenPeriod {
				weakenRate = 0
			}

			if skill.Name == "誓约" {
				weakenRate = 1
			}
			nonWeakenSkillCount := (1 - weakenRate) * float64(skill.Count)
			nonWeakenPeriodScore := rawSkillScore * nonWeakenSkillCount
			critRate := (stats.CritRate + skill.CritRate) / 100
			if !skill.CanBeCrit {
				critRate = 0
			}
			critDmg := (stats.CritDmg + skill.CritDmg) / 100
			nonWeakenPeriodScore = nonWeakenPeriodScore*(1-critRate) +
				nonWeakenPeriodScore*critRate*critDmg

			// consider weaken period
			weakenSkillCount := weakenRate * float64(skill.Count)
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

	matchingBuff := 1 + stats.MatchingBuff/100.0
	championshipsBuff := 1 + stats.Buff/100.0

	buffedTotal := int(total * matchingBuff * championshipsBuff)

	return models.CombatPower{
		Score:          fmt.Sprintf("%d", int(total)),
		BuffedScore:    fmt.Sprintf("%d", buffedTotal),
		WeakenScore:    fmt.Sprintf("%d", int(weakenScore)),
		NonWeakenScore: fmt.Sprintf("%d", int(nonWeakenScore)),
	}
}

func printCompanionFlow(flow models.CompanionFlow) {
	fmt.Println("=== Companion Flow Debug ===")
	fmt.Println()

	for periodIdx, period := range flow.Periods {
		fmt.Printf("Period %d: | WeakenRate: %.1f%% | Period Boost: %.1f%%\n", periodIdx+1, period.WeakenRate*100, period.Boost)
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
			if skill.OathBoost > 0 {
				fmt.Printf("    OathBoost: %.1f%%", skill.OathBoost)
			}
			if skill.EnemyDefenceReduction > 0 {
				fmt.Printf("    DefenceReduction: %.1f%%", skill.EnemyDefenceReduction)
			}

			fmt.Printf("    CanBeCrit: %t\n", skill.CanBeCrit)
			fmt.Printf("    NoWeakenPeriod: %t\n", skill.NoWeakenPeriod)
			fmt.Println()
		}
	}
	fmt.Println("========================")
}
