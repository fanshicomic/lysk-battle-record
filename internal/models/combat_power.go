package models

const (
	NoData = "无数据"
)

type CombatPower struct {
	Score          string
	BuffedScore    string
	WeakenScore    string
	NonWeakenScore string
}

type Stats struct {
	Attack       int
	HP           int
	Defense      int
	Matching     string
	MatchingBuff float64
	CritRate     float64
	CritDmg      float64
	EnergyRegen  float64
	WeakenBoost  float64
	OathBoost    float64
	OathRegen    float64
	TotalLevel   int
	Partner      string
	SetCard      string
	Stage        string
	Weapon       string
	Buff         float64
}

type Skill struct {
	Name        string
	Base        float64
	HpRate      float64
	AttackRate  float64
	DefenseRate float64
	CritRate    float64
	CritDmg     float64
	WeakenBoost float64
	DamageBoost float64
	Count       int
}

type PartnerSkillSet struct {
	Skills []Skill
}

type PartnerPeriod struct {
	SkillSet PartnerSkillSet
}

type PartnerFlow struct {
	Periods    []PartnerPeriod
	WeakenRate float64
	Boost      float64
}

type SkillBuff struct {
	CritRate         float64
	CritDmg          float64
	WeakenBoost      float64
	DamageBoost      float64
	DefenceReduction float64
	CountBonus       float64
}

type StageBuff struct {
	Buffs map[string]SkillBuff
}

func GetLightSeekerFlow(stats Stats) PartnerFlow {
	activeSkill := stats.GetActiveSkill()
	heavyAttack := stats.GetHeavyAttack()
	resonanceSkill := stats.GetResonanceSkill()
	oathSkill := stats.GetOathSkill()
	supportSkill := stats.GetSupportSkill()

	// others
	passiveSkill := Skill{
		Name:       "溯光共鸣",
		Base:       150,
		AttackRate: 200,
	}

	weakenRate := 0.17 // 逐光虚弱期10s
	if stats.Matching == "顺" {
		weakenRate = 0.34
	}

	flow := PartnerFlow{
		Periods: []PartnerPeriod{
			{
				SkillSet: PartnerSkillSet{
					Skills: []Skill{
						activeSkill,
						heavyAttack,
						resonanceSkill,
						oathSkill,
						supportSkill,
						passiveSkill,
					},
				},
			},
		},
		Boost:      21, // 逐光骑士溯光力场10%攻击提升+溯光剑破盾后的20%增伤
		WeakenRate: weakenRate,
	}

	setCardBuff := stats.GetSetCardBuff()
	applySetCardBuff(&flow, setCardBuff)

	return flow
}

func applySetCardBuff(flow *PartnerFlow, setCardBuff StageBuff) {
	for periodIdx := range flow.Periods {
		period := &flow.Periods[periodIdx]
		for i, skill := range period.SkillSet.Skills {
			if allBuff, exists := setCardBuff.Buffs["所有"]; exists {
				skill.CritRate += allBuff.CritRate
				skill.CritDmg += allBuff.CritDmg
				skill.WeakenBoost += allBuff.WeakenBoost
				skill.DamageBoost += allBuff.DamageBoost
				if allBuff.CountBonus > 1 {
					skill.Count = int(float64(skill.Count) * allBuff.CountBonus)
				}
			}

			if skillBuff, exists := setCardBuff.Buffs[skill.Name]; exists {
				skill.CritRate += skillBuff.CritRate
				skill.CritDmg += skillBuff.CritDmg
				skill.WeakenBoost += skillBuff.WeakenBoost
				skill.DamageBoost += skillBuff.DamageBoost
				if skillBuff.CountBonus > 1 {
					skill.Count = int(float64(skill.Count) * skillBuff.CountBonus)
				}
			}

			period.SkillSet.Skills[i] = skill
		}
	}
}

func (s Stats) GetSetCardBuff() StageBuff {
	switch s.SetCard {
	case "逐光":
		switch s.Stage {
		case "IV":
			return StageBuff{
				Buffs: map[string]SkillBuff{
					"所有": {DamageBoost: 8},
					"主动": {DamageBoost: 25, CountBonus: 1.4},
				},
			}
		default:
			return StageBuff{
				Buffs: map[string]SkillBuff{
					"所有": {DamageBoost: 8},
					"主动": {DamageBoost: 25},
				},
			}
		}
	}
	return StageBuff{Buffs: make(map[string]SkillBuff)}
}

func (s Stats) GetActiveSkill() Skill {
	activeSkill := Skill{
		Name: "主动",
	}
	energy := s.GetEnergy()

	if s.Weapon == "专武" {
		switch s.Partner {
		case "逐光骑士":
			activeSkill.Base = 341
			activeSkill.AttackRate = 455
			activeSkill.Count = energy - 8
			activeSkill.Count *= 2 // 主动命中减冷却加能量
			break
		default:

		}
	} else if s.Weapon == "重剑" {
		activeSkill.Base = 621
		activeSkill.AttackRate = 829
		activeSkill.Count = int(float64(energy-8) / 1.2)
		activeSkill.DamageBoost = 150
	}

	return activeSkill
}

func (s Stats) GetHeavyAttack() Skill {
	heavyAttack := Skill{
		Name:  "重击",
		Count: 30, // hard code it as 30
	}

	if s.Weapon == "专武" {
		switch s.Partner {
		case "逐光骑士":
			heavyAttack.Base = 118
			heavyAttack.AttackRate = 157
			break
		default:

		}
	} else if s.Weapon == "重剑" {
		heavyAttack.Base = 337
		heavyAttack.AttackRate = 449
		heavyAttack.DamageBoost = 112.9
		heavyAttack.Count = 10
	}

	return heavyAttack
}

func (s Stats) GetResonanceSkill() Skill {
	resonanceSkill := Skill{
		Name:  "共鸣",
		Count: 4,
	}

	switch s.Partner {
	case "逐光骑士":
		resonanceSkill.Base = 641
		resonanceSkill.AttackRate = 854
	default:

	}

	return resonanceSkill
}

func (s Stats) GetOathSkill() Skill {
	oathSkill := Skill{
		Name:        "誓约",
		DamageBoost: s.OathBoost * 100,
	}

	if s.Stage != "无套装" && s.Stage != "I" {
		oathSkill.Count = 1
	}
	if s.OathRegen >= 17 {
		oathSkill.Count = 1
	}

	switch s.Partner {
	case "逐光骑士":
		oathSkill.Base = 1440
		oathSkill.AttackRate = 1920
	case "终极兵器X-02":
		// should change skill name and count
	default:

	}

	return oathSkill
}

func (s Stats) GetSupportSkill() Skill {
	supportSkill := Skill{
		Name: "协助",
	}

	switch s.Partner {
	case "逐光骑士":
		supportSkill.Base = 400
		supportSkill.AttackRate = 400
		supportSkill.Count = 6
	default:

	}

	return supportSkill
}

func (s Stats) GetEnergy() int {
	energy := 8
	if s.EnergyRegen >= 39.6 {
		energy = 12
	} else if s.EnergyRegen >= 30 {
		energy = 11
	} else if s.EnergyRegen >= 10.8 {
		energy = 10
	} else if s.EnergyRegen >= 6 {
		energy = 9
	}

	if s.Stage == "III" || s.Stage == "IV" {
		energy += 1
	}

	return energy
}
