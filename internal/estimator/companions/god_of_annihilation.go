package companions

import "lysk-battle-record/internal/models"

type GodOfAnnihilation struct{}

func (p GodOfAnnihilation) GetName() string {
	return "终末之神"
}

func (p GodOfAnnihilation) GetCompanionFlow(stats models.Stats) models.CompanionFlow {
	activeSkill := p.GetActiveSkill(stats)
	basicAttack := p.GetBasicAttack(stats)
	heavyAttack := p.GetHeavyAttack(stats)
	featherAttack := p.GetFeatherAttack(stats)
	resonanceSkill := p.GetResonanceSkill(stats)
	oathSkill := p.GetOathSkill(stats)
	supportSkill := p.GetSupportSkill(stats)

	godPeriodHeavyAttack := p.GetGodPeriodHeavyAttack(stats)
	godPeriodSupportSkill := p.GetGodPeriodSupportSkill(stats)
	soulAttack := p.GetSoulAttack(stats)

	weakenRate := getWeakenRate(stats.Matching)
	damageBoost := 0.0
	if stats.SetCard == "神谕" {
		damageBoost = 8.0
	}
	return models.CompanionFlow{
		Periods: []models.CompanionPeriod{
			{
				SkillSet: models.CompanionSkillSet{
					Skills: []models.Skill{
						activeSkill,
						basicAttack,
						heavyAttack,
						featherAttack,
						resonanceSkill,
						supportSkill,
					},
				},
				WeakenRate: 0,
			},
			{
				SkillSet: models.CompanionSkillSet{
					Skills: []models.Skill{
						activeSkill,
						godPeriodHeavyAttack,
						soulAttack,
						godPeriodSupportSkill,
						oathSkill,
					},
				},
				WeakenRate: weakenRate * 2,
				Boost:      damageBoost,
			},
		},
	}
}

func (p GodOfAnnihilation) GetActiveSkill(stats models.Stats) models.Skill {
	energy := stats.GetEnergy()
	if stats.Weapon == "专武" {
		skill := getDefaultActiveSkill()
		skill.Base = 338 + 389
		skill.AttackRate = 180 + 207
		skill.HpRate = 16.2 + 18.7
		skill.Count = p.GetNormalPeriodCount()
		return skill
	}

	skill := getActiveSkillForWeapon(stats.Weapon, energy)
	skill.Count /= 2
	return skill
}

func (p GodOfAnnihilation) GetBasicAttack(stats models.Stats) models.Skill {
	if stats.Weapon == "专武" {
		skill := getDefaultBasicAttack()

		skill.Base = 105.0 + 116.0 + 116.0 + 147.0
		skill.AttackRate = 56.0 + 62.0 + 62.0 + 78.0
		skill.HpRate = 5.1 + 5.6 + 5.6 + 7.1
		skill.Count = 2 * p.GetNormalPeriodCount()
		if stats.SetCard == "神谕" && stats.Stage == "IV" {
			skill.Base = 105.0 + 116.0
			skill.AttackRate = 56.0 + 62.0
			skill.HpRate = 5.1 + 5.6
			skill.Count = 3 * p.GetNormalPeriodCount()
		}

		return skill
	}

	return models.Skill{}
}

func (p GodOfAnnihilation) GetHeavyAttack(stats models.Stats) models.Skill {
	if stats.Weapon == "专武" {
		skill := getDefaultBasicAttack()
		skill.Name = "重击"
		skill.Base = 112 + 83*2
		skill.AttackRate = 60 + 44*2
		skill.HpRate = 5.3 + 4*2
		skill.Count = 2 * p.GetNormalPeriodCount()

		return skill
	}

	return getBasicAttackForWeapon(stats.Weapon)
}

func (p GodOfAnnihilation) GetGodPeriodHeavyAttack(stats models.Stats) models.Skill {
	if stats.Weapon == "专武" {
		skill := getDefaultBasicAttack()
		skill.Name = "重击-飞升"
		skill.Base = 112 + 83*2
		skill.AttackRate = 60 + 44*2
		skill.HpRate = 5.3 + 4*2
		skill.Count = 3 * p.GetNormalPeriodCount()

		return skill
	}

	return getBasicAttackForWeapon(stats.Weapon)
}

func (p GodOfAnnihilation) GetFeatherAttack(stats models.Stats) models.Skill {
	if stats.Weapon == "专武" {
		skill := models.Skill{
			Name:       "金箭羽",
			Base:       88*4 + 354,
			AttackRate: 47*4 + 189,
			HpRate:     4.3*4 + 17,
			CanBeCrit:  true,
			Count:      2 * p.GetNormalPeriodCount(),
		}
		return skill
	}

	return models.Skill{}
}

func (p GodOfAnnihilation) GetSoulAttack(stats models.Stats) models.Skill {
	if stats.Weapon == "专武" {
		skill := models.Skill{
			Name:       "魂隙击破",
			Base:       315,
			AttackRate: 168,
			HpRate:     15.1,
			CanBeCrit:  true,
			Count:      2 * p.GetNormalPeriodCount(),
		}
		return skill
	}

	return models.Skill{}
}

func (p GodOfAnnihilation) GetOathSkill(stats models.Stats) models.Skill {
	skill := getDefaultOathSkill()
	skill.Base = 1800
	skill.AttackRate = 960
	skill.HpRate = 86
	skill.OathBoost = stats.OathBoost
	skill.Count = getOathCount(stats)
	return skill
}

func (p GodOfAnnihilation) GetResonanceSkill(stats models.Stats) models.Skill {
	skill := getDefaultResonanceSkill()
	skill.Base = 1262
	skill.AttackRate = 674
	skill.HpRate = 60.6
	return skill
}

func (p GodOfAnnihilation) GetSupportSkill(stats models.Stats) models.Skill {
	skill := getDefaultSupportSkill()
	skill.Base = 98 + 272
	skill.AttackRate = 52 + 143
	skill.HpRate = 4.7 + 13.1
	skill.Count = 2 * p.GetNormalPeriodCount()
	return skill
}

func (p GodOfAnnihilation) GetGodPeriodSupportSkill(stats models.Stats) models.Skill {
	skill := getDefaultSupportSkill()
	skill.Name = "协助-飞升"
	skill.Base = 852
	skill.AttackRate = 456
	skill.HpRate = 40.8
	skill.Count = 2 * p.GetGodPeriodCount()
	return skill
}

func (p GodOfAnnihilation) GetPassiveSkill(stats models.Stats) models.Skill {
	return models.Skill{
		Name: "被动",
	}
}

func (p GodOfAnnihilation) GetNormalPeriodCount() int {
	return 2
}

func (p GodOfAnnihilation) GetGodPeriodCount() int {
	return 2
}
