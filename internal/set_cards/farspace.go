package set_cards

import "lysk-battle-record/internal/models"

type Farspace struct{}

func (c Farspace) GetName() string {
	return "无套装"
}

func (c Farspace) GetSetCardBuff() map[string]models.StageBuff {
	return map[string]models.StageBuff{
		"IV": {
			Buffs: map[string]models.SkillBuff{
				"所有":     {DamageBoost: 30},
				"纵深打击": {DamageBoost: 40},
			},
		},
		"III": {
			Buffs: map[string]models.SkillBuff{
				"所有":     {DamageBoost: 17},
				"纵深打击": {DamageBoost: 40},
			},
		},
		"II": {
			Buffs: map[string]models.SkillBuff{
				"所有":     {DamageBoost: 8},
				"纵深打击": {DamageBoost: 40},
			},
		},
		"I": {
			Buffs: map[string]models.SkillBuff{
				"所有": {DamageBoost: 8},
			},
		},
	}
}
