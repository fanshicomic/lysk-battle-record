package companions

import "lysk-battle-record/internal/models"

type AbysmSovereign struct{}

func (p AbysmSovereign) GetName() string {
	return "深渊主宰"
}

func (p AbysmSovereign) GetCompanionFlow(stats models.Stats) models.CompanionFlow {
	activeSkill := p.GetActiveSkill(stats)
	heavyAttack := p.GetBasicAttack(stats)
	resonanceSkill := p.GetResonanceSkill(stats)
	resonanceAltSkill := p.GetAltResonanceSkill(stats)
	oathSkill := p.GetOathSkill(stats)
	supportSkill := p.GetSupportSkill(stats)
	supportExtraSkill := p.GetSupportExtra(stats)
	passiveSkill := p.GetPassiveSkill(stats)

	weakenRate := getWeakenRate(stats.Matching)
	return models.CompanionFlow{
		Periods: []models.CompanionPeriod{
			{
				SkillSet: models.CompanionSkillSet{
					Skills: []models.Skill{
						activeSkill,
						heavyAttack,
						resonanceSkill,
						resonanceAltSkill,
						oathSkill,
						supportSkill,
						supportExtraSkill,
						passiveSkill,
					},
				},
				WeakenRate: weakenRate,
			},
		},
	}
}

func (p AbysmSovereign) GetActiveSkill(stats models.Stats) models.Skill {
	if stats.SetCard == "深渊" && (stats.Stage == "IV" || stats.Stage == "III") {
		stats.EnergyRegen += 24
	}
	energy := stats.GetEnergy()
	if stats.Weapon == "专武" {
		skill := getDefaultActiveSkill()
		skill.Base = 397
		skill.AttackRate = 212
		skill.HpRate = 19.1
		skill.Count = energy - 8
		skill.DamageBoost = p.getExtraBuff("主动", stats)

		if stats.SetCard == "深渊" && (stats.Stage == "IV" || stats.Stage == "III") {
			skill.Count += 3
		}

		if stats.SetCard == "深渊" && stats.Stage == "IV" {
			skill.Base = 632
			skill.AttackRate = 337
			skill.HpRate = 30.3
		}

		return skill
	}

	return getActiveSkillForWeapon(stats.Weapon, energy)
}

func (p AbysmSovereign) GetBasicAttack(stats models.Stats) models.Skill {
	if stats.Weapon == "专武" {
		return p.GetLightAttack(stats)
	}

	skill := getBasicAttackForWeapon(stats.Weapon)
	skill.DamageBoost = p.getExtraBuff("普攻", stats)
	return skill
}

func (p AbysmSovereign) GetLightAttack(stats models.Stats) models.Skill {
	activeSkillCount := p.GetActiveSkill(stats).Count
	if stats.Weapon == "专武" {
		skill := getDefaultBasicAttack()
		skill.Name = "普攻"
		skill.Base = 162
		skill.AttackRate = 87
		skill.HpRate = 7.8
		skill.Count = 11*4 - activeSkillCount*2
		skill.DamageBoost = p.getExtraBuff("普攻", stats)
		return skill
	}

	return getBasicAttackForWeapon(stats.Weapon)
}

func (p AbysmSovereign) GetOathSkill(stats models.Stats) models.Skill {
	oathSkill := getDefaultOathSkill()
	oathSkill.Base = 1440
	oathSkill.AttackRate = 780
	oathSkill.HpRate = 69.4
	oathSkill.OathBoost = stats.OathBoost
	oathSkill.Count = getOathCount(stats)

	return oathSkill
}

func (p AbysmSovereign) GetResonanceSkill(stats models.Stats) models.Skill {
	resonanceSkill := getDefaultResonanceSkill()
	resonanceSkill.Base = 512
	resonanceSkill.AttackRate = 273
	resonanceSkill.HpRate = 24.6

	return resonanceSkill
}

// 二段伤害，第二段可暴击可虚弱
func (p AbysmSovereign) GetAltResonanceSkill(stats models.Stats) models.Skill {
	resonanceSkill := getDefaultResonanceSkill()
	resonanceSkill.Name = "共鸣下半"
	resonanceSkill.Base = 512
	resonanceSkill.AttackRate = 273
	resonanceSkill.HpRate = 24.6
	resonanceSkill.NoWeakenPeriod = false

	return resonanceSkill
}

func (p AbysmSovereign) GetSupportSkill(stats models.Stats) models.Skill {
	supportSkill := getDefaultSupportSkill()
	supportSkill.Base = 508
	supportSkill.AttackRate = 271
	supportSkill.HpRate = 24.4
	supportSkill.Count = 6
	supportSkill.DamageBoost = p.getExtraBuff("协助", stats)

	return supportSkill
}

func (p AbysmSovereign) GetPassiveSkill(stats models.Stats) models.Skill {
	return models.Skill{
		Name:  "被动",
		Count: 0,
	}
}

func (p AbysmSovereign) GetSupportExtra(stats models.Stats) models.Skill {
	supportSkillCount := p.GetSupportSkill(stats).Count
	skill := getDefaultSupportSkill()
	skill.Name = "魔魇之爪范围"
	skill.Base = 121
	skill.HpRate = 5.8
	skill.AttackRate = 64.8
	skill.Count = supportSkillCount * 3

	return skill
}

func (p AbysmSovereign) getExtraBuff(skill string, stats models.Stats) float64 {
	// 74 to 46 for 8 times, first 2 assigns to active and support rest 7 to light
	// 74 70 66 62 58 54 50 46
	if skill == "主动" {
		if stats.SetCard == "深渊" && stats.Stage == "IV" {
			return 74.0
		}

		return 74.0 / 2.0 // 非三阶时只有一半的时候使用渊怒
	}

	if skill == "协助" {
		if stats.SetCard == "深渊" && stats.Stage == "IV" {
			return 70.0
		}

		return 70.0 / 2.0 // 非三阶时只有一半的时候使用渊怒
	}
	if stats.SetCard == "深渊" && stats.Stage == "IV" {
		return (66.0 + 46.0) * 7.0 / (2.0 * 11.0) // only 7 out of 11 can enjoy buff
	}

	return (66.0 + 46.0) * 7.0 / (4.0 * 11.0)
}
