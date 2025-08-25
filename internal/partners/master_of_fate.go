package partners

import "lysk-battle-record/internal/models"

type MasterOfFate struct{}

func (p MasterOfFate) GetName() string {
	return "九黎司命"
}

func (p MasterOfFate) GetPartnerFlow(stats models.Stats) models.PartnerFlow {
	activeSkill := p.GetActiveSkill(stats)
	heavyAttack := p.GetBasicAttack(stats)
	resonanceSkill := p.GetResonanceSkill(stats)
	oathSkill := p.GetOathSkill(stats)
	supportSkill := p.GetSupportSkill(stats)
	passiveSkill := p.GetPassiveSkill(stats)
	altPassiveSkill := p.GetAltResonanceSkill(stats)

	weakenRate := getWeakenRate(stats.Matching)
	if stats.SetCard == "拥雪" && stats.Stage != "I" && stats.Stage != "无套装" {
		weakenRate *= 1.1
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
						altPassiveSkill,
					},
				},
				WeakenRate: weakenRate,
			},
		},
	}

	return flow
}

func (p MasterOfFate) GetActiveSkill(stats models.Stats) models.Skill {
	energy := stats.GetEnergy()
	if stats.SetCard == "拥雪" && (stats.Stage == "III" || stats.Stage == "IV") {
		energy += 2
	}
	if stats.Weapon == "专武" {
		skill := getDefaultActiveSkill()
		skill.Base = 404
		skill.AttackRate = 539
		skill.Count = energy - 8
		return skill
	}

	return getActiveSkillForWeapon(stats.Weapon, energy)
}

func (p MasterOfFate) GetBasicAttack(stats models.Stats) models.Skill {
	if stats.Weapon == "专武" {
		skill := getDefaultBasicAttack()
		skill.Base = 141
		skill.AttackRate = 188
		skill.Count = 30
		return skill
	}

	return getBasicAttackForWeapon(stats.Weapon)
}

func (p MasterOfFate) GetResonanceSkill(stats models.Stats) models.Skill {
	skill := getDefaultResonanceSkill()
	skill.Base = 632
	skill.AttackRate = 842
	return skill
}

func (p MasterOfFate) GetOathSkill(stats models.Stats) models.Skill {
	skill := getDefaultOathSkill()
	skill.Base = 1440
	skill.AttackRate = 1920
	skill.OathBoost = stats.OathBoost
	skill.Count = getOathCount(stats)
	return skill
}

func (p MasterOfFate) GetSupportSkill(stats models.Stats) models.Skill {
	skill := getDefaultSupportSkill()
	skill.Base = 260
	skill.AttackRate = 348
	skill.Count = 6
	return skill
}

func (p MasterOfFate) GetPassiveSkill(stats models.Stats) models.Skill {
	activeSkillCount := p.GetActiveSkill(stats).Count
	supportSkillCount := p.GetSupportSkill(stats).Count
	altResonanceSkillCount := p.GetAltResonanceSkill(stats).Count
	partnerCount := 5
	passiveSkill := models.Skill{
		Name:       "断玉诀",
		Base:       233,
		AttackRate: 310,
		Count:      (activeSkillCount*4+supportSkillCount+altResonanceSkillCount)/3.0 + partnerCount + 6, // 6 from normal attacks
		CanBeCrit:  true,
	}

	return passiveSkill
}

func (p MasterOfFate) GetAltResonanceSkill(stats models.Stats) models.Skill {
	skill := models.Skill{
		Name:       "穿雨",
		Base:       205,
		AttackRate: 273,
		Count:      3 * 4,
		CanBeCrit:  true,
	}

	if stats.Weapon != "专武" {
		skill.Count = 0
	}

	return skill
}
