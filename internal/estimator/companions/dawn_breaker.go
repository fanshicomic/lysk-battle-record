package companions

import "lysk-battle-record/internal/models"

type DawnBreaker struct{}

func (p DawnBreaker) GetName() string {
	return "黎明抹杀者"
}

func (p DawnBreaker) GetCompanionFlow(stats models.Stats) models.CompanionFlow {
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

func (p DawnBreaker) GetActiveSkill(stats models.Stats) models.Skill {
	energy := stats.GetEnergy()
	return getActiveSkillForWeapon(stats.Weapon, energy)
}

func (p DawnBreaker) GetBasicAttack(stats models.Stats) models.Skill {
	return getBasicAttackForWeapon(stats.Weapon)
}

func (p DawnBreaker) GetOathSkill(stats models.Stats) models.Skill {
	skill := getDefaultOathSkill()
	skill.Base = 1200
	skill.AttackRate = 1600
	skill.OathBoost = stats.OathBoost
	skill.Count = getOathCount(stats)
	return skill
}

func (p DawnBreaker) GetResonanceSkill(stats models.Stats) models.Skill {
	skill := getDefaultResonanceSkill()
	skill.Base = 252
	skill.AttackRate = 336
	return skill
}

func (p DawnBreaker) GetSupportSkill(stats models.Stats) models.Skill {
	skill := getDefaultSupportSkill()
	skill.Base = 626
	skill.AttackRate = 835
	skill.Count = 6
	return skill
}

func (p DawnBreaker) GetPassiveSkill(stats models.Stats) models.Skill {
	return models.Skill{}
}
