package models

import (
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	ID        int       `json:"id" gorm:"primarykey;type:integer;autoIncrement"`
	CreatedAt time.Time `json:"created_at" gorm:"type:datetime"`
}

func (b *BaseModel) BeforeCreate(tx *gorm.DB) error {
	b.CreatedAt = time.Now().UTC()
	return nil
}
