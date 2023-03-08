package types

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Author struct {
	AuthorID  uuid.UUID `json:"ID" gorm:"primary_key"`
	FirstName string    `json:"FirstName" gorm:"type:varchar(255);not null"`
	LastName  string    `json:"LastName" gorm:"type:varchar(255);not null"`
	Movie     []Movie   `gorm:"ForeignKey:AuthorID"`
}

type Movie struct {
	ID        uuid.UUID `gorm:"primary_key"`
	MovieYear int       `json:"MovieYear" gorm:"type:int;not null"`
	MovieName string    `json:"MovieName" gorm:"type:varchar(255);not null"`
	AuthorID  string    `json:"AuthorID"`
}

func (m *Movie) BeforeCreate(tx *gorm.DB) (err error) {
	m.ID = uuid.New()
	return
}

func (a *Author) BeforeCreate(tx *gorm.DB) (err error) {
	a.AuthorID = uuid.New()
	return
}
