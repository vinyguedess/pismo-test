package repositories

import (
	"context"
	"pismo/models"

	"gorm.io/gorm"
)

type AccountRepository interface {
	Create(ctx context.Context, data models.Account) (*models.Account, error)
}

type accountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) AccountRepository {
	return &accountRepository{db: db}
}

func (r *accountRepository) Create(ctx context.Context, data models.Account) (*models.Account, error) {
	err := r.db.WithContext(ctx).Create(&data).Error
	if err != nil {
		return nil, err
	}

	return &data, nil
}
