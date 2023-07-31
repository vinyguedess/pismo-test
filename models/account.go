package models

type Account struct {
	*BaseModel
	DocumentNumber string `json:"document_number" gorm:"unique;not null"`
}
