package services

import (
	"context"

	"pismo/entities"
	"pismo/models"
	"pismo/repositories"
)

type TransactionService interface {
	Create(ctx context.Context, data models.Transaction) (*models.Transaction, error)
}

func NewTransactionService(
	accountRepository repositories.AccountRepository,
	operationTypeRepository repositories.OperationTypeRepository,
	transactionRepository repositories.TransactionRepository,
) TransactionService {
	return &transactionService{
		accountRepository:       accountRepository,
		operationTypeRepository: operationTypeRepository,
		transactionRepository:   transactionRepository,
	}
}

type transactionService struct {
	accountRepository       repositories.AccountRepository
	operationTypeRepository repositories.OperationTypeRepository
	transactionRepository   repositories.TransactionRepository
}

func (s *transactionService) Create(
	ctx context.Context, data models.Transaction,
) (*models.Transaction, error) {
	account, err := s.accountRepository.FindByID(ctx, data.AccountID)
	if err != nil {
		return nil, err
	} else if account == nil {
		return nil, entities.NewItemNotFoundError("Account", data.AccountID)
	}

	operationType, err := s.operationTypeRepository.FindByID(ctx, data.OperationTypeID)
	if err != nil {
		return nil, err
	} else if operationType == nil {
		return nil, entities.NewItemNotFoundError("OperationType", data.OperationTypeID)
	}

	return s.transactionRepository.Create(ctx, data)
}
