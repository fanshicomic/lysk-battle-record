package set_cards

import "lysk-battle-record/internal/models"

type Temple struct{}

func (c Temple) GetName() string {
	return "神殿"
}

func (c Temple) GetSetCardBuff() map[string]models.StageBuff {
	return map[string]models.StageBuff{
		"IV": {
			Buffs: map[string]models.SkillBuff{
				"所有": {DamageBoost: 16, CritDmg: 10},
				"海灵": {DamageBoost: 12.5, CountBonus: 1.7},
			},
		},
		"III": {
			Buffs: map[string]models.SkillBuff{
				"所有": {DamageBoost: 8, CritDmg: 9},
				"海灵": {DamageBoost: 12.5, CountBonus: 1.7}, // 下雨期间额外25%增伤，III阶大约3次下雨
			},
		},
		"II": {
			Buffs: map[string]models.SkillBuff{
				"所有": {DamageBoost: 8, CritDmg: 8},
				"海灵": {DamageBoost: 8.3, CountBonus: 1.7}, // 攻击次数由7变12，下雨期间额外25%增伤，II阶大约2次下雨
			},
		},
		"I": {
			Buffs: map[string]models.SkillBuff{
				"所有": {DamageBoost: 8, CritDmg: 7},
			},
		},
	}
}
