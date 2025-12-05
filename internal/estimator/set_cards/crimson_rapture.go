package set_cards

import "lysk-battle-record/internal/models"

type CrimsonRapture struct{}

func (c CrimsonRapture) GetName() string {
	return "猩红"
}

func (c CrimsonRapture) GetSetCardBuff() map[string]models.StageBuff {
	return map[string]models.StageBuff{
		"IV": {
			Buffs: map[string]models.SkillBuff{
				"所有":          {DamageBoost: 16},
				"蔷薇地棘":      {DamageBoost: 10 + 10},
				"蔷薇地棘-血誓": {DamageBoost: 10 + 10, EnemyWeakenBoost: (8.0 * 2.0) / 3.0}, // 一共3次，只有两次能吃到buff
				"蔷薇之雨":      {DamageBoost: 10, EnemyWeakenBoost: 4},                      // 一共2次，只有一次能吃到buff
				"誓约":          {DamageBoost: 20, EnemyWeakenBoost: 8},
				"重击-血誓":     {DamageBoost: 10, EnemyWeakenBoost: 8},
				"血破之湮":      {DamageBoost: 10, EnemyWeakenBoost: 8},
				"堙界之棂":      {DamageBoost: 20, EnemyWeakenBoost: 8},
				"血蔷薇子弹":    {DamageBoost: 10, EnemyWeakenBoost: 8},
				"主动":          {CountBonus: 1.6},
			},
		},
		"III": {
			Buffs: map[string]models.SkillBuff{
				"所有":          {DamageBoost: 8},
				"蔷薇地棘":      {DamageBoost: 10},
				"蔷薇地棘-血誓": {DamageBoost: 10, EnemyWeakenBoost: (8.0 * 2.0) / 3.0}, // 一共3次，只有两次能吃到buff
				"蔷薇之雨":      {EnemyWeakenBoost: 4},                                  // 一共2次，只有一次能吃到buff
				"誓约":          {EnemyWeakenBoost: 8},
				"重击-血誓":     {EnemyWeakenBoost: 8},
				"血破之湮":      {EnemyWeakenBoost: 8},
				"堙界之棂":      {EnemyWeakenBoost: 8},
				"血蔷薇子弹":    {CountBonus: 0, EnemyWeakenBoost: 8},
				"主动":          {CountBonus: 1.6},
			},
		},
		"II": {
			Buffs: map[string]models.SkillBuff{
				"所有":          {DamageBoost: 8},
				"蔷薇地棘":      {DamageBoost: 10},
				"蔷薇地棘-血誓": {DamageBoost: 10},
				"蔷薇之雨":      {CountBonus: 2},
				"血蔷薇子弹":    {CountBonus: 0},
				"主动":          {CountBonus: 1.6},
			},
		},
		"I": {
			Buffs: map[string]models.SkillBuff{
				"所有":       {DamageBoost: 8},
				"血蔷薇子弹": {CountBonus: 0},
				"主动":       {CountBonus: 1.6},
			},
		},
	}
}
