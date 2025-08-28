package set_cards

import "lysk-battle-record/internal/models"

type FourStar struct{}

func (c FourStar) GetName() string {
	return "四星"
}

func (c FourStar) GetSetCardBuff() map[string]models.StageBuff {
	return map[string]models.StageBuff{
		"IV": {
			Buffs: map[string]models.SkillBuff{
				"所有": {DamageBoost: 10},
			},
		},
		"III": {
			Buffs: map[string]models.SkillBuff{
				"所有": {DamageBoost: 5},
			},
		},
		"II": {
			Buffs: map[string]models.SkillBuff{
				"所有": {DamageBoost: 5},
			},
		},
		"I": {
			Buffs: map[string]models.SkillBuff{
				"所有": {DamageBoost: 5},
			},
		},
	}
}
