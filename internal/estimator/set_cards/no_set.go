package set_cards

import "lysk-battle-record/internal/models"

type NoSet struct{}

func (c NoSet) GetName() string {
	return "无套装"
}

func (c NoSet) GetSetCardBuff() map[string]models.StageBuff {
	return map[string]models.StageBuff{
		"IV": {
			Buffs: map[string]models.SkillBuff{},
		},
		"III": {
			Buffs: map[string]models.SkillBuff{},
		},
		"II": {
			Buffs: map[string]models.SkillBuff{},
		},
		"I": {
			Buffs: map[string]models.SkillBuff{},
		},
	}
}
