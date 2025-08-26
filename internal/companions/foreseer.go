package companions

import "lysk-battle-record/internal/models"

type Foreseer struct{}

func (p Foreseer) GetName() string {
	return "永恒先知"
}

func (p Foreseer) GetCompanionFlow(stats models.Stats) models.CompanionFlow {
	activeSkill := p.GetActiveSkill(stats)
	heavyAttack := p.GetBasicAttack(stats)
	resonanceSkill := p.GetResonanceSkill(stats)
	oathSkill := p.GetOathSkill(stats)
	supportSkill := p.GetSupportSkill(stats)
	passiveSkill := p.GetPassiveSkill(stats)

	weakenRate := getWeakenRate(stats.Matching)

	flow := models.CompanionFlow{
		Periods: []models.CompanionPeriod{
			{
				SkillSet: models.CompanionSkillSet{
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
	}

	return flow
}

func (p Foreseer) GetActiveSkill(stats models.Stats) models.Skill {
	energy := stats.GetEnergy()

	if stats.Weapon == "专武" {
		skill := getDefaultActiveSkill()
		skill.Base = 52
		skill.AttackRate = 28
		skill.DefenseRate = 111
		skill.Count = (energy - 8) * 6
		return skill
	}

	return getActiveSkillForWeapon(stats.Weapon, energy)
}

func (p Foreseer) GetBasicAttack(stats models.Stats) models.Skill {
	if stats.Weapon == "专武" {
		lightAttackDamageAdjustment := 0.0
		if stats.Stage != "IV" {
			lightAttackDamageAdjustment = -20.0
		}

		skill := getDefaultBasicAttack()
		skill.Base = 167
		skill.AttackRate = 89
		skill.DefenseRate = 353
		skill.Count = 25
		skill.DamageBoost = lightAttackDamageAdjustment
		return skill
	}

	return getBasicAttackForWeapon(stats.Weapon)
}

func (p Foreseer) GetResonanceSkill(stats models.Stats) models.Skill {
	skill := getDefaultResonanceSkill()
	skill.Base = 790
	skill.AttackRate = 421
	skill.DefenseRate = 1670
	return skill
}

func (p Foreseer) GetOathSkill(stats models.Stats) models.Skill {
	skill := getDefaultOathSkill()
	skill.Base = 1440
	skill.AttackRate = 780
	skill.DefenseRate = 3060
	skill.OathBoost = stats.OathBoost
	skill.Count = getOathCount(stats)
	return skill
}

func (p Foreseer) GetSupportSkill(stats models.Stats) models.Skill {
	return getDefaultSupportSkill()
}

func (p Foreseer) GetPassiveSkill(stats models.Stats) models.Skill {
	passiveSkill := models.Skill{
		Name:        "恒之罪",
		Base:        198,
		AttackRate:  102,
		DefenseRate: 406,
		Count:       12,
		CanBeCrit:   true,
	}

	if stats.SetCard == "永恒" {
		switch stats.Stage {
		case "IV":
			passiveSkill.Count += 10
			passiveSkill.Count += 4
		case "II", "III":
			passiveSkill.Count += 4
		}
	}

	if stats.SetCard == "永恒" && stats.Stage == "IV" {
		passiveSkill.Count += 10
	}

	if stats.Weapon != "专武" {
		passiveSkill.Count -= p.GetActiveSkill(stats).Count
	}

	return passiveSkill

}
