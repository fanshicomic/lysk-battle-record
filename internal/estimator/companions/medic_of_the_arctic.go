package companions

import "lysk-battle-record/internal/models"

type MedicOfTheArctic struct{}

func (p MedicOfTheArctic) GetName() string {
	return "极地军医"
}

func (p MedicOfTheArctic) GetCompanionFlow(stats models.Stats) models.CompanionFlow {
	activeSkill := p.GetActiveSkill(stats)
	basicAttack := p.GetBasicAttack(stats)
	resonanceSkill := p.GetResonanceSkill(stats)
	oathSkill := p.GetOathSkill(stats)
	supportSkill := p.GetSupportSkill(stats)
	passiveSkill := p.GetPassiveSkill(stats)

	weakenRate := getWeakenRate(stats.Matching)
	return models.CompanionFlow{
		Periods: []models.CompanionPeriod{
			{
				SkillSet: models.CompanionSkillSet{
					Skills: []models.Skill{
						activeSkill,
						basicAttack,
						resonanceSkill,
						oathSkill,
						supportSkill,
						passiveSkill,
					},
				},
				WeakenRate: weakenRate,
			},
		},
	}
}

func (p MedicOfTheArctic) GetActiveSkill(stats models.Stats) models.Skill {
	stats.EnergyRegen += 24
	energy := stats.GetEnergy()
	energy += 1
	skill := getActiveSkillForWeapon(stats.Weapon, energy)
	skill.DamageBoost = 40 * 8.0 / 15.0
	skill.WeakenBoost = 34 // 军医破斩调参，应该有多一个破斩在虚弱期
	return skill
}

func (p MedicOfTheArctic) GetBasicAttack(stats models.Stats) models.Skill {
	skill := getBasicAttackForWeapon(stats.Weapon)
	skill.DamageBoost = 40 * 8.0 / 15.0
	return skill
}

func (p MedicOfTheArctic) GetOathSkill(stats models.Stats) models.Skill {
	skill := getDefaultOathSkill()
	skill.Base = 1200
	skill.AttackRate = 1600
	skill.OathBoost = stats.OathBoost
	skill.Count = getOathCount(stats)
	return skill
}

func (p MedicOfTheArctic) GetResonanceSkill(stats models.Stats) models.Skill {
	skill := getDefaultResonanceSkill()
	skill.Base = 362
	skill.AttackRate = 482
	return skill
}

func (p MedicOfTheArctic) GetSupportSkill(stats models.Stats) models.Skill {
	skill := getDefaultSupportSkill()
	skill.Base = 208
	skill.AttackRate = 275
	skill.Count = 6
	return skill
}

func (p MedicOfTheArctic) GetPassiveSkill(stats models.Stats) models.Skill {
	return models.Skill{
		Name: "被动",
	}
}
