package partners

import "lysk-battle-record/internal/models"

type PartnerLightSeeker struct{}

func (p PartnerLightSeeker) GetName() string {
	return "逐光骑士"
}

func (p PartnerLightSeeker) GetPartnerFlow(stats models.Stats) models.PartnerFlow {
	activeSkill := p.GetActiveSkill(stats)
	heavyAttack := p.GetHeavyAttack(stats)
	resonanceSkill := p.GetResonanceSkill()
	oathSkill := p.GetOathSkill(stats)
	supportSkill := p.GetSupportSkill()

	passiveSkill := models.Skill{
		Name:       "溯光共鸣",
		Base:       150,
		AttackRate: 200,
	}

	weakenRate := 0.17
	if stats.Matching == "顺" {
		weakenRate = 0.34
	}

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
		Boost:      21,
		WeakenRate: weakenRate,
	}

	return flow
}

func (p PartnerLightSeeker) GetActiveSkill(stats models.Stats) models.Skill {
	energy := stats.GetEnergy()
	activeSkill := models.Skill{
		Name:  "主动",
		Count: energy - 8,
	}

	switch stats.Weapon {
	case "专武":
		activeSkill.Base = 341
		activeSkill.AttackRate = 455
		activeSkill.Count = energy - 8
		activeSkill.Count *= 2 // 主动命中减冷却加能量
		break
	case "重剑":
		activeSkill.Name = "重剑主动"
		activeSkill.Base = 621
		activeSkill.AttackRate = 829
		activeSkill.DamageBoost = 150
		break
	case "单手剑":
		activeSkill.Name = "单手剑主动"
		activeSkill.Base = 341
		activeSkill.AttackRate = 455
		break
	case "法杖":
		activeSkill.Name = "法杖主动"
		activeSkill.Base = 204
		activeSkill.AttackRate = 270
		break
	case "手枪":
		activeSkill.Name = "手枪主动"
		activeSkill.Base = 160
		activeSkill.AttackRate = 213
		break
	}

	return activeSkill
}

func (p PartnerLightSeeker) GetHeavyAttack(stats models.Stats) models.Skill {
	heavyAttack := models.Skill{
		Name:  "重击",
		Count: 30,
	}

	switch stats.Weapon {
	case "专武":
		heavyAttack.Base = 118
		heavyAttack.AttackRate = 157
	case "重剑":
		heavyAttack.Base = 337
		heavyAttack.AttackRate = 449
		heavyAttack.DamageBoost = 112.9
		heavyAttack.Count = 10
	case "单手剑":
		heavyAttack.Base = 250
		heavyAttack.AttackRate = 333
		heavyAttack.DamageBoost = 14
	case "法杖":
		heavyAttack.Base = 122
		heavyAttack.AttackRate = 162
		heavyAttack.DamageBoost = 28
		heavyAttack.Count = 15
	case "手枪":
		heavyAttack.Base = 120
		heavyAttack.AttackRate = 160
		heavyAttack.DamageBoost = 25
		heavyAttack.Count = 35
	}

	return heavyAttack
}

func (p PartnerLightSeeker) GetResonanceSkill() models.Skill {
	resonanceSkill := models.Skill{
		Name:       "共鸣",
		Count:      4,
		Base:       641,
		AttackRate: 854,
	}

	return resonanceSkill
}

func (p PartnerLightSeeker) GetOathSkill(stats models.Stats) models.Skill {
	oathSkill := models.Skill{
		Name:        "誓约",
		Base:        1440,
		AttackRate:  1920,
		DamageBoost: stats.OathBoost * 100,
	}

	if stats.Stage != "无套装" && stats.Stage != "I" {
		oathSkill.Count = 1
	}
	if stats.OathRegen >= 17 {
		oathSkill.Count = 1
	}

	return oathSkill
}

func (p PartnerLightSeeker) GetSupportSkill() models.Skill {
	supportSkill := models.Skill{
		Name:       "协助",
		Base:       400,
		AttackRate: 400,
		Count:      6,
	}

	return supportSkill
}
