package set_cards

import "lysk-battle-record/internal/models"

type Diviner struct{}

func (c Diviner) GetName() string {
	return "神谕"
}

func (c Diviner) GetSetCardBuff() map[string]models.StageBuff {
	return map[string]models.StageBuff{
		"IV": {
			Buffs: map[string]models.SkillBuff{
				"所有":      {DamageBoost: 16 + 30*0.8},
				"协助-飞升": {CountBonus: 1.5, DamageBoost: 28.4},
				"魂隙击破":  {CountBonus: 1.5},
				"主动":      {DamageBoost: 20},
				"重击":      {CountBonus: 2},
				"金箭羽":    {CountBonus: 2.5},
			},
		},
		"III": {
			Buffs: map[string]models.SkillBuff{
				"所有":      {DamageBoost: 8},
				"协助-飞升": {CountBonus: 1.5, DamageBoost: 28.4},
				"主动":      {DamageBoost: 20},
			},
		},
		"II": {
			Buffs: map[string]models.SkillBuff{
				"所有":      {DamageBoost: 8},
				"协助-飞升": {CountBonus: 1.5, DamageBoost: 28.4},
			},
		},
		"I": {
			Buffs: map[string]models.SkillBuff{
				"所有": {DamageBoost: 8},
			},
		},
	}
}
