package companions

import "lysk-battle-record/internal/models"

type Lumiere struct{}

func (p Lumiere) GetName() string {
	return "光猎"
}

func (p Lumiere) GetCompanionFlow(stats models.Stats) models.CompanionFlow {
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

func (p Lumiere) GetActiveSkill(stats models.Stats) models.Skill {
	energy := stats.GetEnergy()
	var skill models.Skill

	if stats.Weapon == "专武" {
		skill = getDefaultActiveSkill()
		skill.Base = 403
		skill.AttackRate = 215
		skill.DefenseRate = 852
		skill.Count = energy - 8
	} else {
		skill = getActiveSkillForWeapon(stats.Weapon, energy)
	}

	if stats.SetCard == "末夜" {
		if stats.Stage == "III" || stats.Stage == "IV" {
			skill.Count += 6
		}

		if stats.Stage == "IV" {
			skill.Count += 4
		}
	}

	return skill
}

func (p Lumiere) GetBasicAttack(stats models.Stats) models.Skill {
	if stats.Weapon == "专武" {
		skill := getDefaultBasicAttack()
		skill.Base = 150
		skill.AttackRate = 80
		skill.DefenseRate = 317
		skill.Count = 35
		return skill
	}

	return getBasicAttackForWeapon(stats.Weapon)
}

func (p Lumiere) GetResonanceSkill(stats models.Stats) models.Skill {
	skill := getDefaultResonanceSkill()
	skill.Base = 686
	skill.AttackRate = 366
	skill.DefenseRate = 1450
	return skill
}

func (p Lumiere) GetOathSkill(stats models.Stats) models.Skill {
	skill := getDefaultOathSkill()
	skill.Base = 1440
	skill.AttackRate = 780
	skill.DefenseRate = 3060
	skill.OathBoost = stats.OathBoost
	skill.Count = getOathCount(stats)
	return skill
}

func (p Lumiere) GetSupportSkill(stats models.Stats) models.Skill {
	skill := getDefaultSupportSkill()
	skill.Count = 3
	return skill
}

func (p Lumiere) GetPassiveSkill(stats models.Stats) models.Skill {
	activeSkillCount := p.GetActiveSkill(stats).Count
	heavyAttackCount := p.GetBasicAttack(stats).Count
	supportSkillCount := p.GetSupportSkill(stats).Count
	partnerCount := 26 // tested

	if stats.Weapon != "专武" {
		activeSkillCount = 0
	}

	count := partnerCount + activeSkillCount + supportSkillCount + heavyAttackCount/4 + 4 // last 4 is from 共鸣
	if stats.SetCard == "末夜" && stats.Stage == "IV" {
		count = partnerCount*(60.0-4.0*8.0)/60.0 + (partnerCount*8.0*4.0/60.0)*4 + // 非朦胧期 + 朦胧期
			activeSkillCount*(60.0-4.0*8.0)/60.0 + (activeSkillCount*8.0*4.0/60.0)*4 + // 非朦胧期 + 朦胧期
			//supportSkillCount +
			(heavyAttackCount/4)*(60.0-4.0*8.0)/60.0 +
			4
	}
	passiveSkill := models.Skill{
		Name:        "月光",
		Base:        92,
		AttackRate:  49,
		DefenseRate: 194,
		Count:       count,
		CanBeCrit:   true,
	}

	return passiveSkill
}
