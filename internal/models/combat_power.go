package models

const (
	NoData = "无数据"
)

type CombatPower struct {
	Score         string
	BufferedScore string
}

type Stats struct {
	Attack       string
	HP           string
	Defense      string
	Matching     string
	MatchingBuff string
	CritRate     string
	CritDmg      string
	EnergyRegen  string
	WeakenBoost  string
	OathBoost    string
	OathRegen    string
	TotalLevel   string
	Partner      string
	SetCard      string
	Stage        string
	Weapon       string
	Buff         string
}

type Skill struct {
	Base        int
	HpRate      float64
	AttackRate  float64
	DefenceRate float64
	StageBuff   map[string]float64
	Count       int
}

type PartnerSkillSet struct {
	NormalAttack   Skill
	HeavyAttack    Skill
	ActiveSkill    Skill
	ResonanceSkill Skill
	OathSkill      Skill
	SupportSkill   Skill
	PassiveSkill   Skill
}

type PartnerPeriod struct {
	SkillSet PartnerSkillSet
	Count    int
	IsWeaken bool
}

type PartnerFlow struct {
	Periods []PartnerPeriod
}
