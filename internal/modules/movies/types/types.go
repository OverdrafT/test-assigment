package types

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Movie struct {
	ID        uuid.UUID `gorm:"primary_key"`
	MovieYear int       `gorm:"type:int;not null"` //`json:"movieyear"`
	MovieName string    `gorm:"type:varchar(255);not null"`
}

func (m *Movie) BeforeCreate(tx *gorm.DB) (err error) {
	m.ID = uuid.New()
	return
}
