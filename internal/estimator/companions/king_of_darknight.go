package companions

import "lysk-battle-record/internal/models"

type KingOfDarknight struct{}

func (p KingOfDarknight) GetName() string {
	return "暗蚀国王"
}

func (p KingOfDarknight) GetCompanionFlow(stats models.Stats) models.CompanionFlow {
	activeSkill := p.GetActiveSkill(stats)
	basicAttack := p.GetBasicAttack(stats)
	resonanceSkill := p.GetResonanceSkill(stats)
	oathSkill := p.GetOathSkill(stats)
	supportSkill := p.GetSupportSkill(stats)

	lordPeriodActiveSkill := p.GetLordPeriodActiveSkill(stats)
	lordPeriodBasicAttack := p.GetLordPeriodBasicAttack(stats)
	lordPeriodSupportSkill := p.GetLordPeriodSupportSkill(stats)

	weakenRate := getWeakenRate(stats.Matching)
	return models.CompanionFlow{
		Periods: []models.CompanionPeriod{
			{
				SkillSet: models.CompanionSkillSet{
					Skills: []models.Skill{
						activeSkill,
						basicAttack,
						resonanceSkill,
						supportSkill,
					},
				},
				WeakenRate: 0,
			},
			{
				SkillSet: models.CompanionSkillSet{
					Skills: []models.Skill{
						lordPeriodActiveSkill,
						lordPeriodBasicAttack,
						lordPeriodSupportSkill,
						oathSkill,
					},
				},
				WeakenRate: weakenRate * 2,
			},
		},
	}
}

func (p KingOfDarknight) GetActiveSkill(stats models.Stats) models.Skill {
	energy := stats.GetEnergy()
	if stats.Weapon == "专武" {
		skill := getDefaultActiveSkill()
		skill.Base = (351*2 + 376*3) / 5
		skill.AttackRate = (187*2 + 200*3) / 5
		skill.HpRate = (16.8*2 + 18*3) / 5
		skill.Count = 5 * p.GetNormalPeriodCount()
		return skill
	}

	skill := getActiveSkillForWeapon(stats.Weapon, energy)
	skill.Count /= 2
	return skill
}

func (p KingOfDarknight) GetLordPeriodActiveSkill(stats models.Stats) models.Skill {
	energy := stats.GetEnergy()
	if stats.Weapon == "专武" {
		skill := getDefaultActiveSkill()
		skill.Name = "主动-加冕"
		skill.Base = 520
		skill.AttackRate = 277
		skill.HpRate = 25
		skill.Count = p.GetLordPeriodCount() * 3
		return skill
	}

	skill := getActiveSkillForWeapon(stats.Weapon, energy)
	skill.Count /= 2
	return skill
}

func (p KingOfDarknight) GetBasicAttack(stats models.Stats) models.Skill {
	if stats.Weapon == "专武" {
		skill := getDefaultBasicAttack()
		skill.Base = 165
		skill.AttackRate = 88
		skill.HpRate = 7.9
		skill.Count = 2 * p.GetNormalPeriodCount()

		return skill
	}

	return getBasicAttackForWeapon(stats.Weapon)
}

func (p KingOfDarknight) GetLordPeriodBasicAttack(stats models.Stats) models.Skill {
	if stats.Weapon == "专武" {
		skill := getDefaultBasicAttack()
		skill.Name = "普攻-加冕"
		skill.Base = 165
		skill.AttackRate = 88
		skill.HpRate = 7.9
		skill.Count = 3 * p.GetLordPeriodCount()

		return skill
	}

	return getBasicAttackForWeapon(stats.Weapon)
}

func (p KingOfDarknight) GetOathSkill(stats models.Stats) models.Skill {
	skill := getDefaultOathSkill()
	skill.Base = 1800
	skill.AttackRate = 960
	skill.HpRate = 86
	skill.OathBoost = stats.OathBoost
	skill.Count = getOathCount(stats)
	return skill
}

func (p KingOfDarknight) GetResonanceSkill(stats models.Stats) models.Skill {
	skill := getDefaultResonanceSkill()
	skill.Base = 1767
	skill.AttackRate = 942
	skill.HpRate = 84
	return skill
}

func (p KingOfDarknight) GetSupportSkill(stats models.Stats) models.Skill {
	skill := getDefaultSupportSkill()
	skill.Base = 364
	skill.AttackRate = 194
	skill.HpRate = 17
	skill.Count = 2 * p.GetNormalPeriodCount()
	return skill
}

func (p KingOfDarknight) GetLordPeriodSupportSkill(stats models.Stats) models.Skill {
	baseCount := 4
	if stats.SetCard != "夜誓" {
		baseCount = 1
	}
	skill := getDefaultSupportSkill()
	skill.Name = "协助-加冕"
	skill.Base = 720
	skill.AttackRate = 384
	skill.HpRate = 35
	skill.Count = baseCount * p.GetNormalPeriodCount()
	return skill
}

func (p KingOfDarknight) GetPassiveSkill(stats models.Stats) models.Skill {
	return models.Skill{
		Name: "被动",
	}
}

func (p KingOfDarknight) GetNormalPeriodCount() int {
	return 2
}

func (p KingOfDarknight) GetLordPeriodCount() int {
	return 2
}
