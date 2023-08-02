package providers

import (
	"gorm.io/gorm"

	"pismo/common"
	"pismo/models"
)

func SeedDatabase(db *gorm.DB) {
	seedOperationTypes(db)
}

func seedOperationTypes(db *gorm.DB) {
	operationTypes := []string{
		common.CashPayment,
		common.InstallmentPurchase,
		common.Withdraw,
		common.Payment,
	}

	for _, operationTypeDesc := range operationTypes {
		var operationTypeTotal int64
		db.Model(&models.OperationType{}).
			Where("description = ?", operationTypeDesc).
			Count(&operationTypeTotal)
		if operationTypeTotal > 0 {
			continue
		}

		operationType := &models.OperationType{
			Description: operationTypeDesc,
		}
		db.Create(operationType)
	}
}
