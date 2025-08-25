package partners

import (
	"lysk-battle-record/internal/models"
	"math"
)

type FarspaceColonel struct{}

func (p FarspaceColonel) GetName() string {
	return "远空执舰官"
}

func (p FarspaceColonel) GetPartnerFlow(stats models.Stats) models.PartnerFlow {
	activeSkill := p.GetActiveSkill(stats)
	lightAttack := p.GetLightAttack(stats)
	lightAttackSecondPeriod := p.GetLightAttackSecondPeriod(stats)
	heavyAttack := p.GetHeavyAttack(stats)
	resonanceSkill := p.GetResonanceSkill(stats)
	resonanceAltSkill := p.GetResonanceAltSkill(stats)
	oathSkill := p.GetOathSkill(stats)
	supportSkill := p.GetSupportSkill(stats)
	passiveSkill := p.GetPassiveSkill(stats)
	fireAltSkill := p.GetFireAltSkill(stats)

	if stats.Weapon != "专武" {
		return models.PartnerFlow{
			Periods: []models.PartnerPeriod{
				{
					SkillSet: models.PartnerSkillSet{
						Skills: []models.Skill{
							activeSkill,
							heavyAttack,
							oathSkill,
							supportSkill,
						},
					},
					WeakenRate: 0,
				},
			},
		}
	}

	// 阵地外
	activeSkill.Count -= 1
	firstPeriod := models.PartnerPeriod{
		SkillSet: models.PartnerSkillSet{
			Skills: []models.Skill{
				activeSkill,
				lightAttack,
				resonanceSkill,
				supportSkill,
				passiveSkill,
			},
		},
		WeakenRate: 0,
	}

	// 阵地内
	activeSkill.Count += 2
	boost := 0
	if stats.SetCard == "远空" {
		boost = 20
	}
	secondPeriod := models.PartnerPeriod{
		SkillSet: models.PartnerSkillSet{
			Skills: []models.Skill{
				activeSkill,
				heavyAttack,
				lightAttackSecondPeriod,
				resonanceAltSkill,
				oathSkill,
				supportSkill,
				passiveSkill,
				fireAltSkill,
			},
		},
		WeakenRate: getWeakenRate(stats.Matching) * 2,
		Boost:      float64(boost),
	}

	return models.PartnerFlow{
		Periods: []models.PartnerPeriod{
			firstPeriod,
			secondPeriod,
		},
	}
}

func (p FarspaceColonel) GetActiveSkill(stats models.Stats) models.Skill {
	energy := stats.GetEnergy()
	if stats.Weapon == "专武" {
		skill := getDefaultActiveSkill()
		skill.Base = 200 + 185
		skill.AttackRate = 105 + 100
		skill.DefenseRate = 420 + 395
		skill.Count = int(math.Ceil(float64(energy) / (3.0 * 2.0)))
		return skill
	}

	return getActiveSkillForWeapon(stats.Weapon, energy)
}

func (p FarspaceColonel) GetHeavyAttack(stats models.Stats) models.Skill {
	if stats.Weapon == "专武" {
		skill := getDefaultHeavyAttack()
		skill.Base = 133
		skill.AttackRate = 71
		skill.DefenseRate = 281
		skill.Count = 2 * p.GetSecondPeriodCount(stats)
		return skill
	}

	return getHeavyAttackForWeapon(stats.Weapon)
}

func (p FarspaceColonel) GetLightAttack(stats models.Stats) models.Skill {
	if stats.Weapon == "专武" {
		// 4段普攻攒15点火力值， 共需6轮普攻触发一次共鸣
		skill := models.Skill{
			Name:        "普攻",
			Base:        47 + 70 + 68 + 85,
			AttackRate:  25 + 37 + 36 + 45,
			DefenseRate: 99 + 148 + 144 + 180,
			CanBeCrit:   true,
			Count:       6 * p.GetFirstPeriodCount(stats),
		}
		return skill
	}
	return models.Skill{}
}

func (p FarspaceColonel) GetLightAttackSecondPeriod(stats models.Stats) models.Skill {
	if stats.Weapon == "专武" {
		skill := models.Skill{
			Name:        "普攻",
			Base:        68 + 85,
			AttackRate:  36 + 45,
			DefenseRate: 144 + 180,
			CanBeCrit:   true,
			Count:       (3 + 3) * p.GetSecondPeriodCount(stats), // first 3 is after each heavy attach / active attack, other 3 is extra 3 hits
		}
		return skill
	}
	return models.Skill{}
}

func (p FarspaceColonel) GetOathSkill(stats models.Stats) models.Skill {
	skill := getDefaultOathSkill()
	skill.Base = 1440
	skill.AttackRate = 780
	skill.DefenseRate = 3060
	skill.DamageBoost = stats.OathBoost
	skill.Count = getOathCount(stats)
	if stats.Weapon != "专武" {
		skill.Name = "非专武誓约(无虚弱期)"
		skill.NoWeakenPeriod = false
	}
	return skill
}

func (p FarspaceColonel) GetResonanceSkill(stats models.Stats) models.Skill {
	skill := getDefaultResonanceSkill()
	skill.Base = 245
	skill.AttackRate = 131
	skill.DefenseRate = 519
	skill.Count = p.GetSecondPeriodCount(stats)
	return skill
}

func (p FarspaceColonel) GetResonanceAltSkill(stats models.Stats) models.Skill {
	skill := getDefaultResonanceSkill()
	skill.Name = "纵深打击"
	skill.Base = 512
	skill.AttackRate = 273
	skill.DefenseRate = 1082
	skill.DamageBoost = 80
	skill.Count = p.GetSecondPeriodCount(stats)
	return skill
}

func (p FarspaceColonel) GetSupportSkill(stats models.Stats) models.Skill {
	skill := getDefaultSupportSkill()
	skill.Base = 284
	skill.AttackRate = 151
	skill.DefenseRate = 599
	skill.Count = 3
	return skill
}

func (p FarspaceColonel) GetPassiveSkill(stats models.Stats) models.Skill {
	return models.Skill{
		Name: "被动",
	}
}

func (p FarspaceColonel) GetFireAltSkill(stats models.Stats) models.Skill {
	if stats.SetCard == "远空" && stats.Stage == "IV" {
		skill := models.Skill{
			Name:        "集火引力波",
			Base:        225,
			AttackRate:  120,
			DefenseRate: 476,
			CanBeCrit:   true,
			Count:       4 * p.GetSecondPeriodCount(stats),
		}
		return skill
	}
	return models.Skill{
		Name: "集火引力波",
	}
}

func (p FarspaceColonel) GetFirstPeriodCount(stats models.Stats) int {
	if stats.Weapon != "专武" {
		return 4
	}
	return 2
}

func (p FarspaceColonel) GetSecondPeriodCount(stats models.Stats) int {
	if stats.Weapon != "专武" {
		return 0
	}
	return 2
}
