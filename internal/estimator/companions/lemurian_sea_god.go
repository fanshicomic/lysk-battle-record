package companions

import "lysk-battle-record/internal/models"

type LemurianSeaGod struct{}

func (p LemurianSeaGod) GetName() string {
	return "利莫里亚海神"
}

func (p LemurianSeaGod) GetCompanionFlow(stats models.Stats) models.CompanionFlow {
	activeSkill := p.GetActiveSkill(stats)
	basicAttack := p.GetBasicAttack(stats)
	boostedHeavyAttack := p.GetBoostedHeavyAttack(stats)
	resonanceSkill := p.GetResonanceSkill(stats)
	oathSkill := p.GetOathSkill(stats)
	supportSkill := p.GetSupportSkill(stats)

	godPeriodActiveSkill := p.GetGodPeriodActiveSkill(stats)
	godPeriodBasicAttack := p.GetGodPeriodBasicAttack(stats)
	godPeriodBoostedHeavyAttack := p.GetGodPeriodBoostedHeavyAttack(stats)
	godPeriodThunder := p.GetPartnerThunderSkill(stats)

	thunderBallSkill := p.GetThunderBallSkill(stats)
	thunderWaveSkill := p.GetThunderWaveSkill(stats)

	weakenRate := getWeakenRate(stats.Matching)
	if stats.SetCard == "雾海" {
		weakenRate *= 1.1
	}
	return models.CompanionFlow{
		Periods: []models.CompanionPeriod{
			{
				SkillSet: models.CompanionSkillSet{
					Skills: []models.Skill{
						activeSkill,
						basicAttack,
						boostedHeavyAttack,
						resonanceSkill,
						supportSkill,
					},
				},
				WeakenRate: 0,
			},
			{
				SkillSet: models.CompanionSkillSet{
					Skills: []models.Skill{
						godPeriodActiveSkill,
						godPeriodBasicAttack,
						godPeriodBoostedHeavyAttack,
						oathSkill,
						thunderBallSkill,
						thunderWaveSkill,
						godPeriodThunder,
					},
				},
				WeakenRate: weakenRate,
			},
		},
	}
}

func (p LemurianSeaGod) GetActiveSkill(stats models.Stats) models.Skill {
	energy := stats.GetEnergy()
	if stats.Weapon == "专武" {
		skill := getDefaultActiveSkill()
		skill.Base = 312
		skill.AttackRate = 166
		skill.DefenseRate = 660
		skill.Count = p.GetNormalPeriodCount()
		return skill
	}

	skill := getActiveSkillForWeapon(stats.Weapon, energy)
	skill.Count /= 2
	return skill
}

func (p LemurianSeaGod) GetGodPeriodActiveSkill(stats models.Stats) models.Skill {
	energy := stats.GetEnergy()
	if stats.Weapon == "专武" {
		skill := getDefaultActiveSkill()
		skill.Name = "主动-神眷"
		skill.Base = 312
		skill.AttackRate = 166
		skill.DefenseRate = 660
		skill.Count = p.GetGodPeriodCount() * 3
		return skill
	}

	skill := getActiveSkillForWeapon(stats.Weapon, energy)
	skill.Count /= 2
	return skill
}

func (p LemurianSeaGod) GetBasicAttack(stats models.Stats) models.Skill {
	if stats.Weapon == "专武" {
		skill := getDefaultBasicAttack()
		skill.Base = (49 + 53 + 83) / 3
		skill.AttackRate = (26 + 28 + 44) / 3
		skill.DefenseRate = (104 + 112 + 175) / 3
		skill.Count = 6 * p.GetNormalPeriodCount()

		return skill
	}

	return getBasicAttackForWeapon(stats.Weapon)
}

func (p LemurianSeaGod) GetGodPeriodBasicAttack(stats models.Stats) models.Skill {
	if stats.Weapon == "专武" {
		skill := getDefaultBasicAttack()
		skill.Name = "普攻-神眷"
		skill.Base = (49 + 53 + 83) / 3
		skill.AttackRate = (26 + 28 + 44) / 3
		skill.DefenseRate = (104 + 112 + 175) / 3
		skill.Count = 9 * p.GetGodPeriodCount()

		return skill
	}

	return getBasicAttackForWeapon(stats.Weapon)
}

func (p LemurianSeaGod) GetBoostedHeavyAttack(stats models.Stats) models.Skill {
	activeSkillCount := p.GetActiveSkill(stats).Count
	if stats.Weapon == "专武" {
		return models.Skill{
			Name:        "武器被动重击",
			CanBeCrit:   true,
			Base:        265,
			AttackRate:  141,
			DefenseRate: 560,
			Count:       activeSkillCount * 2,
			DamageBoost: 50,
		}
	}
	return models.Skill{}
}

func (p LemurianSeaGod) GetGodPeriodBoostedHeavyAttack(stats models.Stats) models.Skill {
	if stats.Weapon == "专武" {
		return models.Skill{
			Name:        "武器被动重击-神眷",
			CanBeCrit:   true,
			Base:        265,
			AttackRate:  141,
			DefenseRate: 560,
			Count:       3 * p.GetGodPeriodCount(),
			DamageBoost: 50,
		}
	}
	return models.Skill{}
}

func (p LemurianSeaGod) GetOathSkill(stats models.Stats) models.Skill {
	skill := getDefaultOathSkill()
	skill.Base = 1800
	skill.AttackRate = 960
	skill.DefenseRate = 3820
	skill.OathBoost = stats.OathBoost
	skill.Count = getOathCount(stats)
	return skill
}

func (p LemurianSeaGod) GetResonanceSkill(stats models.Stats) models.Skill {
	skill := getDefaultResonanceSkill()
	skill.Base = 1311
	skill.AttackRate = 699
	skill.DefenseRate = 2773
	return skill
}

func (p LemurianSeaGod) GetSupportSkill(stats models.Stats) models.Skill {
	skill := getDefaultSupportSkill()
	skill.Base = 360
	skill.AttackRate = 192
	skill.DefenseRate = 761
	skill.Count = 4
	return skill
}

func (p LemurianSeaGod) GetPassiveSkill(stats models.Stats) models.Skill {
	return models.Skill{
		Name: "被动",
	}
}

func (p LemurianSeaGod) GetThunderBallSkill(stats models.Stats) models.Skill {
	return models.Skill{
		Name:        "雷晶",
		CanBeCrit:   true,
		Base:        32,
		AttackRate:  17,
		DefenseRate: 68,
		Count:       3 * p.GetGodPeriodCount(),
	}
}

func (p LemurianSeaGod) GetThunderWaveSkill(stats models.Stats) models.Skill {
	return models.Skill{
		Name:        "雷潮",
		CanBeCrit:   true,
		Base:        450,
		AttackRate:  240,
		DefenseRate: 951,
		Count:       3 * p.GetGodPeriodCount(),
	}
}

func (p LemurianSeaGod) GetPartnerThunderSkill(stats models.Stats) models.Skill {
	skill := models.Skill{
		Name:        "落雷",
		CanBeCrit:   true,
		Base:        270,
		AttackRate:  144,
		DefenseRate: 571,
		Count:       10 * p.GetGodPeriodCount(),
	}

	if stats.SetCard != "雾海" || (stats.Stage == "I" || stats.Stage == "无套装") {
		skill.Count = 0
	}

	return skill
}

func (p LemurianSeaGod) GetNormalPeriodCount() int {
	return 2
}

func (p LemurianSeaGod) GetGodPeriodCount() int {
	return 2
}
