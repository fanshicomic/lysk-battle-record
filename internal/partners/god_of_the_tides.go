package partners

import (
	"lysk-battle-record/internal/models"
	"math"
)

type GodOfTheTides struct{}

func (p GodOfTheTides) GetName() string {
	return "潮汐之神"
}

func (p GodOfTheTides) GetPartnerFlow(stats models.Stats) models.PartnerFlow {
	activeSkill := p.GetActiveSkill(stats)
	heavyAttack := p.GetHeavyAttack(stats)
	resonanceSkill := p.GetResonanceSkill(stats)
	oathSkill := p.GetOathSkill(stats)
	supportSkill := p.GetSupportSkill()
	passiveSkill := p.GetPassiveSkill(stats)

	weakenRate := getWeakenRate(stats.Matching)

	flow := models.PartnerFlow{
		Periods: []models.PartnerPeriod{
			{
				SkillSet: models.PartnerSkillSet{
					Skills: []models.Skill{
						activeSkill,
						heavyAttack,
						resonanceSkill,
						oathSkill,
						supportSkill,
						passiveSkill,
					},
				},
			},
		},
		Boost:      30 * (float64(p.GetRainCount(stats)) / 6), // 下雨30%增伤，持续10秒
		WeakenRate: weakenRate,
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
		return models.Skill{
			Name:       "主动",
			Base:       73,
			AttackRate: 39,
			HpRate:     3.5,
			Count:      count,
			CritRate:   30 * float64(count) * 6 / 60,
			CanBeCrit:  true,
		}
	}

	return getActiveSkillForWeapon(stats.Weapon, energy)
}

func (p GodOfTheTides) GetHeavyAttack(stats models.Stats) models.Skill {
	if stats.Weapon == "专武" {
		return models.Skill{
			Name:       "重击",
			Base:       182,
			AttackRate: 97,
			HpRate:     9,
			Count:      30,
			CritRate:   p.getExtraCritRate(stats),
			CanBeCrit:  true,
		}
	}

	return getHeavyAttackForWeapon(stats.Weapon)
}

func (p GodOfTheTides) GetResonanceSkill(stats models.Stats) models.Skill {
	resonanceSkill := models.Skill{
		Name:       "共鸣",
		Base:       995,
		AttackRate: 531,
		HpRate:     47.8,
		Count:      4,
		CritRate:   p.getExtraCritRate(stats),
		CanBeCrit:  true,
	}

	return resonanceSkill
}

func (p GodOfTheTides) GetOathSkill(stats models.Stats) models.Skill {
	oathSkill := models.Skill{
		Name:        "誓约",
		Base:        1440,
		AttackRate:  780,
		HpRate:      69.4,
		DamageBoost: stats.OathBoost * 100,
		Count:       getOathCount(stats),
	}

	return oathSkill
}

func (p GodOfTheTides) GetSupportSkill() models.Skill {
	supportSkill := models.Skill{
		Name:  "协助",
		Count: 6,
	}

	return supportSkill
}

func (p GodOfTheTides) GetPassiveSkill(stats models.Stats) models.Skill {
	singleTimeCount := 7
	activeSkillCount := p.GetActiveSkill(stats).Count
	supportSkillCount := p.GetSupportSkill().Count

	if stats.Weapon != "专武" {
		activeSkillCount = 0
	}

	passiveSkill := models.Skill{
		Name:        "海灵",
		Base:        47,
		AttackRate:  25,
		HpRate:      2.2,
		Count:       (activeSkillCount + supportSkillCount) * singleTimeCount,
		DamageBoost: ((float64(p.GetRainCount(stats))/6)*1.25 + 5/6) / 6, // 下雨期间海灵升级增益
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
