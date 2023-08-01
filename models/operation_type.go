package models

type OperationType struct {
	BaseModel
	Description string `json:"description" gorm:"unique;not null"`
}
