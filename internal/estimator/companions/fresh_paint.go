package companions

import "lysk-battle-record/internal/models"

type FreshPaint struct{}

func (p FreshPaint) GetName() string {
	return "花坛新锐"
}

func (p FreshPaint) GetCompanionFlow(stats models.Stats) models.CompanionFlow {
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

func (p FreshPaint) GetActiveSkill(stats models.Stats) models.Skill {
	energy := stats.GetEnergy()
	skill := getActiveSkillForWeapon(stats.Weapon, energy)
	skill.EnemyDefenceReduction = 20
	return skill
}

func (p FreshPaint) GetBasicAttack(stats models.Stats) models.Skill {
	return getBasicAttackForWeapon(stats.Weapon)
}

func (p FreshPaint) GetOathSkill(stats models.Stats) models.Skill {
	skill := getDefaultOathSkill()
	skill.Base = 1200
	skill.AttackRate = 1600
	skill.OathBoost = stats.OathBoost
	skill.Count = getOathCount(stats)
	skill.EnemyDefenceReduction = 20
	return skill
}

func (p FreshPaint) GetResonanceSkill(stats models.Stats) models.Skill {
	skill := getDefaultResonanceSkill()
	skill.Base = 922
	skill.AttackRate = 1229
	return skill
}

func (p FreshPaint) GetSupportSkill(stats models.Stats) models.Skill {
	skill := getDefaultSupportSkill()
	skill.Base = 194
	skill.AttackRate = 258
	skill.Count = 6
	return skill
}

func (p FreshPaint) GetPassiveSkill(stats models.Stats) models.Skill {
	count := p.GetActiveSkill(stats).Count
	return models.Skill{
		Name:       "瑰色",
		Base:       188,
		AttackRate: 250,
		Count:      count + 4 + 1,
	}
}
