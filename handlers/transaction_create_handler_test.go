package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	"pismo/entities"
	mock_services "pismo/mocks/services"
	"pismo/models"
)

func TestTransactionCreateHandlerTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(transactionCreateHandlerTestSuite))
}

type transactionCreateHandlerTestSuite struct {
	suite.Suite
	ctrl                   *gomock.Controller
	transactionServiceMock *mock_services.MockTransactionService
	handler                Handler
}

func (s *transactionCreateHandlerTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.transactionServiceMock = mock_services.NewMockTransactionService(s.ctrl)
	s.handler = NewTransactionCreateHandler(s.transactionServiceMock)
}

func (s *transactionCreateHandlerTestSuite) TestRoute() {
	s.Equal("/transactions", s.handler.Route())
}

func (s *transactionCreateHandlerTestSuite) TestMethod() {
	s.Equal([]string{http.MethodPost}, s.handler.Method())
}

func (s *transactionCreateHandlerTestSuite) TestServeHTTP() {
	transactionId := 1
	transaction := &models.Transaction{
		BaseModel: models.BaseModel{ID: transactionId},
	}

	tests := []struct {
		description        string
		payload            string
		invalidPayload     bool
		expectedPayload    models.Transaction
		responseCreate     *models.Transaction
		errorCreate        error
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			description: "Success",
			payload:     `{"account_id":1,"operation_type_id":1,"amount":100.0}`,
			expectedPayload: models.Transaction{
				AccountID: 1, OperationTypeID: 1, Amount: 100.0,
			},
			responseCreate:     transaction,
			expectedStatusCode: http.StatusCreated,
		},
		{
			description:        "Invalid JSON",
			payload:            `{"account_id":1,"operation_type_id":1,"amount":100.0`,
			invalidPayload:     true,
			expectedStatusCode: http.StatusUnprocessableEntity,
			expectedResponse:   `{"message":"Invalid JSON","details":["unexpected EOF"]}`,
		},
		{
			description: "Account not found",
			payload:     `{"account_id":1,"operation_type_id":1,"amount":100.0}`,
			expectedPayload: models.Transaction{
				AccountID: 1, OperationTypeID: 1, Amount: 100.0,
			},
			errorCreate:        entities.NewItemNotFoundError("Account", 1),
			expectedStatusCode: http.StatusNotFound,
			expectedResponse:   `{"message":"Account not found","details":["Account","1"]}`,
		},
		{
			description: "Unexpected error",
			payload:     `{"account_id":1,"operation_type_id":1,"amount":100.0}`,
			expectedPayload: models.Transaction{
				AccountID: 1, OperationTypeID: 1, Amount: 100.0,
			},
			errorCreate:        errors.New("Unexpected error"),
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `{"message":"Internal server error","details":["Unexpected error"]}`,
		},
	}

	for _, test := range tests {
		s.Run(test.description, func() {
			s.SetupTest()

			request := httptest.NewRequest(
				http.MethodPost,
				"/transactions",
				bytes.NewReader([]byte(test.payload)),
			)
			response := httptest.NewRecorder()

			if !test.invalidPayload {
				s.transactionServiceMock.EXPECT().Create(
					request.Context(), test.expectedPayload,
				).Return(test.responseCreate, test.errorCreate)
			}

			s.handler.ServeHTTP(response, request)

			s.Equal(test.expectedStatusCode, response.Code)
			if test.errorCreate != nil || test.invalidPayload {
				s.Equal(test.expectedResponse, response.Body.String())
			} else {
				s.Equal(response.Header().Get("ETag"), fmt.Sprint(transactionId))
			}
		})
	}
}
