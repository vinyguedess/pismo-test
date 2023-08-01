package repositories

import (
	"context"
	"errors"
	"pismo/models"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type accountRepositoryTestSuite struct {
	suite.Suite
	ctx        context.Context
	dbconn     *gorm.DB
	dbmock     sqlmock.Sqlmock
	repository AccountRepository
}

func TestAccountRepository(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(accountRepositoryTestSuite))
}

func (s *accountRepositoryTestSuite) SetupTest() {
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

	s.repository = NewAccountRepository(s.dbconn)
}

func (s *accountRepositoryTestSuite) TestCreate() {
	tests := []struct {
		description      string
		insertQueryError error
	}{
		{
			description: "Success",
		},
		{
			description:      "Error in query",
			insertQueryError: errors.New("error in query"),
		},
	}

	for _, test := range tests {
		s.Run(test.description, func() {
			s.SetupTest()

			s.dbmock.ExpectBegin()

			insertQuery := s.dbmock.ExpectExec("INSERT INTO `accounts`").WithArgs(
				sqlmock.AnyArg(),
				"123456789",
			)
			if test.insertQueryError != nil {
				insertQuery.WillReturnError(test.insertQueryError)
				s.dbmock.ExpectRollback()
			} else {
				insertQuery.WillReturnResult(sqlmock.NewResult(1, 1))
				s.dbmock.ExpectCommit()
			}

			account, err := s.repository.Create(s.ctx, models.Account{DocumentNumber: "123456789"})
			if test.insertQueryError != nil {
				s.Nil(account)
				s.Error(err)
				s.ErrorContains(err, "error in query")
			} else {
				s.NoError(err)
				s.Equal("123456789", account.DocumentNumber)
			}
			s.NoError(s.dbmock.ExpectationsWereMet())
		})
	}
}

func (s *accountRepositoryTestSuite) TestFindByID() {
	accountId := 1

	tests := []struct {
		description  string
		isNotFound   bool
		errorInQuery error
	}{
		{
			description: "Success",
		},
		{
			description: "Not found",
			isNotFound:  true,
		},
		{
			description:  "Error in query",
			errorInQuery: errors.New("error in query"),
		},
	}

	for _, test := range tests {
		s.Run(test.description, func() {
			s.SetupTest()

			expectedQuery := s.dbmock.ExpectQuery(
				regexp.QuoteMeta(
					"SELECT * FROM `accounts` WHERE id = ? ORDER BY `accounts`.`id` LIMIT 1",
				),
			).WithArgs(accountId)
			if test.errorInQuery != nil {
				expectedQuery.WillReturnError(test.errorInQuery)
			} else if test.isNotFound {
				expectedQuery.WillReturnRows(
					sqlmock.NewRows([]string{"id", "document_number"}),
				)
			} else {
				expectedQuery.WillReturnRows(
					sqlmock.NewRows([]string{"id", "document_number"}).AddRow(accountId, "123456789"),
				)
			}

			account, err := s.repository.FindByID(s.ctx, accountId)
			if test.errorInQuery != nil {
				s.Error(err)
				s.ErrorContains(err, "error in query")
			} else if test.isNotFound {
				s.NoError(err)
				s.Nil(account)
			} else {
				s.NoError(err)
				s.Equal("123456789", account.DocumentNumber)
				s.NoError(s.dbmock.ExpectationsWereMet())
			}
		})
	}
}

func (s *accountRepositoryTestSuite) TestFindByDocumentNumber() {
	accountId := 1

	tests := []struct {
		description  string
		isNotFound   bool
		errorInQuery error
	}{
		{
			description: "Success",
		},
		{
			description: "Not found",
			isNotFound:  true,
		},
		{
			description:  "Error in query",
			errorInQuery: errors.New("error in query"),
		},
	}

	for _, test := range tests {
		s.Run(test.description, func() {
			s.SetupTest()

			expectedQuery := s.dbmock.ExpectQuery(
				regexp.QuoteMeta(
					"SELECT * FROM `accounts` WHERE document_number = ? ORDER BY `accounts`.`id` LIMIT 1",
				),
			).WithArgs("123456789")
			if test.errorInQuery != nil {
				expectedQuery.WillReturnError(test.errorInQuery)
			} else if test.isNotFound {
				expectedQuery.WillReturnRows(
					sqlmock.NewRows([]string{"id", "document_number"}),
				)
			} else {
				expectedQuery.WillReturnRows(
					sqlmock.NewRows([]string{"id", "document_number"}).AddRow(accountId, "123456789"),
				)
			}

			account, err := s.repository.FindByDocumentNumber(s.ctx, "123456789")
			if test.errorInQuery != nil {
				s.Error(err)
				s.ErrorContains(err, "error in query")
			} else if test.isNotFound {
				s.NoError(err)
				s.Nil(account)
			} else {
				s.NoError(err)
				s.Equal("123456789", account.DocumentNumber)
				s.NoError(s.dbmock.ExpectationsWereMet())
			}
		})
	}
}
