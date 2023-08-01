package repositories

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type operationTypeRepositoryTestSuite struct {
	suite.Suite
	ctx        context.Context
	dbconn     *gorm.DB
	dbmock     sqlmock.Sqlmock
	repository OperationTypeRepository
}

func TestOperationTypeRepositoryTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(operationTypeRepositoryTestSuite))
}

func (s *operationTypeRepositoryTestSuite) SetupTest() {
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

	s.repository = NewOperationTypeRepository(s.dbconn)
}

func (s *operationTypeRepositoryTestSuite) TestFindByID() {
	tests := []struct {
		description    string
		noResultsFound bool
		errInQuery     error
	}{
		{
			description: "Success",
		},
		{
			description:    "No results found",
			noResultsFound: true,
		},
		{
			description: "Error in query",
			errInQuery:  errors.New("error in query"),
		},
	}

	for _, test := range tests {
		s.Run(test.description, func() {
			expectedQuery := s.dbmock.ExpectQuery(
				regexp.QuoteMeta("SELECT * FROM `operation_types` WHERE `operation_types`.`id` = ? ORDER BY `operation_types`.`id` LIMIT 1"),
			).WithArgs(1)
			if test.errInQuery != nil {
				expectedQuery.WillReturnError(test.errInQuery)
			} else if test.noResultsFound {
				expectedQuery.WillReturnRows(sqlmock.NewRows([]string{"id", "description"}))
			} else {
				expectedQuery.WillReturnRows(sqlmock.NewRows([]string{"id", "description"}).AddRow(1, "description"))
			}

			operationType, err := s.repository.FindByID(s.ctx, 1)
			if test.errInQuery != nil {
				s.Error(err)
				s.Equal(test.errInQuery, err)
			} else if test.noResultsFound {
				s.NoError(err)
				s.Nil(operationType)
			} else {
				s.NoError(err)
				s.NotNil(operationType)
				s.Equal(1, operationType.ID)
				s.Equal("description", operationType.Description)
			}
		})
	}
}
