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
	Companion    string
	SetCard      string
	Stage        string
	Weapon       string
	Buff         float64
}

type Skill struct {
	Name                  string
	Base                  float64
	HpRate                float64
	AttackRate            float64
	DefenseRate           float64
	CritRate              float64
	CritDmg               float64
	WeakenBoost           float64
	DamageBoost           float64
	OathBoost             float64
	EnemyDefenceReduction float64
	Count                 int
	CanBeCrit             bool
	NoWeakenPeriod        bool
}

type CompanionSkillSet struct {
	Skills []Skill
}

type CompanionPeriod struct {
	SkillSet   CompanionSkillSet
	WeakenRate float64
	Boost      float64
}

type CompanionFlow struct {
	Periods []CompanionPeriod
}

type SkillBuff struct {
	CritRate         float64
	CritDmg          float64
	WeakenBoost      float64
	DamageBoost      float64
	OathBoost        float64
	DefenceReduction float64
	CountBonus       float64
}

type StageBuff struct {
	Buffs map[string]SkillBuff
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
