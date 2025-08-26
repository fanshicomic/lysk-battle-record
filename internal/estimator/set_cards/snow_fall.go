package set_cards

import "lysk-battle-record/internal/models"

type SnowFall struct{}

func (c SnowFall) GetName() string {
	return "拥雪"
}

func (c SnowFall) GetSetCardBuff() map[string]models.StageBuff {
	return map[string]models.StageBuff{
		"IV": {
			Buffs: map[string]models.SkillBuff{
				"所有":   {DamageBoost: 16, DefenceReduction: 10, WeakenBoost: 5},
				"誓约":   {WeakenBoost: 5},
				"断玉诀": {DamageBoost: 100},
				"穿雨":   {CountBonus: 1.34},
			},
		},
		"III": {
			Buffs: map[string]models.SkillBuff{
				"所有": {DamageBoost: 8, DefenceReduction: 10, WeakenBoost: 5},
				"誓约": {WeakenBoost: 5},
			},
		},
		"II": {
			Buffs: map[string]models.SkillBuff{
				"所有": {DamageBoost: 8, DefenceReduction: 10, WeakenBoost: 5},
				"誓约": {DamageBoost: 5},
			},
		},
		"I": {
			Buffs: map[string]models.SkillBuff{
				"所有": {DamageBoost: 8, DefenceReduction: 10},
			},
		},
	}
}
