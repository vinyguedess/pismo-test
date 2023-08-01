package services

import (
	"context"

	"pismo/entities"
	"pismo/models"
	"pismo/repositories"
)

type AccountService interface {
	Create(ctx context.Context, data models.Account) (*models.Account, error)
	FindByID(ctx context.Context, accountId int) (*models.Account, error)
}

type accountService struct {
	accountRepository repositories.AccountRepository
}

func NewAccountService(accountRepository repositories.AccountRepository) AccountService {
	return &accountService{
		accountRepository: accountRepository,
	}
}

func (s *accountService) Create(ctx context.Context, data models.Account) (*models.Account, error) {
	account, err := s.accountRepository.FindByDocumentNumber(ctx, data.DocumentNumber)
	if err != nil {
		return nil, err
	} else if account != nil {
		return nil, entities.NewAccountAlreadyExistsError(data.DocumentNumber)
	}

	return s.accountRepository.Create(ctx, data)
}

func (s *accountService) FindByID(ctx context.Context, accountId int) (*models.Account, error) {
	account, err := s.accountRepository.FindByID(ctx, accountId)
	if err != nil {
		return nil, err
	} else if account == nil {
		return nil, entities.NewItemNotFoundError("Account", accountId)
	}

	return account, nil
}
