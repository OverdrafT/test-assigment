package types

type Movie struct {
	ID        uint   `json:"id" gorm:"primary_key"`
	MovieYear int    `json:"movieyear"`
	MovieName string `json:"moviename"`
}
