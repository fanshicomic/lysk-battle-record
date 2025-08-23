package set_cards

import "lysk-battle-record/internal/models"

type SetCard interface {
	GetName() string

	GetSetCardBuff() map[string]models.StageBuff
}
