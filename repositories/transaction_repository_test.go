package repositories

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"pismo/common"
	"pismo/models"
)

type transactionRepositoryTestSuite struct {
	suite.Suite
	ctx        context.Context
	dbconn     *gorm.DB
	dbmock     sqlmock.Sqlmock
	repository TransactionRepository
}

func TestTransactionRepositoryTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(transactionRepositoryTestSuite))
}

func (s *transactionRepositoryTestSuite) SetupTest() {
	s.ctx = context.Background()

	conn, dbmock, _ := sqlmock.New()
	dialector := mysql.Dialector{
		Config: &mysql.Config{
			DSN:                       "sqlmock_db_0",
			Conn:                      conn,
			SkipInitializeWithVersion: true,
		},
	}

	dbconn, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		s.FailNow(err.Error())
	}

	s.dbmock = dbmock
	s.dbconn = dbconn

	s.repository = NewTransactionRepository(s.dbconn)
}

func (s *transactionRepositoryTestSuite) TestCreate() {
	tests := []struct {
		description  string
		errorInQuery error
	}{
		{
			description: "Success",
		},
		{
			description:  "Error in query",
			errorInQuery: errors.New("error in query"),
		},
	}

	for _, test := range tests {
		s.Run(test.description, func() {
			s.dbmock.ExpectBegin()
			expectedInsertQuery := s.dbmock.ExpectExec("INSERT INTO `transactions`").WithArgs(
				sqlmock.AnyArg(), 1, common.CashPaymentID, 10.00,
			)
			if test.errorInQuery != nil {
				expectedInsertQuery.WillReturnError(test.errorInQuery)
				s.dbmock.ExpectRollback()
			} else {
				expectedInsertQuery.WillReturnResult(sqlmock.NewResult(1, 1))
				s.dbmock.ExpectCommit()
			}

			transaction, err := s.repository.Create(
				s.ctx,
				models.Transaction{
					AccountID:       1,
					OperationTypeID: common.CashPaymentID,
					Amount:          10.00,
				},
			)
			if err != nil {
				s.Nil(transaction)
				s.Error(err)
				s.ErrorContains(err, test.errorInQuery.Error())
			} else {
				s.NoError(err)
				s.Equal(transaction.AccountID, 1)
				s.Equal(transaction.OperationTypeID, common.CashPaymentID)
				s.Equal(transaction.Amount, 10.00)
			}
		})
	}
}
