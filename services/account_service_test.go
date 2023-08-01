package services

import (
	"context"
	"errors"
	"pismo/entities"
	mock_repositories "pismo/mocks/repositories"
	"pismo/models"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type accountServiceTestSuite struct {
	suite.Suite
	ctrl                  *gomock.Controller
	accountRepositoryMock *mock_repositories.MockAccountRepository
	service               AccountService
}

func TestAccountServiceTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(accountServiceTestSuite))
}

func (s *accountServiceTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.accountRepositoryMock = mock_repositories.NewMockAccountRepository(s.ctrl)
	s.service = NewAccountService(s.accountRepositoryMock)
}

func (s *accountServiceTestSuite) TestCreate() {
	accountId := uuid.New()
	account := &models.Account{
		BaseModel:      models.BaseModel{ID: accountId},
		DocumentNumber: "123456789",
	}

	tests := []struct {
		description                  string
		data                         models.Account
		responseFindByDocumentNumber *models.Account
		errFindByDocumentNumber      error
		responseCreate               *models.Account
		errCreate                    error
	}{
		{
			description:    "Success",
			data:           models.Account{DocumentNumber: "123456789"},
			responseCreate: account,
		},
		{
			description:             "Error finding by document number",
			data:                    models.Account{DocumentNumber: "123456789"},
			errFindByDocumentNumber: errors.New("error finding by ID"),
		},
		{
			description:                  "Account already exists",
			data:                         models.Account{DocumentNumber: "123456789"},
			responseFindByDocumentNumber: account,
		},
		{
			description: "Error creating",
			data:        models.Account{DocumentNumber: "123456789"},
			errCreate:   errors.New("error creating"),
		},
	}

	for _, test := range tests {
		s.Run(test.description, func() {
			ctx := context.Background()

			s.accountRepositoryMock.EXPECT().
				FindByDocumentNumber(ctx, test.data.DocumentNumber).
				Return(test.responseFindByDocumentNumber, test.errFindByDocumentNumber)

			if test.errFindByDocumentNumber == nil && test.responseFindByDocumentNumber == nil {
				s.accountRepositoryMock.EXPECT().
					Create(ctx, test.data).
					Return(test.responseCreate, test.errCreate)
			}

			response, err := s.service.Create(ctx, test.data)
			if test.errFindByDocumentNumber != nil {
				s.Error(err)
				s.ErrorContains(err, test.errFindByDocumentNumber.Error())
			} else if test.responseFindByDocumentNumber != nil {
				s.Error(err)
				s.ErrorContains(err, "Account already exists")
				s.IsType(&entities.AccountAlreadyExistsError{}, err)
			} else if test.errCreate != nil {
				s.Error(err)
				s.ErrorContains(err, test.errCreate.Error())
			} else {
				s.NoError(err)
				s.Equal(test.responseCreate, response)
			}
		})
	}
}
