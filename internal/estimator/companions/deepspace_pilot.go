package companions

import "lysk-battle-record/internal/models"

type DeepspacePilot struct{}

func (p DeepspacePilot) GetName() string {
	return "深空飞行员"
}

func (p DeepspacePilot) GetCompanionFlow(stats models.Stats) models.CompanionFlow {
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

func (p DeepspacePilot) GetActiveSkill(stats models.Stats) models.Skill {
	energy := stats.GetEnergy()
	return getActiveSkillForWeapon(stats.Weapon, energy)
}

func (p DeepspacePilot) GetBasicAttack(stats models.Stats) models.Skill {
	return getBasicAttackForWeapon(stats.Weapon)
}

func (p DeepspacePilot) GetOathSkill(stats models.Stats) models.Skill {
	skill := getDefaultOathSkill()
	skill.Base = 1200
	skill.AttackRate = 1600
	skill.OathBoost = stats.OathBoost
	skill.Count = getOathCount(stats)
	return skill
}

func (p DeepspacePilot) GetResonanceSkill(stats models.Stats) models.Skill {
	skill := getDefaultResonanceSkill()
	skill.Base = 867
	skill.AttackRate = 1156
	return skill
}

func (p DeepspacePilot) GetSupportSkill(stats models.Stats) models.Skill {
	skill := getDefaultSupportSkill()
	skill.Base = 285
	skill.AttackRate = 381
	skill.Count = 6
	return skill
}

func (p DeepspacePilot) GetPassiveSkill(stats models.Stats) models.Skill {
	activeSKillCount := p.GetActiveSkill(stats).Count
	supportSKillCount := p.GetSupportSkill(stats).Count
	resonanceSKillCount := p.GetResonanceSkill(stats).Count
	return models.Skill{
		Name:       "靶向标记",
		Base:       53 + 113,
		AttackRate: 70 + 150,
		Count:      activeSKillCount + supportSKillCount + resonanceSKillCount,
	}
}
