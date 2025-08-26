package companions

import (
	"math"

	"lysk-battle-record/internal/models"
)

func getActiveSkillForWeapon(weapon string, energy int) models.Skill {
	activeSkill := models.Skill{
		Name:      "主动",
		Count:     energy - 8,
		CanBeCrit: true,
	}

	switch weapon {
	case "重剑":
		activeSkill.Name = "重剑主动"
		activeSkill.Base = 621
		activeSkill.AttackRate = 829
		activeSkill.DamageBoost = 50
		activeSkill.Count = int(math.Min(float64(energy-8), 6))
	case "单手剑":
		activeSkill.Name = "单手剑主动"
		activeSkill.Base = 341
		activeSkill.AttackRate = 455
	case "法杖":
		activeSkill.Name = "法杖主动"
		activeSkill.Base = 204
		activeSkill.AttackRate = 270
	case "手枪":
		activeSkill.Name = "手枪主动"
		activeSkill.Base = 160
		activeSkill.AttackRate = 213
	}

	return activeSkill
}

func getBasicAttackForWeapon(weapon string) models.Skill {
	heavyAttack := models.Skill{
		Name:      "普攻",
		Count:     30,
		CanBeCrit: true,
	}

	switch weapon {
	case "重剑":
		heavyAttack.Base = 337
		heavyAttack.AttackRate = 449
		heavyAttack.DamageBoost = 26
		heavyAttack.Count = 11
	case "单手剑":
		heavyAttack.Base = 250
		heavyAttack.AttackRate = 333
		heavyAttack.DamageBoost = 14
	case "法杖":
		heavyAttack.Base = 122
		heavyAttack.AttackRate = 162
		heavyAttack.DamageBoost = 28
		heavyAttack.Count = 15
	case "手枪":
		heavyAttack.Base = 120
		heavyAttack.AttackRate = 160
		heavyAttack.DamageBoost = 25
		heavyAttack.Count = 35
	}

	return heavyAttack
}

func getOathCount(stats models.Stats) int {
	if stats.Stage != "无套装" && stats.Stage != "I" {
		return 1
	}
	if stats.OathRegen >= 17 {
		return 1
	}

	return 0
}

func getWeakenRate(matching string) float64 {
	if matching == "顺" {
		return 0.50
	}
	return 0.25
}

func getDefaultActiveSkill() models.Skill {
	return models.Skill{
		Name:      "主动",
		CanBeCrit: true,
	}
}

func getDefaultBasicAttack() models.Skill {
	return models.Skill{
		Name:      "普攻",
		CanBeCrit: true,
	}
}

func getDefaultResonanceSkill() models.Skill {
	return models.Skill{
		Name:           "共鸣",
		Count:          4,
		CanBeCrit:      true,
		NoWeakenPeriod: true,
	}
}

func getDefaultOathSkill() models.Skill {
	return models.Skill{
		Name: "誓约",
	}
}

func getDefaultSupportSkill() models.Skill {
	return models.Skill{
		Name:      "协助",
		CanBeCrit: true,
	}
}
