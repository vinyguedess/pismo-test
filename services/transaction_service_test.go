package services

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	"pismo/common"
	"pismo/entities"
	mock_repositories "pismo/mocks/repositories"
	"pismo/models"
)

type transactionServiceTestSuite struct {
	suite.Suite
	ctrl                         *gomock.Controller
	accountRepositoryMock        *mock_repositories.MockAccountRepository
	operationTypeRepoositoryMock *mock_repositories.MockOperationTypeRepository
	transactionRepositoryMock    *mock_repositories.MockTransactionRepository
	service                      TransactionService
}

func TestTransactionServiceTestSuite(t *testing.T) {
	suite.Run(t, new(transactionServiceTestSuite))
}

func (s *transactionServiceTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.accountRepositoryMock = mock_repositories.NewMockAccountRepository(s.ctrl)
	s.operationTypeRepoositoryMock = mock_repositories.NewMockOperationTypeRepository(s.ctrl)
	s.transactionRepositoryMock = mock_repositories.NewMockTransactionRepository(s.ctrl)
	s.service = NewTransactionService(
		s.accountRepositoryMock,
		s.operationTypeRepoositoryMock,
		s.transactionRepositoryMock,
	)
}

func (s *transactionServiceTestSuite) TestCreate() {
	account := models.Account{
		BaseModel: models.BaseModel{ID: 1},
	}

	operationType := models.OperationType{
		BaseModel:   models.BaseModel{ID: common.CashPaymentID},
		Description: common.CashPayment,
	}

	transaction := models.Transaction{
		AccountID:       account.ID,
		OperationTypeID: operationType.ID,
		Amount:          100,
	}

	tests := []struct {
		description                   string
		responseAccountFindById       *models.Account
		errAccountFindById            error
		responseOperationTypeFindById *models.OperationType
		errOperationTypeFindById      error
		responseTransactionCreate     *models.Transaction
		errTransactionCreate          error
	}{
		{
			description:                   "Success",
			responseAccountFindById:       &account,
			responseOperationTypeFindById: &operationType,
			responseTransactionCreate:     &transaction,
		},
		{
			description:        "Error on finding account by ID",
			errAccountFindById: errors.New("Error fetching account"),
		},
		{
			description: "Account not found",
		},
		{
			description:              "Error on finding operation type by ID",
			responseAccountFindById:  &account,
			errOperationTypeFindById: errors.New("Error fetching operation type"),
		},
		{
			description:             "Operation type not found",
			responseAccountFindById: &account,
		},
	}

	for _, test := range tests {
		s.Run(test.description, func() {
			ctx := context.Background()

			s.accountRepositoryMock.EXPECT().FindByID(ctx, account.ID).Return(
				test.responseAccountFindById, test.errAccountFindById,
			)

			if test.responseAccountFindById != nil && test.errAccountFindById == nil {
				s.operationTypeRepoositoryMock.EXPECT().FindByID(ctx, operationType.ID).Return(
					test.responseOperationTypeFindById, test.errOperationTypeFindById,
				)
			}

			if test.responseAccountFindById != nil &&
				test.errAccountFindById == nil &&
				test.errOperationTypeFindById == nil &&
				test.responseOperationTypeFindById != nil {
				s.transactionRepositoryMock.EXPECT().Create(ctx, transaction).Return(
					test.responseTransactionCreate, test.errTransactionCreate,
				)
			}

			newTransaction, err := s.service.Create(ctx, transaction)
			if test.errAccountFindById != nil {
				s.Error(err)
				s.ErrorContains(err, test.errAccountFindById.Error())
			} else if test.responseAccountFindById == nil {
				s.Error(err)
				s.IsType(&entities.ItemNotFoundError{}, err)
				s.ErrorContains(err, "Account not found")
			} else if test.errOperationTypeFindById != nil {
				s.Error(err)
				s.ErrorContains(err, test.errOperationTypeFindById.Error())
			} else if test.responseOperationTypeFindById == nil {
				s.Error(err)
				s.IsType(&entities.ItemNotFoundError{}, err)
				s.ErrorContains(err, "OperationType not found")
			} else if test.errTransactionCreate != nil {
				s.Error(err)
				s.ErrorContains(err, test.errTransactionCreate.Error())
			} else {
				s.NoError(err)
				s.Equal(test.responseTransactionCreate, newTransaction)
			}
		})
	}
}
