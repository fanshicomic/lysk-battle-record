package partners

import "lysk-battle-record/internal/models"

type RelentLessConqueror struct{}

func (p RelentLessConqueror) GetName() string {
	return "无尽掠夺者"
}

func (p RelentLessConqueror) GetPartnerFlow(stats models.Stats) models.PartnerFlow {
	activeSkill := p.GetActiveSkill(stats)
	heavyAttack := p.GetHeavyAttack(stats)
	resonanceSkill := p.GetResonanceSkill(stats)
	oathSkill := p.GetOathSkill(stats)
	supportSkill := p.GetSupportSkill()
	passiveSkill := p.GetPassiveSkill(stats)

	weakenRate := getWeakenRate(stats.Matching)
	boost := (4.0 * 8.0 / 60.0) * 80.0
	if stats.SetCard == "掠心" && stats.Stage != "I" {
		boost = 80
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
		WeakenRate: weakenRate,
		Boost:      boost,
	}

	return flow
}

func (p RelentLessConqueror) GetActiveSkill(stats models.Stats) models.Skill {
	energy := stats.GetEnergy()

	if stats.Weapon == "专武" {
		return models.Skill{
			Name:        "主动",
			Base:        342,
			AttackRate:  456,
			Count:       energy - 8,
			DamageBoost: 10.0 / (10.0 - 6.0*0.5),
			CanBeCrit:   true,
		}
	}

	return getActiveSkillForWeapon(stats.Weapon, energy)
}

func (p RelentLessConqueror) GetHeavyAttack(stats models.Stats) models.Skill {
	if stats.Weapon == "专武" {
		return models.Skill{
			Name:       "重击",
			Base:       160,
			AttackRate: 213,
			Count:      20,
			CanBeCrit:  true,
		}
	}

	return getHeavyAttackForWeapon(stats.Weapon)
}

func (p RelentLessConqueror) GetResonanceSkill(stats models.Stats) models.Skill {
	resonanceSkill := models.Skill{
		Name:       "共鸣",
		Base:       1094,
		AttackRate: 1458,
		Count:      4,
		CanBeCrit:  true,
	}

	return resonanceSkill
}

func (p RelentLessConqueror) GetOathSkill(stats models.Stats) models.Skill {
	oathSkill := models.Skill{
		Name:        "誓约",
		Base:        1440,
		AttackRate:  1920,
		DamageBoost: stats.OathBoost,
		Count:       getOathCount(stats),
	}

	return oathSkill
}

func (p RelentLessConqueror) GetSupportSkill() models.Skill {
	supportSkill := models.Skill{
		Name:       "协助",
		Base:       239,
		AttackRate: 318,
		Count:      6,
		CanBeCrit:  true,
	}

	return supportSkill
}

func (p RelentLessConqueror) GetPassiveSkill(stats models.Stats) models.Skill {
	skill := models.Skill{
		Name:       "掠噬标记",
		Base:       60,
		AttackRate: 80,
		Count:      60 * 2 / 8,
		CanBeCrit:  true,
	}

	if stats.SetCard != "掠心" || stats.Stage != "IV" {
		skill.Count = 0
	}

	return skill
}
