package companions

import "lysk-battle-record/internal/models"

type Netherlord struct{}

func (p Netherlord) GetName() string {
	return "冥罗之主"
}

func (p Netherlord) GetCompanionFlow(stats models.Stats) models.CompanionFlow {
	activeSkill := p.GetActiveSkill(stats)
	yinYangPeriodActiveSkill := p.GetYinYangPeriodActiveSkill(stats)

	basicAttack := p.GetBasicAttack(stats)
	heavyAttack := p.GetHeavyAttack(stats)
	yinYangPeriodHeavyAttack := p.GetYinYangPeriodHeavyAttack(stats)

	resonanceSkill := p.GetResonanceSkill(stats)

	oathSkill := p.GetOathSkill(stats)

	supportSkill := p.GetSupportSkill(stats)
	yinYangPeriodSupportSkill := p.GetYinYangPeriodSupportSkill(stats)

	lotusSkill := p.GetLotusSkill(stats)
	yinYangPeriodLotusSkill := p.GetYinYangPeriodLotusSkill(stats)

	yinYangGapSkill := p.GetYinYangGapSkill(stats)
	ghostClawSkill := p.GetGhostClawSkill(stats)
	yinYangHeavyAttackExtraDamage := p.GetYinYangHeavyAttackExtraDamage(stats)
	ghostClawExtraDamage := p.GetGhostClawExtraDamage(stats)

	weakenRate := getWeakenRate(stats.Matching)

	return models.CompanionFlow{
		Periods: []models.CompanionPeriod{
			{
				SkillSet: models.CompanionSkillSet{
					Skills: []models.Skill{
						activeSkill,
						basicAttack,
						heavyAttack,
						resonanceSkill,
						supportSkill,
						lotusSkill,
					},
				},
				WeakenRate: 0,
			},
			{
				SkillSet: models.CompanionSkillSet{
					Skills: []models.Skill{
						yinYangPeriodActiveSkill,
						yinYangPeriodHeavyAttack,
						yinYangPeriodSupportSkill,
						oathSkill,
						yinYangPeriodLotusSkill,
						yinYangGapSkill,
						ghostClawSkill,
						yinYangHeavyAttackExtraDamage,
						ghostClawExtraDamage,
					},
				},
				WeakenRate: weakenRate * 2,
			},
		},
	}
}

func (p Netherlord) GetActiveSkill(stats models.Stats) models.Skill {
	energy := stats.GetEnergy()
	if stats.Weapon == "专武" {
		skill := getDefaultActiveSkill()
		skill.Name = "断念舞"
		skill.Base = 358.0
		skill.AttackRate = 190.0
		skill.HpRate = 17.2
		skill.Count = 2 * p.GetNormalPeriodCount()
		return skill
	}

	skill := getActiveSkillForWeapon(stats.Weapon, energy)
	skill.Count /= 2
	return skill
}

func (p Netherlord) GetYinYangPeriodActiveSkill(stats models.Stats) models.Skill {
	energy := stats.GetEnergy()
	if stats.Weapon == "专武" {
		skill := getDefaultActiveSkill()
		skill.Name = "碧空谣"
		skill.Base = 671.0
		skill.AttackRate = 358.0
		skill.HpRate = 32.2
		skill.Count = 3 * p.GetYinYangPeriodCount()
		return skill
	}

	skill := getActiveSkillForWeapon(stats.Weapon, energy)
	skill.Count /= 2
	return skill
}

func (p Netherlord) GetBasicAttack(stats models.Stats) models.Skill {
	if stats.Weapon == "专武" {
		skill := getDefaultBasicAttack()

		skill.Base = 74.0 + 66.0 + 135.0 + 135.0
		skill.AttackRate = 39.0 + 34.0 + 72.0 + 72.0
		skill.HpRate = 3.5 + 3.2 + 6.5 + 6.5
		skill.Count = p.GetNormalPeriodCount() // 感觉这个排轴，两次共鸣期间算上主动，打一轮顶天了

		return skill
	}

	return models.Skill{}
}

func (p Netherlord) GetHeavyAttack(stats models.Stats) models.Skill {
	if stats.Weapon == "专武" {
		skill := getDefaultBasicAttack()
		skill.Name = "重击"
		skill.Base = 174.0
		skill.AttackRate = 93.0
		skill.HpRate = 8.4
		skill.Count = 2 * p.GetNormalPeriodCount()

		return skill
	}

	return getBasicAttackForWeapon(stats.Weapon)
}

func (p Netherlord) GetYinYangPeriodHeavyAttack(stats models.Stats) models.Skill {
	if stats.Weapon == "专武" {
		skill := getDefaultBasicAttack()
		skill.Name = "重击-阴阳"
		skill.Base = 174.0
		skill.AttackRate = 93.0
		skill.HpRate = 8.4
		skill.Count = 3 * p.GetYinYangPeriodCount() // 仅三阶有，记得在日卡非三阶里删掉

		return skill
	}

	return getBasicAttackForWeapon(stats.Weapon)
}

func (p Netherlord) GetOathSkill(stats models.Stats) models.Skill {
	skill := getDefaultOathSkill()
	skill.Name = "誓约-尘露不忘"
	skill.Base = 1800
	skill.AttackRate = 960
	skill.HpRate = 86
	skill.OathBoost = stats.OathBoost
	skill.Count = getOathCount(stats) // 三阶两倍，二阶1.8倍。。服了这个怎么算
	if stats.SetCard == "沉冥" && stats.Stage == "IV" {
		skill.CanBeCrit = true
	}
	return skill
}

func (p Netherlord) GetResonanceSkill(stats models.Stats) models.Skill {
	skill := getDefaultResonanceSkill()
	skill.Base = 1361
	skill.AttackRate = 725
	skill.HpRate = 65.3
	skill.Count = p.GetNormalPeriodCount()
	return skill
}

func (p Netherlord) GetSupportSkill(stats models.Stats) models.Skill {
	skill := getDefaultSupportSkill()
	skill.Name = "镇狱斩"
	skill.Base = 336.0
	skill.AttackRate = 179.0
	skill.HpRate = 16.2
	skill.Count = 2 * p.GetNormalPeriodCount()
	return skill
}

func (p Netherlord) GetYinYangPeriodSupportSkill(stats models.Stats) models.Skill {
	skill := getDefaultSupportSkill()
	skill.Name = "幽冥斩"
	skill.Base = 296.0
	skill.AttackRate = 158.0
	skill.HpRate = 14.2
	skill.Count = 3 * p.GetYinYangPeriodCount()
	return skill
}

func (p Netherlord) GetLotusSkill(stats models.Stats) models.Skill {
	return models.Skill{
		Name:       "并蒂莲",
		Base:       491,
		AttackRate: 262,
		HpRate:     23.6,
		CanBeCrit:  true,
		Count:      2 * p.GetNormalPeriodCount(), // 一阶以上多1
	}
}

func (p Netherlord) GetYinYangPeriodLotusSkill(stats models.Stats) models.Skill {
	return models.Skill{
		Name:       "并蒂莲-阴阳",
		Base:       491,
		AttackRate: 262,
		HpRate:     23.6,
		CanBeCrit:  true,
		Count:      p.GetYinYangPeriodCount(), // 有点微妙，感觉也是三阶才能打出来一个半花，因为其他时候主动中间的3s不够打完一整轮
	}
}

func (p Netherlord) GetYinYangGapSkill(stats models.Stats) models.Skill {
	return models.Skill{
		Name:       "阴阳之隙",
		Base:       645,
		AttackRate: 344,
		HpRate:     31,
		CanBeCrit:  true,
		Count:      3 * p.GetYinYangPeriodCount(),
	}
}

func (p Netherlord) GetGhostClawSkill(stats models.Stats) models.Skill {
	return models.Skill{
		Name:       "鬼爪",
		Base:       327 + 251*2, // 2的部分是恶鬼重复标记的伤害
		AttackRate: 175 + 134*2,
		HpRate:     15.7 + 12.1*2,
		CanBeCrit:  true,
		Count:      p.GetYinYangPeriodCount(),
	}
}

func (p Netherlord) GetYinYangHeavyAttackExtraDamage(stats models.Stats) models.Skill {
	return models.Skill{
		Name:       "三阶阴阳重击额外伤害",
		Base:       327,
		AttackRate: 175,
		HpRate:     15.7,
		CanBeCrit:  true,
		Count:      3 * p.GetYinYangPeriodCount(),
	}
}

func (p Netherlord) GetGhostClawExtraDamage(stats models.Stats) models.Skill {
	return models.Skill{
		Name:       "三阶额外鬼爪",
		Base:       251,
		AttackRate: 134,
		HpRate:     12.1,
		CanBeCrit:  true,
		Count:      3 * p.GetYinYangPeriodCount(),
	}
}

func (p Netherlord) GetNormalPeriodCount() int {
	return 2
}

func (p Netherlord) GetYinYangPeriodCount() int {
	return 2
}
