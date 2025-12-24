package companions

import "lysk-battle-record/internal/models"

type EvolPolice struct{}

func (p EvolPolice) GetName() string {
	return "Evol特警"
}

func (p EvolPolice) GetCompanionFlow(stats models.Stats) models.CompanionFlow {
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
				Boost:      20 * 0.25,
			},
		},
	}
}

func (p EvolPolice) GetActiveSkill(stats models.Stats) models.Skill {
	energy := stats.GetEnergy()
	return getActiveSkillForWeapon(stats.Weapon, energy)
}

func (p EvolPolice) GetBasicAttack(stats models.Stats) models.Skill {
	return getBasicAttackForWeapon(stats.Weapon)
}

func (p EvolPolice) GetOathSkill(stats models.Stats) models.Skill {
	skill := getDefaultOathSkill()
	skill.Base = 1200
	skill.AttackRate = 1600
	skill.OathBoost = stats.OathBoost
	skill.Count = getOathCount(stats)
	return skill
}

func (p EvolPolice) GetResonanceSkill(stats models.Stats) models.Skill {
	skill := getDefaultResonanceSkill()
	skill.Base = 729
	skill.AttackRate = 968
	return skill
}

func (p EvolPolice) GetSupportSkill(stats models.Stats) models.Skill {
	skill := getDefaultSupportSkill()
	skill.Base = 306
	skill.AttackRate = 408
	skill.Count = 6
	return skill
}

func (p EvolPolice) GetPassiveSkill(stats models.Stats) models.Skill {
	return models.Skill{}
}
