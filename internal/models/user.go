package models

type User struct {
	ID        string `json:"id"`
	Nickname  string `json:"nickname"`
	RowNumber int    `json:"row_number"`
}
