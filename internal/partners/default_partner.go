package partners

import "lysk-battle-record/internal/models"

type DefaultPartner struct{}

func (p DefaultPartner) GetName() string {
	return "默认搭档"
}

func (p DefaultPartner) GetPartnerFlow(stats models.Stats) models.PartnerFlow {
	activeSkill := p.GetActiveSkill(stats)
	heavyAttack := p.GetHeavyAttack(stats)
	resonanceSkill := p.GetResonanceSkill(stats)
	oathSkill := p.GetOathSkill(stats)
	supportSkill := p.GetSupportSkill()
	passiveSkill := p.GetPassiveSkill(stats)

	weakenRate := getWeakenRate(stats.Matching)
	return models.PartnerFlow{
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
	}
}

func (p DefaultPartner) GetActiveSkill(stats models.Stats) models.Skill {
	energy := stats.GetEnergy()
	if stats.Weapon == "专武" {
		return models.Skill{
			Name:       "主动",
			Base:       309,
			AttackRate: 412,
			Count:      energy - 8,
			CanBeCrit:  true,
		}
	}

	return getActiveSkillForWeapon(stats.Weapon, energy)
}

func (p DefaultPartner) GetHeavyAttack(stats models.Stats) models.Skill {
	if stats.Weapon == "专武" {
		return models.Skill{
			Name:      "重击",
			CanBeCrit: true,
		}
	}

	return getHeavyAttackForWeapon(stats.Weapon)
}

func (p DefaultPartner) GetOathSkill(stats models.Stats) models.Skill {
	oathSkill := models.Skill{
		Name:        "誓约",
		DamageBoost: stats.OathBoost,
		Count:       getOathCount(stats),
	}

	return oathSkill
}

func (p DefaultPartner) GetResonanceSkill(stats models.Stats) models.Skill {
	resonanceSkill := models.Skill{
		Name:      "共鸣",
		Count:     4,
		CanBeCrit: true,
	}

	return resonanceSkill
}

func (p DefaultPartner) GetSupportSkill() models.Skill {
	supportSkill := models.Skill{
		Name:      "协助",
		Count:     6,
		CanBeCrit: true,
	}

	return supportSkill
}

func (p DefaultPartner) GetPassiveSkill(stats models.Stats) models.Skill {
	return models.Skill{}
}
