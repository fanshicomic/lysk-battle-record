package set_cards

import "lysk-battle-record/internal/models"

func GetDefaultCardBuff() map[string]models.StageBuff {
	return map[string]models.StageBuff{
		"IV": {
			Buffs: map[string]models.SkillBuff{
				"所有": {DamageBoost: 16},
			},
		},
		"III": {
			Buffs: map[string]models.SkillBuff{
				"所有": {DamageBoost: 8},
			},
		},
		"II": {
			Buffs: map[string]models.SkillBuff{
				"所有": {DamageBoost: 8},
			},
		},
		"I": {
			Buffs: map[string]models.SkillBuff{
				"所有": {DamageBoost: 8},
			},
		},
	}
}
