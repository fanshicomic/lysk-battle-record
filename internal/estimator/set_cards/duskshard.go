package set_cards

import "lysk-battle-record/internal/models"

type Duskshard struct{}

func (c Duskshard) GetName() string {
	return "沉冥"
}

func (c Duskshard) GetSetCardBuff() map[string]models.StageBuff {
	return map[string]models.StageBuff{
		"IV": {
			Buffs: map[string]models.SkillBuff{
				"所有":          {DamageBoost: 16},
				"誓约-尘露不忘": {CountBonus: 2},
				"并蒂莲":        {CountBonus: 1.5},
				"鬼爪":          {DamageBoost: 6},
				"三阶额外鬼爪":  {DamageBoost: 10},
			},
		},
		"III": {
			Buffs: map[string]models.SkillBuff{
				"所有":                 {DamageBoost: 8},
				"誓约-尘露不忘":        {DamageBoost: 80}, // 理论上来说是加载得更快，但因为不到2倍，所以改用伤害提升来计算了
				"并蒂莲":               {CountBonus: 1.5},
				"鬼爪":                 {DamageBoost: 6},
				"重击-阴阳":            {NotApplicable: true},
				"并蒂莲-阴阳":          {NotApplicable: true},
				"三阶阴阳重击额外伤害": {NotApplicable: true},
				"三阶额外鬼爪":         {NotApplicable: true},
			},
		},
		"II": {
			Buffs: map[string]models.SkillBuff{
				"所有":                 {DamageBoost: 8},
				"并蒂莲":               {CountBonus: 1.5},
				"重击-阴阳":            {NotApplicable: true},
				"并蒂莲-阴阳":          {NotApplicable: true},
				"三阶阴阳重击额外伤害": {NotApplicable: true},
				"三阶额外鬼爪":         {NotApplicable: true},
			},
		},
		"I": {
			Buffs: map[string]models.SkillBuff{
				"所有":                 {DamageBoost: 8},
				"重击-阴阳":            {NotApplicable: true},
				"并蒂莲-阴阳":          {NotApplicable: true},
				"三阶阴阳重击额外伤害": {NotApplicable: true},
				"三阶额外鬼爪":         {NotApplicable: true},
			},
		},
	}
}
