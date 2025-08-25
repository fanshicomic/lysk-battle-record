package set_cards

import "lysk-battle-record/internal/models"

type LoneRoad struct{}

func (c LoneRoad) GetName() string {
	return "寂路"
}

func (c LoneRoad) GetSetCardBuff() map[string]models.StageBuff {
	return map[string]models.StageBuff{
		"IV": {
			Buffs: map[string]models.SkillBuff{
				"所有":          {DamageBoost: 16},
				"誓约-同频觉醒": {OathBoost: 20, DamageBoost: 2 + 8.0/12.0},
				"誓约-同频攻击": {OathBoost: 20, DamageBoost: 32.8 + 2 + 8.0/12.0, CountBonus: 1.67},
				"普攻":          {DamageBoost: 10},
				"主动":          {DamageBoost: 10},
				"协助":          {DamageBoost: 10},
				"共鸣":          {DamageBoost: 2},
			},
		},
		"III": {
			Buffs: map[string]models.SkillBuff{
				"所有":          {DamageBoost: 8},
				"誓约-同频觉醒": {OathBoost: 20, DamageBoost: 2 + 8.0/12.0},
				"誓约-同频攻击": {OathBoost: 20, DamageBoost: 2 + 8.0/12.0},
				"普攻":          {DamageBoost: 10},
				"主动":          {DamageBoost: 10},
				"协助":          {DamageBoost: 10},
				"共鸣":          {DamageBoost: 2},
			},
		},
		"II": {
			Buffs: map[string]models.SkillBuff{
				"所有":          {DamageBoost: 8},
				"誓约-同频觉醒": {OathBoost: 20, DamageBoost: 8.0 / 12.0},
				"誓约-同频攻击": {OathBoost: 20, DamageBoost: 8.0 / 12.0},
				"普攻":          {DamageBoost: 10},
				"主动":          {DamageBoost: 10},
				"协助":          {DamageBoost: 10},
			},
		},
		"I": {
			Buffs: map[string]models.SkillBuff{
				"所有":          {DamageBoost: 8},
				"誓约-同频觉醒": {OathBoost: 20},
				"誓约-同频攻击": {OathBoost: 20},
			},
		},
	}
}
