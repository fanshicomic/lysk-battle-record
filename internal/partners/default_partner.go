package partners

import "lysk-battle-record/internal/models"

type DefaultPartner struct{}

func (p DefaultPartner) GetName() string {
	return "默认搭档"
}

func (p DefaultPartner) GetPartnerFlow(stats models.Stats) models.PartnerFlow {
	activeSkill := p.GetActiveSkill(stats)
	heavyAttack := p.GetBasicAttack(stats)
	resonanceSkill := p.GetResonanceSkill(stats)
	oathSkill := p.GetOathSkill(stats)
	supportSkill := p.GetSupportSkill(stats)
	passiveSkill := p.GetPassiveSkill(stats)

	weakenRate := getWeakenRate(stats.Matching)
	return models.PartnerFlow{
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
				WeakenRate: weakenRate,
			},
		},
	}
}

func (p DefaultPartner) GetActiveSkill(stats models.Stats) models.Skill {
	energy := stats.GetEnergy()
	if stats.Weapon == "专武" {
		skill := getDefaultActiveSkill()
		skill.Count = energy - 8
		return skill
	}

	return getActiveSkillForWeapon(stats.Weapon, energy)
}

func (p DefaultPartner) GetBasicAttack(stats models.Stats) models.Skill {
	if stats.Weapon == "专武" {
		return getDefaultBasicAttack()
	}

	return getBasicAttackForWeapon(stats.Weapon)
}

func (p DefaultPartner) GetOathSkill(stats models.Stats) models.Skill {
	skill := getDefaultOathSkill()
	skill.OathBoost = stats.OathBoost
	skill.Count = getOathCount(stats)
	return skill
}

func (p DefaultPartner) GetResonanceSkill(stats models.Stats) models.Skill {
	skill := getDefaultResonanceSkill()
	return skill
}

func (p DefaultPartner) GetSupportSkill(stats models.Stats) models.Skill {
	skill := getDefaultSupportSkill()
	skill.Count = 6
	return skill
}

func (p DefaultPartner) GetPassiveSkill(stats models.Stats) models.Skill {
	return models.Skill{
		Name: "被动",
	}
}
