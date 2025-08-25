package partners

import "lysk-battle-record/internal/models"

type UltimateWeaponX02 struct{}

func (p UltimateWeaponX02) GetName() string {
	return "终极兵器X-02"
}

func (p UltimateWeaponX02) GetPartnerFlow(stats models.Stats) models.PartnerFlow {
	activeSkill := p.GetActiveSkill(stats)
	basicAttack := p.GetBasicAttack(stats)
	resonanceSkill := p.GetResonanceSkill(stats)
	oathSkill := p.GetOathSkill(stats)
	oathExtraSKill := p.GetOathExtraSkill(stats)
	supportSkill := p.GetSupportSkill(stats)

	weakenRate := getWeakenRate(stats.Matching)
	if stats.SetCard != "寂路" || stats.Stage == "I" || stats.Stage == "无套装" {
		weakenRate *= 0.75
	}
	return models.PartnerFlow{
		Periods: []models.PartnerPeriod{
			{
				SkillSet: models.PartnerSkillSet{
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
				SkillSet: models.PartnerSkillSet{
					Skills: []models.Skill{
						oathSkill,
						oathExtraSKill,
					},
				},
				WeakenRate: weakenRate * 2,
			},
		},
	}
}

func (p UltimateWeaponX02) GetActiveSkill(stats models.Stats) models.Skill {
	energy := stats.GetEnergy()
	if stats.Weapon == "专武" {
		skill := getDefaultActiveSkill()
		skill.Base = 72
		skill.AttackRate = 96
		skill.Count = 2 * 12
		skill.DamageBoost = 20
		return skill
	}

	return getActiveSkillForWeapon(stats.Weapon, energy)
}

func (p UltimateWeaponX02) GetBasicAttack(stats models.Stats) models.Skill {
	if stats.Weapon == "专武" {
		skill := getDefaultBasicAttack()
		skill.Base = 77 + 74 + 109 + 137
		skill.AttackRate = 103 + 99 + 145 + 182
		skill.Count = 3 * 2
		skill.DamageBoost = 10
		return skill
	}

	return getBasicAttackForWeapon(stats.Weapon)
}

func (p UltimateWeaponX02) GetOathSkill(stats models.Stats) models.Skill {
	skill := getDefaultOathSkill()
	skill.Name = "誓约-同频觉醒"
	skill.Base = 3800
	skill.AttackRate = 5000
	skill.DamageBoost = 5.0 / 12.0
	skill.OathBoost = stats.OathBoost
	skill.CanBeCrit = true
	skill.Count = 2

	if stats.SetCard != "寂路" || stats.Stage == "I" || stats.Stage == "无套装" {
		skill.DamageBoost *= 0.75
	}

	if stats.Weapon != "专武" {
		skill.Count = 1
	}
	return skill
}

func (p UltimateWeaponX02) GetOathExtraSkill(stats models.Stats) models.Skill {
	oathSkillCount := p.GetOathSkill(stats).Count
	skill := getDefaultOathSkill()
	skill.Name = "誓约-同频攻击"
	skill.Base = 380
	skill.AttackRate = 500
	skill.DamageBoost = 5.0 / 12.0
	skill.OathBoost = stats.OathBoost
	skill.CanBeCrit = true
	skill.Count = 3 * oathSkillCount

	if stats.SetCard != "寂路" || stats.Stage == "I" || stats.Stage == "无套装" {
		skill.DamageBoost *= 0.75
	}
	return skill
}

func (p UltimateWeaponX02) GetResonanceSkill(stats models.Stats) models.Skill {
	skill := getDefaultResonanceSkill()
	skill.Base = 990
	skill.AttackRate = 1322
	if stats.SetCard != "寂路" || stats.Stage == "I" || stats.Stage == "无套装" {
		skill.Count -= 1
	}
	return skill
}

func (p UltimateWeaponX02) GetSupportSkill(stats models.Stats) models.Skill {
	skill := getDefaultSupportSkill()
	skill.Base = 461
	skill.AttackRate = 615
	skill.Count = 4
	skill.DamageBoost = 10
	return skill
}

func (p UltimateWeaponX02) GetPassiveSkill(stats models.Stats) models.Skill {
	return models.Skill{
		Name: "被动",
	}
}
