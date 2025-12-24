package companions

import "lysk-battle-record/internal/models"

type Artist struct{}

func (p Artist) GetName() string {
	return "艺术家"
}

func (p Artist) GetCompanionFlow(stats models.Stats) models.CompanionFlow {
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
				Boost:      5 * 5 * 0.5,
			},
		},
	}
}

func (p Artist) GetActiveSkill(stats models.Stats) models.Skill {
	energy := stats.GetEnergy()
	return getActiveSkillForWeapon(stats.Weapon, energy)
}

func (p Artist) GetBasicAttack(stats models.Stats) models.Skill {
	return getBasicAttackForWeapon(stats.Weapon)
}

func (p Artist) GetOathSkill(stats models.Stats) models.Skill {
	skill := getDefaultOathSkill()
	skill.Base = 1200
	skill.AttackRate = 1600
	skill.OathBoost = stats.OathBoost
	skill.Count = getOathCount(stats)
	return skill
}

func (p Artist) GetResonanceSkill(stats models.Stats) models.Skill {
	skill := getDefaultResonanceSkill()
	skill.Base = 118 * 5
	skill.AttackRate = 250 * 5
	return skill
}

func (p Artist) GetSupportSkill(stats models.Stats) models.Skill {
	skill := getDefaultSupportSkill()
	skill.Base = 84
	skill.AttackRate = 112
	skill.Count = 6
	return skill
}

func (p Artist) GetPassiveSkill(stats models.Stats) models.Skill {
	return models.Skill{
		Name:       "火焰陷阱",
		Base:       15,
		AttackRate: 20,
		Count:      30,
	}
}

func (p Artist) GetBurnSkill() models.Skill {
	return models.Skill{
		Name:       "灼烧",
		Base:       20,
		AttackRate: 26,
		Count:      10 * 6,
	}
}
