package set_cards

import "lysk-battle-record/internal/models"

type LightSeeking struct{}

func (c LightSeeking) GetName() string {
	return "逐光"
}

func (c LightSeeking) GetSetCardBuff() map[string]models.StageBuff {
	return map[string]models.StageBuff{
		"IV": {
			Buffs: map[string]models.SkillBuff{
				"所有":       {DamageBoost: 16},
				"主动":       {DamageBoost: 25, CountBonus: 1.4},
				"重剑主动":   {DamageBoost: 25},
				"单手剑主动": {DamageBoost: 25},
				"法杖主动":   {DamageBoost: 25},
				"手枪主动":   {DamageBoost: 25},
			},
		},
		"III": {
			Buffs: map[string]models.SkillBuff{
				"所有":       {DamageBoost: 8},
				"主动":       {DamageBoost: 25},
				"重剑主动":   {DamageBoost: 25},
				"单手剑主动": {DamageBoost: 25},
				"法杖主动":   {DamageBoost: 25},
				"手枪主动":   {DamageBoost: 25},
			},
		},
		"II": {
			Buffs: map[string]models.SkillBuff{
				"所有":       {DamageBoost: 8},
				"主动":       {DamageBoost: 25},
				"重剑主动":   {DamageBoost: 25},
				"单手剑主动": {DamageBoost: 25},
				"法杖主动":   {DamageBoost: 25},
				"手枪主动":   {DamageBoost: 25},
			},
		},
		"I": {
			Buffs: map[string]models.SkillBuff{
				"所有":       {DamageBoost: 8},
				"主动":       {DamageBoost: 25},
				"重剑主动":   {DamageBoost: 25},
				"单手剑主动": {DamageBoost: 25},
				"法杖主动":   {DamageBoost: 25},
				"手枪主动":   {DamageBoost: 25},
			},
		},
	}
}
