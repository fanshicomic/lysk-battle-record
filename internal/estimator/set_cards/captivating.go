package set_cards

import "lysk-battle-record/internal/models"

type Captivating struct{}

func (c Captivating) GetName() string {
	return "掠心"
}

func (c Captivating) GetSetCardBuff() map[string]models.StageBuff {
	return map[string]models.StageBuff{
		"IV": {
			Buffs: map[string]models.SkillBuff{
				"所有":       {DamageBoost: 8.0 + 4.0},
				"主动":       {CountBonus: 1.6},
				"重剑主动":   {CountBonus: 1.6},
				"单手剑主动": {CountBonus: 1.6},
				"法杖主动":   {CountBonus: 1.6},
				"手枪主动":   {CountBonus: 1.6},
			},
		},
		"III": {
			Buffs: map[string]models.SkillBuff{
				"所有":       {DamageBoost: 8.0 + 4.0},
				"主动":       {CountBonus: 1.6},
				"重剑主动":   {CountBonus: 1.6},
				"单手剑主动": {CountBonus: 1.6},
				"法杖主动":   {CountBonus: 1.6},
				"手枪主动":   {CountBonus: 1.6},
			},
		},
		"II": {
			Buffs: map[string]models.SkillBuff{
				"所有": {DamageBoost: 8.0 + 4.0},
			},
		},
		"I": {
			Buffs: map[string]models.SkillBuff{
				"所有": {DamageBoost: 8.0 + 4.0},
			},
		},
	}
}
