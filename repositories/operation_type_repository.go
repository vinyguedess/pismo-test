package repositories

import (
	"context"

	"gorm.io/gorm"

	"pismo/models"
)

type OperationTypeRepository interface {
	FindByID(ctx context.Context, operationTypeId int) (*models.OperationType, error)
}

func NewOperationTypeRepository(db *gorm.DB) OperationTypeRepository {
	return &operationTypeRepository{db: db}
}

type operationTypeRepository struct {
	db *gorm.DB
}

func (r *operationTypeRepository) FindByID(
	ctx context.Context, operationTypeId int,
) (*models.OperationType, error) {
	var operationType models.OperationType
	err := r.db.WithContext(ctx).
		First(&operationType, operationTypeId).
		Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, err
	}

	return &operationType, nil
}
