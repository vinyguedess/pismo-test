package models

type Transaction struct {
	BaseModel
	AccountID       int     `json:"account_id"`
	OperationTypeID int     `json:"operation_type_id"`
	Amount          float64 `json:"amount" gorm:"type:decimal(8,2);not null"`

	Account       Account       `json:"account" gorm:"foreignKey:AccountID"`
	OperationType OperationType `json:"operation_type" gorm:"foreignKey:OperationTypeID"`
}
