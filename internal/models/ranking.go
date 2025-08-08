package models

type RankingItem struct {
	OpenID       string `json:"openid"`
	Contribution int32  `json:"contribution"`
	Rank         int    `json:"rank"`
}
