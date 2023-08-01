package repositories

import (
	"context"
	"pismo/models"

	"gorm.io/gorm"
)

type AccountRepository interface {
	Create(ctx context.Context, data models.Account) (*models.Account, error)
	FindByID(ctx context.Context, accountId int) (*models.Account, error)
	FindByDocumentNumber(ctx context.Context, documentNumber string) (*models.Account, error)
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

func (r *accountRepository) FindByID(ctx context.Context, accountId int) (*models.Account, error) {
	var account models.Account
	err := r.db.WithContext(ctx).
		Where("id = ?", accountId).
		First(&account).
		Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, err
	}

	return &account, nil
}

func (r *accountRepository) FindByDocumentNumber(
	ctx context.Context, documentNumber string,
) (*models.Account, error) {
	var account models.Account
	err := r.db.WithContext(ctx).
		Where("document_number = ?", documentNumber).
		First(&account).
		Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, err
	}

	return &account, nil
}
