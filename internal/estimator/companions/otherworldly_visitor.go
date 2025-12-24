package companions

import "lysk-battle-record/internal/models"

type OtherworldlyVisitor struct{}

func (p OtherworldlyVisitor) GetName() string {
	return "异界来客"
}

func (p OtherworldlyVisitor) GetCompanionFlow(stats models.Stats) models.CompanionFlow {
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
				Boost:      8 * 0.6,
			},
		},
	}
}

func (p OtherworldlyVisitor) GetActiveSkill(stats models.Stats) models.Skill {
	energy := stats.GetEnergy()
	return getActiveSkillForWeapon(stats.Weapon, energy)
}

func (p OtherworldlyVisitor) GetBasicAttack(stats models.Stats) models.Skill {
	return getBasicAttackForWeapon(stats.Weapon)
}

func (p OtherworldlyVisitor) GetOathSkill(stats models.Stats) models.Skill {
	skill := getDefaultOathSkill()
	skill.Base = 1200
	skill.AttackRate = 1600
	skill.OathBoost = stats.OathBoost
	skill.Count = getOathCount(stats)
	return skill
}

func (p OtherworldlyVisitor) GetResonanceSkill(stats models.Stats) models.Skill {
	skill := getDefaultResonanceSkill()
	skill.Base = 950
	skill.AttackRate = 1266
	return skill
}

func (p OtherworldlyVisitor) GetSupportSkill(stats models.Stats) models.Skill {
	skill := getDefaultSupportSkill()
	skill.Base = 321
	skill.AttackRate = 429
	skill.Count = 4
	return skill
}

func (p OtherworldlyVisitor) GetPassiveSkill(stats models.Stats) models.Skill {
	return models.Skill{}
}
