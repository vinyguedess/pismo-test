package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uuid.UUID `json:"id" gorm:"primarykey;type:varchar(36)"`
	CreatedAt time.Time `json:"created_at" gorm:"type:datetime"`
}

func (b *BaseModel) BeforeCreate(tx *gorm.DB) error {
	b.ID = uuid.New()
	b.CreatedAt = time.Now().UTC()
	return nil
}
