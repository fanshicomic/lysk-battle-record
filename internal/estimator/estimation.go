package estimator

import (
	"fmt"

	"lysk-battle-record/internal/estimator/companions"
	"lysk-battle-record/internal/estimator/set_cards"
	"lysk-battle-record/internal/models"
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
	myCompanion := getCompanion(stats)
	flow := myCompanion.GetCompanionFlow(stats)
	setCard := getSetCard(stats)
	setCardBuff := getSetCardBuff(stats, setCard)
	applySetCardBuff(&flow, setCardBuff)

	score := estimate(stats, flow)
	//printCompanionFlow(flow)
	//fmt.Printf("Final Score: %+v\n", score)
	return score
}

func getCompanion(stats models.Stats) companions.Companion {
	switch stats.Companion {
	case "逐光骑士":
		return companions.LightSeeker{}
	case "永恒先知":
		return companions.Foreseer{}
	case "深海潜行者":
		return companions.AbyssWalker{}
	case "潮汐之神":
		return companions.GodOfTheTides{}
	case "光猎":
		return companions.Lumiere{}
	case "九黎司命":
		return companions.MasterOfFate{}
	case "无尽掠夺者":
		return companions.RelentLessConqueror{}
	case "深渊主宰":
		return companions.AbysmSovereign{}
	case "远空执舰官":
		return companions.FarspaceColonel{}
	case "终极兵器X-02":
		return companions.UltimateWeaponX02{}
	case "利莫里亚海神":
		return companions.LemurianSeaGod{}
	case "暗蚀国王":
		return companions.KingOfDarknight{}
	case "极地军医":
		return companions.MedicOfTheArctic{}
	default:
		return companions.AbysmSovereign{}
	}
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
	case "雾海":
		setCard = set_cards.Mistsea{}
	case "夜誓":
		setCard = set_cards.Nightvow{}
	default:
		setCard = set_cards.NoSet{}
	}
	return setCard
}

func getSetCardBuff(stats models.Stats, setCard set_cards.SetCard) models.StageBuff {
	setMap := map[string]string{
		"逐光骑士":     "逐光",
		"永恒先知":     "永恒",
		"深海潜行者":   "深海",
		"潮汐之神":     "神殿",
		"光猎":         "末夜",
		"九黎司命":     "拥雪",
		"无尽掠夺者":   "掠心",
		"深渊主宰":     "深渊",
		"远空执舰官":   "远空",
		"终极兵器X-02": "寂路",
		"利莫里亚海神": "雾海",
		"暗蚀国王":     "夜誓",
	}
	var setCardBuff models.StageBuff
	if setCard.GetName() == setMap[stats.Companion] {
		setCardBuff = setCard.GetSetCardBuff()[stats.Stage]
	} else if setCard.GetName() != "无套装" {
		setCardBuff = set_cards.GetDefaultCardBuff()[stats.Stage]
	} else {
		setCardBuff = set_cards.NoSet{}.GetSetCardBuff()[stats.Stage]
	}

	return setCardBuff
}

func applySetCardBuff(flow *models.CompanionFlow, buff models.StageBuff) {
	for periodIdx := range flow.Periods {
		period := &flow.Periods[periodIdx]
		for i, skill := range period.SkillSet.Skills {
			if allBuff, exists := buff.Buffs["所有"]; exists {
				skill.CritRate += allBuff.CritRate
				skill.CritDmg += allBuff.CritDmg
				skill.WeakenBoost += allBuff.WeakenBoost
				skill.DamageBoost += allBuff.DamageBoost
				skill.OathBoost += allBuff.OathBoost
				if allBuff.CountBonus > 1 {
					skill.Count = int(float64(skill.Count) * allBuff.CountBonus)
				}
			}

			if skillBuff, exists := buff.Buffs[skill.Name]; exists {
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
