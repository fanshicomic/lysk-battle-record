package set_cards

import "lysk-battle-record/internal/models"

type Midnight struct{}

func (c Midnight) GetName() string {
	return "末夜"
}

func (c Midnight) GetSetCardBuff() map[string]models.StageBuff {
	return map[string]models.StageBuff{
		"IV": {
			Buffs: map[string]models.SkillBuff{
				"所有": {DamageBoost: 16},
				"月光": {DamageBoost: 25, CountBonus: 1.34, CritDmg: 30},
			},
		},
		"III": {
			Buffs: map[string]models.SkillBuff{
				"所有": {DamageBoost: 8},
				"月光": {DamageBoost: 25, CountBonus: 1.34},
			},
		},
		"II": {
			Buffs: map[string]models.SkillBuff{
				"所有": {DamageBoost: 8},
				"月光": {DamageBoost: 25, CountBonus: 1.34},
			},
		},
		"I": {
			Buffs: map[string]models.SkillBuff{
				"所有": {DamageBoost: 8},
			},
		},
	}
}
