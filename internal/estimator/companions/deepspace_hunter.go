package companions

import "lysk-battle-record/internal/models"

type DeepspaceHunter struct{}

func (p DeepspaceHunter) GetName() string {
	return "深空猎人"
}

func (p DeepspaceHunter) GetCompanionFlow(stats models.Stats) models.CompanionFlow {
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

func (p DeepspaceHunter) GetActiveSkill(stats models.Stats) models.Skill {
	energy := stats.GetEnergy()
	skill := getActiveSkillForWeapon(stats.Weapon, energy)
	skill.CritRate = 7
	return skill
}

func (p DeepspaceHunter) GetBasicAttack(stats models.Stats) models.Skill {
	skill := getBasicAttackForWeapon(stats.Weapon)
	skill.CritRate = 7
	return skill
}

func (p DeepspaceHunter) GetOathSkill(stats models.Stats) models.Skill {
	skill := getDefaultOathSkill()
	skill.Base = 1200
	skill.AttackRate = 1600
	skill.OathBoost = stats.OathBoost
	skill.Count = getOathCount(stats)
	skill.CritRate = 7
	return skill
}

func (p DeepspaceHunter) GetResonanceSkill(stats models.Stats) models.Skill {
	skill := getDefaultResonanceSkill()
	skill.Base = 318
	skill.AttackRate = 424
	skill.CritRate = 7
	return skill
}

func (p DeepspaceHunter) GetSupportSkill(stats models.Stats) models.Skill {
	skill := getDefaultSupportSkill()
	skill.Base = 330
	skill.AttackRate = 440
	skill.Count = 6
	skill.CritRate = 7
	return skill
}

func (p DeepspaceHunter) GetPassiveSkill(stats models.Stats) models.Skill {
	return models.Skill{}
}
