package partners

import "lysk-battle-record/internal/models"

type AbyssWalker struct{}

func (a AbyssWalker) GetName() string {
	return "深海潜行者"
}

func (a AbyssWalker) GetPartnerFlow(stats models.Stats) models.PartnerFlow {
	activeSkill := a.GetActiveSkill(stats)
	heavyAttack := a.GetHeavyAttack(stats)
	resonanceSkill := a.GetResonanceSkill()
	oathSkill := a.GetOathSkill(stats)
	supportSkill := a.GetSupportSkill()

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
		CritRate:   a.getExtraCritRate(),
		CanBeCrit:  true,
	}
	if stats.SetCard == "深海" && stats.Stage == "IV" {
		passiveSkillSlash.Count += 2
	}

	weakenRate := getWeakenRate(stats.Matching)

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
						passiveSkillBurn,
						passiveSkillSlash,
					},
				},
			},
		},
		// 4*0.8: 潜能回复时的攻击增益参数
		Boost:      8 * 0.7,
		WeakenRate: weakenRate,
	}

	return flow
}

func (a AbyssWalker) GetActiveSkill(stats models.Stats) models.Skill {
	energy := stats.GetEnergy()

	if stats.Weapon == "专武" {
		bonusCount := 0
		if stats.SetCard == "深海" && (stats.Stage == "III" || stats.Stage == "IV") {
			bonusCount = 2
		}

		return models.Skill{
			Name:       "主动",
			Base:       309,
			AttackRate: 412,
			Count:      energy - 8 + bonusCount,
			CritRate:   a.getExtraCritRate(),
			CanBeCrit:  true,
		}
	}

	return getActiveSkillForWeapon(stats.Weapon, energy)
}

func (a AbyssWalker) GetHeavyAttack(stats models.Stats) models.Skill {
	if stats.Weapon == "专武" {
		return models.Skill{
			Name:       "重击",
			Base:       144,
			AttackRate: 192,
			Count:      35,
			CritRate:   a.getExtraCritRate(),
			CanBeCrit:  true,
		}
	}

	return getHeavyAttackForWeapon(stats.Weapon)
}

func (a AbyssWalker) GetResonanceSkill() models.Skill {
	resonanceSkill := models.Skill{
		Name:       "共鸣",
		Base:       785,
		AttackRate: 1047,
		Count:      4,
		CritRate:   a.getExtraCritRate(),
		CanBeCrit:  true,
	}

	return resonanceSkill
}

func (a AbyssWalker) GetOathSkill(stats models.Stats) models.Skill {
	oathSkill := models.Skill{
		Name:        "誓约",
		Base:        1440,
		AttackRate:  1920,
		DamageBoost: stats.OathBoost * 100,
		Count:       getOathCount(stats),
	}

	return oathSkill
}

func (a AbyssWalker) GetSupportSkill() models.Skill {
	supportSkill := models.Skill{
		Name:       "协助",
		Base:       264,
		AttackRate: 352,
		Count:      4,
		CanBeCrit:  true,
	}

	return supportSkill
}

func (a AbyssWalker) getExtraCritRate() float64 {
	// 2: 专武普攻被动参数 15*2*5/60: 潜能充满时增加15%暴击持续5秒，60秒可触发两次，平均到60秒内的增益
	critRate := float64(2 + 15*2*5/60)
	return critRate
}
