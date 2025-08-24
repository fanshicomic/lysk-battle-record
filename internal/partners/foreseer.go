package partners

import "lysk-battle-record/internal/models"

type Foreseer struct{}

func (p Foreseer) GetName() string {
	return "永恒先知"
}

func (p Foreseer) GetPartnerFlow(stats models.Stats) models.PartnerFlow {
	activeSkill := p.GetActiveSkill(stats)
	heavyAttack := p.GetHeavyAttack(stats)
	resonanceSkill := p.GetResonanceSkill(stats)
	oathSkill := p.GetOathSkill(stats)
	supportSkill := p.GetSupportSkill()
	passiveSkill := p.GetPassiveSkill(stats)

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
						passiveSkill,
					},
				},
			},
		},
		WeakenRate: weakenRate,
	}

	return flow
}

func (p Foreseer) GetActiveSkill(stats models.Stats) models.Skill {
	energy := stats.GetEnergy()

	if stats.Weapon == "专武" {
		return models.Skill{
			Name:        "主动",
			Base:        52,
			AttackRate:  28,
			DefenseRate: 111,
			Count:       (energy - 8) * 6,
			CanBeCrit:   true,
		}
	}

	return getActiveSkillForWeapon(stats.Weapon, energy)
}

func (p Foreseer) GetHeavyAttack(stats models.Stats) models.Skill {
	lightAttackDamageAdjustment := 0.0
	if stats.Stage != "IV" { // 非满阶时轻重击交替，调整伤害参数
		lightAttackDamageAdjustment = -20.0
	}
	if stats.Weapon == "专武" {
		return models.Skill{
			Name:        "重击",
			Base:        167,
			AttackRate:  89,
			DefenseRate: 353,
			Count:       25,
			DamageBoost: lightAttackDamageAdjustment,
			CanBeCrit:   true,
		}
	}

	return getHeavyAttackForWeapon(stats.Weapon)
}

func (p Foreseer) GetResonanceSkill(stats models.Stats) models.Skill {
	resonanceSkill := models.Skill{
		Name:        "共鸣",
		Count:       4,
		Base:        790,
		AttackRate:  421,
		DefenseRate: 1670,
		CanBeCrit:   true,
	}

	return resonanceSkill
}

func (p Foreseer) GetOathSkill(stats models.Stats) models.Skill {
	oathSkill := models.Skill{
		Name:        "誓约",
		Base:        1440,
		AttackRate:  780,
		DefenseRate: 3060,
		DamageBoost: stats.OathBoost,
		Count:       getOathCount(stats),
	}

	return oathSkill
}

func (p Foreseer) GetSupportSkill() models.Skill {
	supportSkill := models.Skill{
		Name: "协助",
	}

	return supportSkill
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
