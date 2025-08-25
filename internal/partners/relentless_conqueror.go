package partners

import "lysk-battle-record/internal/models"

type RelentLessConqueror struct{}

func (p RelentLessConqueror) GetName() string {
	return "无尽掠夺者"
}

func (p RelentLessConqueror) GetPartnerFlow(stats models.Stats) models.PartnerFlow {
	activeSkill := p.GetActiveSkill(stats)
	heavyAttack := p.GetHeavyAttack(stats)
	resonanceSkill := p.GetResonanceSkill(stats)
	oathSkill := p.GetOathSkill(stats)
	supportSkill := p.GetSupportSkill(stats)
	passiveSkill := p.GetPassiveSkill(stats)

	weakenRate := getWeakenRate(stats.Matching)
	boost := (4.0 * 8.0 / 60.0) * 80.0
	if stats.SetCard == "掠心" && stats.Stage != "I" {
		boost = 80
	}

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
		Boost: boost,
	}

	return flow
}

func (p RelentLessConqueror) GetActiveSkill(stats models.Stats) models.Skill {
	energy := stats.GetEnergy()

	if stats.Weapon == "专武" {
		skill := getDefaultActiveSkill()
		skill.Base = 342
		skill.AttackRate = 456
		skill.Count = energy - 8
		skill.DamageBoost = 10.0 / (10.0 - 6.0*0.5)
		return skill
	}

	return getActiveSkillForWeapon(stats.Weapon, energy)
}

func (p RelentLessConqueror) GetHeavyAttack(stats models.Stats) models.Skill {
	if stats.Weapon == "专武" {
		skill := getDefaultHeavyAttack()
		skill.Base = 160
		skill.AttackRate = 213
		skill.Count = 20
		return skill
	}

	return getHeavyAttackForWeapon(stats.Weapon)
}

func (p RelentLessConqueror) GetResonanceSkill(stats models.Stats) models.Skill {
	skill := getDefaultResonanceSkill()
	skill.Base = 1094
	skill.AttackRate = 1458
	return skill
}

func (p RelentLessConqueror) GetOathSkill(stats models.Stats) models.Skill {
	skill := getDefaultOathSkill()
	skill.Base = 1440
	skill.AttackRate = 1920
	skill.DamageBoost = stats.OathBoost
	skill.Count = getOathCount(stats)
	return skill
}

func (p RelentLessConqueror) GetSupportSkill(stats models.Stats) models.Skill {
	skill := getDefaultSupportSkill()
	skill.Base = 239
	skill.AttackRate = 318
	skill.Count = 6
	return skill
}

func (p RelentLessConqueror) GetPassiveSkill(stats models.Stats) models.Skill {
	skill := models.Skill{
		Name:       "掠噬标记",
		Base:       60,
		AttackRate: 80,
		Count:      60 * 2 / 8,
		CanBeCrit:  true,
	}

	if stats.SetCard != "掠心" || stats.Stage != "IV" {
		skill.Count = 0
	}

	return skill
}
