package repositories

import (
	"context"

	"gorm.io/gorm"

	"pismo/models"
)

type TransactionRepository interface {
	Create(ctx context.Context, data models.Transaction) (*models.Transaction, error)
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db}
}

type transactionRepository struct {
	db *gorm.DB
}

func (r *transactionRepository) Create(
	ctx context.Context, data models.Transaction,
) (*models.Transaction, error) {
	if err := r.db.WithContext(ctx).Create(&data).Error; err != nil {
		return nil, err
	}

	return &data, nil
}
