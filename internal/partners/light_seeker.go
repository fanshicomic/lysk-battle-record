package partners

import "lysk-battle-record/internal/models"

type LightSeeker struct{}

func (p LightSeeker) GetName() string {
	return "逐光骑士"
}

func (p LightSeeker) GetPartnerFlow(stats models.Stats) models.PartnerFlow {
	activeSkill := p.GetActiveSkill(stats)
	heavyAttack := p.GetHeavyAttack(stats)
	resonanceSkill := p.GetResonanceSkill(stats)
	oathSkill := p.GetOathSkill(stats)
	supportSkill := p.GetSupportSkill()

	passiveSkill := models.Skill{
		Name:       "溯光共鸣",
		Base:       150,
		AttackRate: 200,
		Count:      int(float64(activeSkill.Count) * 0.8),
		CanBeCrit:  true,
	}

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
		Boost:      21, // 溯光力场内10%攻击增益+破盾后增伤20%
		WeakenRate: weakenRate,
	}

	return flow
}

func (p LightSeeker) GetActiveSkill(stats models.Stats) models.Skill {
	energy := stats.GetEnergy()

	if stats.Weapon == "专武" {
		return models.Skill{
			Name:       "主动",
			Base:       341,
			AttackRate: 455,
			Count:      (energy - 8) * 2,
			CanBeCrit:  true,
		}
	}

	return getActiveSkillForWeapon(stats.Weapon, energy)
}

func (p LightSeeker) GetHeavyAttack(stats models.Stats) models.Skill {
	if stats.Weapon == "专武" {
		return models.Skill{
			Name:       "重击",
			Base:       118,
			AttackRate: 157,
			Count:      30,
			CanBeCrit:  true,
		}
	}

	return getHeavyAttackForWeapon(stats.Weapon)
}

func (p LightSeeker) GetResonanceSkill(stats models.Stats) models.Skill {
	resonanceSkill := models.Skill{
		Name:       "共鸣",
		Base:       641,
		AttackRate: 854,
		Count:      4,
		CanBeCrit:  true,
	}

	return resonanceSkill
}

func (p LightSeeker) GetOathSkill(stats models.Stats) models.Skill {
	oathSkill := models.Skill{
		Name:        "誓约",
		Base:        1440,
		AttackRate:  1920,
		DamageBoost: stats.OathBoost * 100,
		Count:       getOathCount(stats),
	}

	return oathSkill
}

func (p LightSeeker) GetSupportSkill() models.Skill {
	supportSkill := models.Skill{
		Name:       "协助",
		Base:       400,
		AttackRate: 400,
		Count:      6,
		CanBeCrit:  true,
	}

	return supportSkill
}
