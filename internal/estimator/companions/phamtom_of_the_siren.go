package companions

import "lysk-battle-record/internal/models"

type PhantomOfTheSiren struct{}

func (p PhantomOfTheSiren) GetName() string {
	return "海妖魅影"
}

func (p PhantomOfTheSiren) GetCompanionFlow(stats models.Stats) models.CompanionFlow {
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
				Boost:      10 * 0.5,
			},
		},
	}
}

func (p PhantomOfTheSiren) GetActiveSkill(stats models.Stats) models.Skill {
	energy := stats.GetEnergy()
	return getActiveSkillForWeapon(stats.Weapon, energy)
}

func (p PhantomOfTheSiren) GetBasicAttack(stats models.Stats) models.Skill {
	return getBasicAttackForWeapon(stats.Weapon)
}

func (p PhantomOfTheSiren) GetOathSkill(stats models.Stats) models.Skill {
	skill := getDefaultOathSkill()
	skill.Base = 1200
	skill.AttackRate = 1600
	skill.OathBoost = stats.OathBoost
	skill.Count = getOathCount(stats)
	return skill
}

func (p PhantomOfTheSiren) GetResonanceSkill(stats models.Stats) models.Skill {
	skill := getDefaultResonanceSkill()
	skill.Base = 508 + 254
	skill.AttackRate = 678 + 339
	return skill
}

func (p PhantomOfTheSiren) GetSupportSkill(stats models.Stats) models.Skill {
	skill := getDefaultSupportSkill()
	skill.Base = 218
	skill.AttackRate = 291
	skill.Count = 6
	return skill
}

func (p PhantomOfTheSiren) GetPassiveSkill(stats models.Stats) models.Skill {
	return models.Skill{
		Name:       "回声",
		Base:       18,
		AttackRate: 24,
		Count:      5 * 6,
	}
}
