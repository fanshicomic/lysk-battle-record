package partners

import "lysk-battle-record/internal/models"

type MasterOfFate struct{}

func (p MasterOfFate) GetName() string {
	return "九黎司命"
}

func (p MasterOfFate) GetPartnerFlow(stats models.Stats) models.PartnerFlow {
	activeSkill := p.GetActiveSkill(stats)
	heavyAttack := p.GetHeavyAttack(stats)
	resonanceSkill := p.GetResonanceSkill(stats)
	oathSkill := p.GetOathSkill(stats)
	supportSkill := p.GetSupportSkill()
	passiveSkill := p.GetPassiveSkill(stats)

	weakenRate := getWeakenRate(stats.Matching)
	if stats.SetCard == "拥雪" && stats.Stage != "I" {
		weakenRate *= 1.2
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
		Boost:      21, // 溯光力场内10%攻击增益+破盾后增伤20%
		WeakenRate: weakenRate,
	}

	return flow
}

func (p MasterOfFate) GetActiveSkill(stats models.Stats) models.Skill {
	energy := stats.GetEnergy()
	if stats.SetCard == "拥雪" && (stats.Stage == "III" || stats.Stage == "IV") {
		energy += 2
	}
	if stats.Weapon == "专武" {
		return models.Skill{
			Name:       "主动",
			Base:       404,
			AttackRate: 539,
			Count:      energy - 8,
			CanBeCrit:  true,
		}
	}

	return getActiveSkillForWeapon(stats.Weapon, energy)
}

func (p MasterOfFate) GetHeavyAttack(stats models.Stats) models.Skill {
	if stats.Weapon == "专武" {
		return models.Skill{
			Name:       "重击",
			Base:       141,
			AttackRate: 188,
			Count:      30,
			CanBeCrit:  true,
		}
	}

	return getHeavyAttackForWeapon(stats.Weapon)
}

func (p MasterOfFate) GetResonanceSkill(stats models.Stats) models.Skill {
	resonanceSkill := models.Skill{
		Name:       "共鸣",
		Base:       632,
		AttackRate: 842,
		Count:      4,
		CanBeCrit:  true,
	}

	return resonanceSkill
}

func (p MasterOfFate) GetOathSkill(stats models.Stats) models.Skill {
	oathSkill := models.Skill{
		Name:        "誓约",
		Base:        1440,
		AttackRate:  1920,
		DamageBoost: stats.OathBoost,
		Count:       getOathCount(stats),
	}

	return oathSkill
}

func (p MasterOfFate) GetSupportSkill() models.Skill {
	supportSkill := models.Skill{
		Name:       "协助",
		Base:       260,
		AttackRate: 348,
		Count:      6,
		CanBeCrit:  true,
	}

	return supportSkill
}

func (p MasterOfFate) GetPassiveSkill(stats models.Stats) models.Skill {
	activeSkillCount := p.GetActiveSkill(stats).Count
	supportSkillCount := p.GetSupportSkill().Count
	altResonanceSkillCount := p.GetAltResonanceSkill(stats).Count
	passiveSkill := models.Skill{
		Name:       "断玉诀",
		Base:       233,
		AttackRate: 310,
		Count:      (activeSkillCount*4+supportSkillCount+altResonanceSkillCount)/3 + 6, // 6 from normal attacks
		CanBeCrit:  true,
	}

	return passiveSkill
}

func (p MasterOfFate) GetAltResonanceSkill(stats models.Stats) models.Skill {
	skill := models.Skill{
		Name:       "穿雨",
		Base:       205,
		AttackRate: 273,
		Count:      3 * 4,
		CanBeCrit:  true,
	}

	if stats.Weapon != "专武" {
		skill.Count = 0
	}

	return skill
}
