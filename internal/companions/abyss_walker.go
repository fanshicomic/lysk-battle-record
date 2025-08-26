package companions

import "lysk-battle-record/internal/models"

type AbyssWalker struct{}

func (p AbyssWalker) GetName() string {
	return "深海潜行者"
}

func (p AbyssWalker) GetCompanionFlow(stats models.Stats) models.CompanionFlow {
	activeSkill := p.GetActiveSkill(stats)
	heavyAttack := p.GetBasicAttack(stats)
	resonanceSkill := p.GetResonanceSkill(stats)
	oathSkill := p.GetOathSkill(stats)
	supportSkill := p.GetSupportSkill(stats)

	passiveSkillBurn := models.Skill{
		Name:       "灼烧",
		Base:       23,
		AttackRate: 31,
		Count:      4 * 7,
	}
	if stats.Weapon == "专武" {
		passiveSkillBurn.Count += activeSkill.Count * 7
	}

	passiveSkillSlash := models.Skill{
		Name:       "强力斩击",
		Base:       540,
		AttackRate: 720,
		Count:      2,
		CritRate:   p.getExtraCritRate(stats),
		CanBeCrit:  true,
	}
	if stats.SetCard == "深海" && stats.Stage == "IV" {
		passiveSkillSlash.Count += 2
	}

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
						passiveSkillBurn,
						passiveSkillSlash,
					},
				},
				WeakenRate: weakenRate,
				Boost:      8 * 0.7, // 4*0.8: 潜能回复时的攻击增益参数
			},
		},
	}

	return flow
}

func (p AbyssWalker) GetActiveSkill(stats models.Stats) models.Skill {
	energy := stats.GetEnergy()

	if stats.Weapon == "专武" {
		bonusCount := 0
		if stats.SetCard == "深海" && (stats.Stage == "III" || stats.Stage == "IV") {
			bonusCount = 2
		}

		skill := getDefaultActiveSkill()
		skill.Base = 309
		skill.AttackRate = 412
		skill.Count = energy - 8 + bonusCount
		skill.CritRate = p.getExtraCritRate(stats)
		return skill
	}

	return getActiveSkillForWeapon(stats.Weapon, energy)
}

func (p AbyssWalker) GetBasicAttack(stats models.Stats) models.Skill {
	if stats.Weapon == "专武" {
		skill := getDefaultBasicAttack()
		skill.Base = 144
		skill.AttackRate = 192
		skill.Count = 35
		skill.CritRate = p.getExtraCritRate(stats)
		return skill
	}

	return getBasicAttackForWeapon(stats.Weapon)
}

func (p AbyssWalker) GetResonanceSkill(stats models.Stats) models.Skill {
	skill := getDefaultResonanceSkill()
	skill.Base = 785
	skill.AttackRate = 1047
	skill.CritRate = p.getExtraCritRate(stats)
	return skill
}

func (p AbyssWalker) GetOathSkill(stats models.Stats) models.Skill {
	skill := getDefaultOathSkill()
	skill.Base = 1440
	skill.AttackRate = 1920
	skill.OathBoost = stats.OathBoost
	skill.Count = getOathCount(stats)
	return skill
}

func (p AbyssWalker) GetSupportSkill(stats models.Stats) models.Skill {
	skill := getDefaultSupportSkill()
	skill.Base = 264
	skill.AttackRate = 352
	skill.Count = 4
	return skill
}

func (p AbyssWalker) getExtraCritRate(stats models.Stats) float64 {
	// 2: 专武普攻被动参数 15*2*5/60: 潜能充满时增加15%暴击持续5秒，60秒可触发两次，平均到60秒内的增益
	critRate := float64(2 + 15*2*5/60)

	if stats.Weapon != "专武" {
		critRate -= 2
	}

	return critRate
}
