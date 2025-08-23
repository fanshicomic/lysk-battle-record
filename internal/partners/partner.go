package partners

import "lysk-battle-record/internal/models"

type Partner interface {
	GetName() string

	GetPartnerFlow(stats models.Stats) models.PartnerFlow

	GetActiveSkill(stats models.Stats) models.Skill
	GetHeavyAttack(stats models.Stats) models.Skill
	GetOathSkill(stats models.Stats) models.Skill
	GetResonanceSkill() models.Skill
	GetSupportSkill() models.Skill
}
