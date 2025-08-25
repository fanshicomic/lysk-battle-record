package set_cards

import "lysk-battle-record/internal/models"

type Abyssal struct{}

func (c Abyssal) GetName() string {
	return "无套装"
}

func (c Abyssal) GetSetCardBuff() map[string]models.StageBuff {
	return map[string]models.StageBuff{
		"IV": {
			Buffs: map[string]models.SkillBuff{
				"所有":         {DamageBoost: 16},
				"主动":         {DamageBoost: 30},
				"协助":         {CountBonus: 1.34},
				"魔魇之爪范围": {CountBonus: 1.34},
			},
		},
		"III": {
			Buffs: map[string]models.SkillBuff{
				"所有":         {DamageBoost: 8},
				"协助":         {CountBonus: 1.34},
				"魔魇之爪范围": {CountBonus: 1.34},
			},
		},
		"II": {
			Buffs: map[string]models.SkillBuff{
				"所有":         {DamageBoost: 8},
				"协助":         {CountBonus: 1.34},
				"魔魇之爪范围": {CountBonus: 1.34},
			},
		},
		"I": {
			Buffs: map[string]models.SkillBuff{
				"所有":         {DamageBoost: 8},
				"魔魇之爪范围": {CountBonus: 0},
			},
		},
	}
}
