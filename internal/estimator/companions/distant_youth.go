package companions

import "lysk-battle-record/internal/models"

type DistantYouth struct{}

func (p DistantYouth) GetName() string {
	return "遥远少年"
}

func (p DistantYouth) GetCompanionFlow(stats models.Stats) models.CompanionFlow {
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

func (p DistantYouth) GetActiveSkill(stats models.Stats) models.Skill {
	energy := stats.GetEnergy()
	skill := getActiveSkillForWeapon(stats.Weapon, energy)
	skill.DamageBoost = 30
	return skill
}

func (p DistantYouth) GetBasicAttack(stats models.Stats) models.Skill {
	return getBasicAttackForWeapon(stats.Weapon)
}

func (p DistantYouth) GetOathSkill(stats models.Stats) models.Skill {
	skill := getDefaultOathSkill()
	skill.Base = 1200
	skill.AttackRate = 1600
	skill.OathBoost = stats.OathBoost
	skill.Count = getOathCount(stats)
	return skill
}

func (p DistantYouth) GetResonanceSkill(stats models.Stats) models.Skill {
	skill := getDefaultResonanceSkill()
	skill.Base = 739
	skill.AttackRate = 986
	return skill
}

func (p DistantYouth) GetSupportSkill(stats models.Stats) models.Skill {
	skill := getDefaultSupportSkill()
	skill.Base = 340
	skill.AttackRate = 453
	skill.Count = 6
	return skill
}

func (p DistantYouth) GetPassiveSkill(stats models.Stats) models.Skill {
	return models.Skill{
		Name:       "剑意",
		AttackRate: 100,
		Count:      6,
	}
}
