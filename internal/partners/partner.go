package partners

import "lysk-battle-record/internal/models"

type Partner interface {
	GetName() string

	GetPartnerFlow(stats models.Stats) models.PartnerFlow

	GetActiveSkill(stats models.Stats) models.Skill
	GetBasicAttack(stats models.Stats) models.Skill
	GetOathSkill(stats models.Stats) models.Skill
	GetResonanceSkill(stats models.Stats) models.Skill
	GetSupportSkill(stats models.Stats) models.Skill
}
