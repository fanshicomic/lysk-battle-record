package partners

import "lysk-battle-record/internal/models"

type LightSeeker struct{}

func (p LightSeeker) GetName() string {
	return "逐光骑士"
}

func (p LightSeeker) GetPartnerFlow(stats models.Stats) models.PartnerFlow {
	activeSkill := p.GetActiveSkill(stats)
	heavyAttack := p.GetHeavyAttack(stats)
	resonanceSkill := p.GetResonanceSkill(stats)
	oathSkill := p.GetOathSkill(stats)
	supportSkill := p.GetSupportSkill(stats)
	passiveSkill := p.GetPassiveSkill(stats)

	weakenRate := getWeakenRate(stats.Matching)

	flow := models.PartnerFlow{
		Periods: []models.PartnerPeriod{
			{
				SkillSet: models.PartnerSkillSet{
					Skills: []models.Skill{
						activeSkill,
						heavyAttack,
						resonanceSkill,
						oathSkill,
						supportSkill,
						passiveSkill,
					},
				},
				WeakenRate: weakenRate,
			},
		},
		Boost: 25, // 溯光力场内10%攻击增益+破盾后增伤20%
	}

	return flow
}

func (p LightSeeker) GetActiveSkill(stats models.Stats) models.Skill {
	energy := stats.GetEnergy()

	if stats.Weapon == "专武" {
		skill := getDefaultActiveSkill()
		skill.Base = 341
		skill.AttackRate = 455
		skill.Count = (energy - 8) * 2
		return skill
	}

	return getActiveSkillForWeapon(stats.Weapon, energy)
}

func (p LightSeeker) GetHeavyAttack(stats models.Stats) models.Skill {
	if stats.Weapon == "专武" {
		skill := getDefaultHeavyAttack()
		skill.Base = 118
		skill.AttackRate = 157
		skill.Count = 30
		return skill
	}

	return getHeavyAttackForWeapon(stats.Weapon)
}

func (p LightSeeker) GetResonanceSkill(stats models.Stats) models.Skill {
	skill := getDefaultResonanceSkill()
	skill.Base = 641
	skill.AttackRate = 854
	return skill
}

func (p LightSeeker) GetOathSkill(stats models.Stats) models.Skill {
	skill := getDefaultOathSkill()
	skill.Base = 1440
	skill.AttackRate = 1920
	skill.DamageBoost = stats.OathBoost
	skill.Count = getOathCount(stats)
	return skill
}

func (p LightSeeker) GetSupportSkill(stats models.Stats) models.Skill {
	skill := getDefaultSupportSkill()
	skill.Base = 400
	skill.AttackRate = 400
	skill.Count = 6
	return skill
}

func (p LightSeeker) GetPassiveSkill(stats models.Stats) models.Skill {
	activeSKillCount := p.GetActiveSkill(stats).Count
	skill := models.Skill{
		Name:       "溯光共鸣",
		Base:       150,
		AttackRate: 200,
		Count:      int(float64(activeSKillCount) * 0.8),
		CanBeCrit:  true,
	}

	return skill
}
