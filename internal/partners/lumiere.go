package partners

import "lysk-battle-record/internal/models"

type Lumiere struct{}

func (p Lumiere) GetName() string {
	return "光猎"
}

func (p Lumiere) GetPartnerFlow(stats models.Stats) models.PartnerFlow {
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

func (p Lumiere) GetActiveSkill(stats models.Stats) models.Skill {
	energy := stats.GetEnergy()
	var skill models.Skill
	if stats.Weapon == "专武" {
		skill = models.Skill{
			Name:        "主动",
			Base:        403,
			AttackRate:  215,
			DefenseRate: 852,
			Count:       energy - 8,
			CanBeCrit:   true,
		}
	} else {
		skill = getActiveSkillForWeapon(stats.Weapon, energy)
	}

	if stats.SetCard == "末夜" {
		if stats.Stage == "III" || stats.Stage == "IV" {
			skill.Count += 3
		}

		if stats.Stage == "IV" {
			skill.Count += 4
		}
	}

	return skill
}

func (p Lumiere) GetHeavyAttack(stats models.Stats) models.Skill {
	if stats.Weapon == "专武" {
		return models.Skill{
			Name:        "重击",
			Base:        150,
			AttackRate:  80,
			DefenseRate: 317,
			Count:       35,
			CanBeCrit:   true,
		}
	}

	return getHeavyAttackForWeapon(stats.Weapon)
}

func (p Lumiere) GetResonanceSkill(stats models.Stats) models.Skill {
	resonanceSkill := models.Skill{
		Name:        "共鸣",
		Base:        686,
		AttackRate:  366,
		DefenseRate: 1450,
		Count:       4,
		CanBeCrit:   true,
	}

	return resonanceSkill
}

func (p Lumiere) GetOathSkill(stats models.Stats) models.Skill {
	oathSkill := models.Skill{
		Name:        "誓约",
		Base:        1440,
		AttackRate:  780,
		DefenseRate: 3060,
		DamageBoost: stats.OathBoost * 100,
		Count:       getOathCount(stats),
	}

	return oathSkill
}

func (p Lumiere) GetSupportSkill() models.Skill {
	supportSkill := models.Skill{
		Name:  "协助",
		Count: 3,
	}

	return supportSkill
}

func (p Lumiere) GetPassiveSkill(stats models.Stats) models.Skill {
	activeSkillCount := p.GetActiveSkill(stats).Count
	heavyAttackCount := p.GetHeavyAttack(stats).Count
	supportSkillCount := p.GetSupportSkill().Count
	partnerCount := 26

	if stats.Weapon != "专武" {
		activeSkillCount = 0
	}

	count := partnerCount + activeSkillCount + supportSkillCount + heavyAttackCount/4 + 4 // last 4 is from 共鸣
	if stats.SetCard == "末夜" && stats.Stage == "IV" {
		count = partnerCount + (8 * 4 / 60) +
			activeSkillCount + activeSkillCount*(8*4/60)*4 +
			supportSkillCount +
			heavyAttackCount/4 + 3*4 + // last 3 * 4 is heavy attack during 朦胧 period
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
