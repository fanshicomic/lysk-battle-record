package companions

import "lysk-battle-record/internal/models"

type Companion interface {
	GetName() string

	GetCompanionFlow(stats models.Stats) models.CompanionFlow

	GetActiveSkill(stats models.Stats) models.Skill
	GetBasicAttack(stats models.Stats) models.Skill
	GetOathSkill(stats models.Stats) models.Skill
	GetResonanceSkill(stats models.Stats) models.Skill
	GetSupportSkill(stats models.Stats) models.Skill
}
