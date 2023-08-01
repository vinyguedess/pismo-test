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

	for _, operationType := range operationTypes {
		var operationTypeTotal int64
		db.Model(&models.OperationType{}).Where("description = ?").Count(&operationTypeTotal)
		if operationTypeTotal > 0 {
			continue
		}

		db.Create(&operationType)
	}
}
