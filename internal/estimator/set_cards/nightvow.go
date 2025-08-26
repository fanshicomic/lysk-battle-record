package set_cards

import "lysk-battle-record/internal/models"

type Nightvow struct{}

func (c Nightvow) GetName() string {
	return "夜誓"
}

func (c Nightvow) GetSetCardBuff() map[string]models.StageBuff {
	return map[string]models.StageBuff{
		"IV": {
			Buffs: map[string]models.SkillBuff{
				"所有":      {DamageBoost: 16},
				"主动-加冕": {DamageBoost: 40, CountBonus: 2, DefenceReduction: 12.5},
				"普攻-加冕": {DefenceReduction: 12.5},
				"协助-加冕": {DefenceReduction: 12.5, DamageBoost: 24},
			},
		},
		"III": {
			Buffs: map[string]models.SkillBuff{
				"所有":      {DamageBoost: 8},
				"主动-加冕": {DamageBoost: 40, DefenceReduction: 12.5},
				"普攻-加冕": {DefenceReduction: 12.5},
				"协助-加冕": {DefenceReduction: 12.5},
			},
		},
		"II": {
			Buffs: map[string]models.SkillBuff{
				"所有":      {DamageBoost: 8},
				"主动-加冕": {DamageBoost: 40},
			},
		},
		"I": {
			Buffs: map[string]models.SkillBuff{
				"所有": {DamageBoost: 8},
			},
		},
	}
}
