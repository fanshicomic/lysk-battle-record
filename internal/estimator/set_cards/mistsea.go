package set_cards

import "lysk-battle-record/internal/models"

type Mistsea struct{}

func (c Mistsea) GetName() string {
	return "雾海"
}

func (c Mistsea) GetSetCardBuff() map[string]models.StageBuff {
	return map[string]models.StageBuff{
		"IV": {
			Buffs: map[string]models.SkillBuff{
				"所有":              {DamageBoost: 16},
				"武器被动重击":      {CountBonus: 1.5, DamageBoost: 50},
				"武器被动重击-神眷": {CountBonus: 1.5, DamageBoost: 50},
				"主动-神眷":         {CountBonus: 1.34},
				"雷晶":              {CountBonus: 1.34},
				"雷潮":              {CountBonus: 1.34, DamageBoost: 70},
				"普攻-神眷":         {CountBonus: 0.88},
			},
		},
		"III": {
			Buffs: map[string]models.SkillBuff{
				"所有":              {DamageBoost: 8},
				"武器被动重击":      {CountBonus: 1.5, DamageBoost: 50},
				"武器被动重击-神眷": {CountBonus: 1.5, DamageBoost: 50},
			},
		},
		"II": {
			Buffs: map[string]models.SkillBuff{
				"所有":              {DamageBoost: 8},
				"武器被动重击":      {CountBonus: 1.5},
				"武器被动重击-神眷": {CountBonus: 1.5},
			},
		},
		"I": {
			Buffs: map[string]models.SkillBuff{
				"所有":              {DamageBoost: 8},
				"武器被动重击":      {CountBonus: 1.5},
				"武器被动重击-神眷": {CountBonus: 1.5},
			},
		},
	}
}
