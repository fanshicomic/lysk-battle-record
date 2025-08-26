package companions

import (
	"lysk-battle-record/internal/models"
	"math"
)

type GodOfTheTides struct{}

func (p GodOfTheTides) GetName() string {
	return "潮汐之神"
}

func (p GodOfTheTides) GetCompanionFlow(stats models.Stats) models.CompanionFlow {
	activeSkill := p.GetActiveSkill(stats)
	heavyAttack := p.GetBasicAttack(stats)
	resonanceSkill := p.GetResonanceSkill(stats)
	oathSkill := p.GetOathSkill(stats)
	supportSkill := p.GetSupportSkill(stats)
	passiveSkill := p.GetPassiveSkill(stats)

	weakenRate := getWeakenRate(stats.Matching)

	flow := models.CompanionFlow{
		Periods: []models.CompanionPeriod{
			{
				SkillSet: models.CompanionSkillSet{
					Skills: []models.Skill{
						activeSkill,
						heavyAttack,
						resonanceSkill,
						oathSkill,
						supportSkill,
						passiveSkill,
					},
				},
				WeakenRate: weakenRate,
				Boost:      30 * (float64(p.GetRainCount(stats)) / 6), // 下雨30%增伤，持续10秒
			},
		},
	}

	return flow
}

func (p GodOfTheTides) GetActiveSkill(stats models.Stats) models.Skill {
	energy := stats.GetEnergy()
	if stats.SetCard == "神殿" && (stats.Stage == "III" || stats.Stage == "IV") {
		energy += 2 * p.GetRainCount(stats) // 神殿III/IV阶增加2点能量
	}

	count := int(math.Min(float64(energy-8), 6))

	if stats.Weapon == "专武" {
		skill := getDefaultActiveSkill()
		skill.Base = 73
		skill.AttackRate = 39
		skill.HpRate = 3.5
		skill.Count = count
		skill.CritRate = 30 * float64(count) * 6 / 60
		return skill
	}

	return getActiveSkillForWeapon(stats.Weapon, energy)
}

func (p GodOfTheTides) GetBasicAttack(stats models.Stats) models.Skill {
	if stats.Weapon == "专武" {
		skill := getDefaultBasicAttack()
		skill.Base = 182
		skill.AttackRate = 97
		skill.HpRate = 9
		skill.Count = 30
		skill.CritRate = p.getExtraCritRate(stats)
		return skill
	}

	return getBasicAttackForWeapon(stats.Weapon)
}

func (p GodOfTheTides) GetResonanceSkill(stats models.Stats) models.Skill {
	skill := getDefaultResonanceSkill()
	skill.Base = 995
	skill.AttackRate = 531
	skill.HpRate = 47.8
	skill.CritRate = p.getExtraCritRate(stats)
	return skill
}

func (p GodOfTheTides) GetOathSkill(stats models.Stats) models.Skill {
	skill := getDefaultOathSkill()
	skill.Base = 1440
	skill.AttackRate = 780
	skill.HpRate = 69.4
	skill.OathBoost = stats.OathBoost
	skill.Count = getOathCount(stats)
	return skill
}

func (p GodOfTheTides) GetSupportSkill(stats models.Stats) models.Skill {
	skill := getDefaultSupportSkill()
	skill.Count = 6
	return skill
}

func (p GodOfTheTides) GetPassiveSkill(stats models.Stats) models.Skill {
	singleTimeCount := 7
	activeSkillCount := p.GetActiveSkill(stats).Count
	supportSkillCount := p.GetSupportSkill(stats).Count

	if stats.Weapon != "专武" {
		activeSkillCount = 0
	}

	passiveSkill := models.Skill{
		Name:        "海灵",
		Base:        47,
		AttackRate:  25,
		HpRate:      2.2,
		Count:       (activeSkillCount + supportSkillCount) * singleTimeCount,
		DamageBoost: ((float64(p.GetRainCount(stats))/6.0)*1.25 + 5.0/6.0) * 100 / 6.0, // 下雨期间海灵升级增益
		CritRate:    p.getExtraCritRate(stats),
		CanBeCrit:   true,
	}

	return passiveSkill
}

func (p GodOfTheTides) getExtraCritRate(stats models.Stats) float64 {
	if stats.Weapon != "专武" {
		return 0
	}

	activeSkillCount := p.GetActiveSkill(stats).Count
	// 主动释放后增加30%暴击率，持续6秒
	critRate := 30 * activeSkillCount * 6 / 60
	return float64(critRate)
}

func (p GodOfTheTides) GetRainCount(stats models.Stats) int {
	if stats.SetCard == "神殿" {
		if stats.Stage == "IV" {
			return 4
		} else if stats.Stage == "III" {
			return 3
		} else if stats.Stage == "II" {
			return 2
		}
	}
	return 1
}
