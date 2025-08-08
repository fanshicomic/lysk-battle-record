package models

import (
	"fmt"
	"unicode/utf8"

	"lysk-battle-record/internal/pkg"
)

type User struct {
	ID        string `json:"id"`
	Nickname  string `json:"nickname"`
	RowNumber int    `json:"row_number"`
}

func (u User) ValidateNickname() error {
	if utf8.RuneCountInString(u.Nickname) > 10 {
		return fmt.Errorf("昵称最长10个字: %s", u.Nickname)
	}

	detector, err := pkg.NewDetector()
	if err != nil {
		return err
	}

	if detector.ContainsSensitiveWords(u.Nickname) {
		return fmt.Errorf("昵称中包含敏感词")
	}

	return nil
}
