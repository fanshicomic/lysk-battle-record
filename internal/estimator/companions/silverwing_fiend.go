package companions

import "lysk-battle-record/internal/models"

type SilverwingFiend struct{}

func (p SilverwingFiend) GetName() string {
	return "银翼恶魔"
}

func (p SilverwingFiend) GetCompanionFlow(stats models.Stats) models.CompanionFlow {
	activeSkill := p.GetActiveSkill(stats)
	bloodPeriodActiveSkill := p.GetBloodPeriodActiveSkill(stats)

	heavyAttack := p.GetHeavyAttack(stats)
	bloodPeriodHeavyAttack := p.GetBloodPeriodHeavyAttack(stats)

	resonanceSkill := p.GetResonanceSkill(stats)
	bloodPeriodResonanceSkill := p.GetBloodPeriodResonanceSkill(stats)

	oathSkill := p.GetOathSkill(stats)
	roseThornsSkill := p.GetRoseThornsSkill(stats)
	bloodPeriodRoseThornsSkill := p.GetBloodPeriodRoseThornsSkill(stats)

	supportSkill := p.GetSupportSkill(stats)
	bloodPeriodSupportSkill := p.GetBloodPeriodSupportSkill(stats)

	roseBulletSkill := p.GetRoseBulletSkill(stats)

	weakenRate := getWeakenRate(stats.Matching)

	return models.CompanionFlow{
		Periods: []models.CompanionPeriod{
			{
				SkillSet: models.CompanionSkillSet{
					Skills: []models.Skill{
						activeSkill,
						heavyAttack,
						resonanceSkill,
						supportSkill,
						roseThornsSkill,
					},
				},
				WeakenRate: 0,
			},
			{
				SkillSet: models.CompanionSkillSet{
					Skills: []models.Skill{
						bloodPeriodActiveSkill,
						bloodPeriodHeavyAttack,
						bloodPeriodResonanceSkill,
						bloodPeriodSupportSkill,
						bloodPeriodRoseThornsSkill,
						roseBulletSkill,
						oathSkill,
					},
				},
				WeakenRate: weakenRate * 2,
			},
		},
	}
}

func (p SilverwingFiend) GetActiveSkill(stats models.Stats) models.Skill {
	energy := stats.GetEnergy()
	if stats.Weapon == "专武" {
		skill := getDefaultActiveSkill()
		skill.Count = 0
		return skill
	}

	skill := getActiveSkillForWeapon(stats.Weapon, energy)
	skill.Count /= 2
	return skill
}

func (p SilverwingFiend) GetBloodPeriodActiveSkill(stats models.Stats) models.Skill {
	energy := stats.GetEnergy()
	if stats.Weapon == "专武" {
		skill := getDefaultActiveSkill()
		skill.Name = "蔷薇之雨"
		skill.Base = 45 * 8
		skill.AttackRate = 24 * 8
		skill.DefenseRate = 95 * 8
		skill.Count = p.GetBloodPeriodCount()
		skill.DamageBoost = 15.0
		return skill
	}

	skill := getActiveSkillForWeapon(stats.Weapon, energy)
	skill.Count /= 2
	return skill
}

func (p SilverwingFiend) GetBasicAttack(stats models.Stats) models.Skill {
	if stats.Weapon == "专武" {
		skill := getDefaultBasicAttack()

		skill.Base = 49.0 + 49.0 + 94.0 + 146.0
		skill.AttackRate = 26.0 + 26.0 + 50.0 + 78.0
		skill.DefenseRate = 103.0 + 103.0 + 198.0 + 309.0
		skill.Count = 2 * p.GetNormalPeriodCount()
		skill.DamageBoost = 15.0

		return skill
	}

	return models.Skill{}
}

func (p SilverwingFiend) GetHeavyAttack(stats models.Stats) models.Skill {
	if stats.Weapon == "专武" {
		skill := getDefaultBasicAttack()
		skill.Name = "重击"
		skill.Base = 180.0 + 49.0
		skill.AttackRate = 96.0 + 26.0
		skill.DefenseRate = 381.0 + 103.0
		skill.Count = 7 * p.GetNormalPeriodCount()
		skill.DamageBoost = 15.0

		return skill
	}

	return getBasicAttackForWeapon(stats.Weapon)
}

func (p SilverwingFiend) GetBloodPeriodHeavyAttack(stats models.Stats) models.Skill {
	if stats.Weapon == "专武" {
		skill := getDefaultBasicAttack()
		skill.Name = "重击-血誓"
		skill.Base = 180.0
		skill.AttackRate = 96.0
		skill.DefenseRate = 381.0
		skill.Count = p.GetBloodPeriodCount()
		skill.DamageBoost = 15.0

		return skill
	}

	return getBasicAttackForWeapon(stats.Weapon)
}

func (p SilverwingFiend) GetOathSkill(stats models.Stats) models.Skill {
	skill := getDefaultOathSkill()
	skill.Base = 1800
	skill.AttackRate = 960
	skill.DefenseRate = 3820
	skill.OathBoost = stats.OathBoost
	skill.Count = getOathCount(stats)
	return skill
}

func (p SilverwingFiend) GetResonanceSkill(stats models.Stats) models.Skill {
	skill := getDefaultResonanceSkill()
	skill.Base = 1663
	skill.AttackRate = 887
	skill.DefenseRate = 3515
	skill.Count = p.GetNormalPeriodCount()
	return skill
}

func (p SilverwingFiend) GetBloodPeriodResonanceSkill(stats models.Stats) models.Skill {
	skill := getDefaultResonanceSkill()
	skill.Name = "堙界之棂"
	skill.Base = 1458
	skill.AttackRate = 778
	skill.DefenseRate = 3084
	skill.DamageBoost = 60 // 两朵血蔷薇引爆增伤
	skill.Count = p.GetBloodPeriodCount()
	skill.NoWeakenPeriod = false
	if stats.Weapon == "专武" {
		skill.DamageBoost += 15.0
	}
	return skill
}

func (p SilverwingFiend) GetSupportSkill(stats models.Stats) models.Skill {
	skill := getDefaultSupportSkill()
	skill.Base = 584.0
	skill.AttackRate = 311.0
	skill.DefenseRate = 1244.0
	skill.Count = p.GetNormalPeriodCount()
	return skill
}

func (p SilverwingFiend) GetBloodPeriodSupportSkill(stats models.Stats) models.Skill {
	skill := getDefaultSupportSkill()
	skill.Name = "血破之湮"
	skill.Base = 628.0
	skill.AttackRate = 335.0
	skill.DefenseRate = 1327.0
	skill.Count = 2 * p.GetBloodPeriodCount()
	return skill
}

func (p SilverwingFiend) GetRoseThornsSkill(stats models.Stats) models.Skill {
	return models.Skill{
		Name:        "蔷薇地棘",
		Base:        675,
		AttackRate:  360,
		DefenseRate: 1427,
		CanBeCrit:   true,
		Count:       2 * p.GetNormalPeriodCount(),
	}
}

func (p SilverwingFiend) GetBloodPeriodRoseThornsSkill(stats models.Stats) models.Skill {
	return models.Skill{
		Name:        "蔷薇地棘-血誓",
		Base:        675,
		AttackRate:  360,
		DefenseRate: 1427,
		CanBeCrit:   true,
		Count:       3 * p.GetBloodPeriodCount(),
	}
}

func (p SilverwingFiend) GetRoseBulletSkill(stats models.Stats) models.Skill {
	count := 6 * p.GetBloodPeriodCount() // 一次重击 (带一次普攻) + 三次蔷薇地棘 + 一次堙界之棂
	if stats.SetCard != "猩红" {
		count = 0
	}

	return models.Skill{
		Name:        "血蔷薇子弹",
		Base:        79,
		AttackRate:  42,
		DefenseRate: 168,
		CanBeCrit:   true,
		DamageBoost: 15,
		Count:       count, // 一次重击 + 三次蔷薇地棘 + 一次堙界之棂
	}
}

func (p SilverwingFiend) GetNormalPeriodCount() int {
	return 2
}

func (p SilverwingFiend) GetBloodPeriodCount() int {
	return 2
}
