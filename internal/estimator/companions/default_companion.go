package companions

import "lysk-battle-record/internal/models"

type DefaultCompanion struct{}

func (p DefaultCompanion) GetName() string {
	return "默认搭档"
}

func (p DefaultCompanion) GetCompanionFlow(stats models.Stats) models.CompanionFlow {
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

func (p DefaultCompanion) GetActiveSkill(stats models.Stats) models.Skill {
	energy := stats.GetEnergy()
	if stats.Weapon == "专武" {
		skill := getDefaultActiveSkill()
		skill.Count = energy - 8
		return skill
	}

	return getActiveSkillForWeapon(stats.Weapon, energy)
}

func (p DefaultCompanion) GetBasicAttack(stats models.Stats) models.Skill {
	if stats.Weapon == "专武" {
		return getDefaultBasicAttack()
	}

	return getBasicAttackForWeapon(stats.Weapon)
}

func (p DefaultCompanion) GetOathSkill(stats models.Stats) models.Skill {
	skill := getDefaultOathSkill()
	skill.OathBoost = stats.OathBoost
	skill.Count = getOathCount(stats)
	return skill
}

func (p DefaultCompanion) GetResonanceSkill(stats models.Stats) models.Skill {
	skill := getDefaultResonanceSkill()
	return skill
}

func (p DefaultCompanion) GetSupportSkill(stats models.Stats) models.Skill {
	skill := getDefaultSupportSkill()
	skill.Count = 6
	return skill
}

func (p DefaultCompanion) GetPassiveSkill(stats models.Stats) models.Skill {
	return models.Skill{
		Name: "被动",
	}
}
