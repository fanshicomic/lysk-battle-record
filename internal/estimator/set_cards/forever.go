package set_cards

import "lysk-battle-record/internal/models"

type Forever struct{}

func (c Forever) GetName() string {
	return "永恒"
}

func (c Forever) GetSetCardBuff() map[string]models.StageBuff {
	return map[string]models.StageBuff{
		"IV": {
			Buffs: map[string]models.SkillBuff{
				"所有": {DamageBoost: 16 +
					0.8*12 + // 咒文守护下12%增伤，持续12s，冷却15s
					0.6*10}, // 共鸣后10s增伤10%，四次共鸣共40s
				"恒之罪": {DamageBoost: 25},
			},
		},
		"III": {
			Buffs: map[string]models.SkillBuff{
				"所有": {
					DamageBoost: 8 +
						0.8*12 + // 咒文守护下12%增伤，持续12s，冷却15s
						0.6*10}, // 共鸣后10s增伤10%，四次共鸣共40s
				"恒之罪": {DamageBoost: 25},
			},
		},
		"II": {
			Buffs: map[string]models.SkillBuff{
				"所有":   {DamageBoost: 8 + 0.8*12}, // 咒文守护下12%增伤，持续12s，冷却15s
				"恒之罪": {DamageBoost: 25},
			},
		},
		"I": {
			Buffs: map[string]models.SkillBuff{
				"所有":   {DamageBoost: 8},
				"恒之罪": {DamageBoost: 25},
			},
		},
	}
}
