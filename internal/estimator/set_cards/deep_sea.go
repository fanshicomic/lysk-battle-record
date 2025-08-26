package set_cards

import "lysk-battle-record/internal/models"

type DeepSea struct{}

func (c DeepSea) GetName() string {
	return "深海"
}

func (c DeepSea) GetSetCardBuff() map[string]models.StageBuff {
	return map[string]models.StageBuff{
		"IV": {
			Buffs: map[string]models.SkillBuff{
				"所有":     {DamageBoost: 16 + 20*0.7},
				"强力斩击": {CritDmg: 30, DamageBoost: 150},
			},
		},
		"III": {
			Buffs: map[string]models.SkillBuff{
				"所有":     {DamageBoost: 8 + 20*0.7}, // 灼烧20%增伤，大概持续70%时间（技能回能增加
				"强力斩击": {CritDmg: 30},
			},
		},
		"II": {
			Buffs: map[string]models.SkillBuff{
				"所有":     {DamageBoost: 8 + 20*0.5}, // 灼烧20%增伤，大概持续一半时间
				"强力斩击": {CritDmg: 30},
			},
		},
		"I": {
			Buffs: map[string]models.SkillBuff{
				"所有":     {DamageBoost: 8},
				"强力斩击": {CritDmg: 30},
			},
		},
	}
}
