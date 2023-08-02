package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/suite"

	"pismo/entities"
	mock_services "pismo/mocks/services"
	"pismo/models"
)

func TestAccountGetByIDHandlerTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(accountGetByIdHandlerTestSuite))
}

type accountGetByIdHandlerTestSuite struct {
	suite.Suite
	ctrl               *gomock.Controller
	accountServiceMock *mock_services.MockAccountService
	handler            Handler
}

func (s *accountGetByIdHandlerTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.accountServiceMock = mock_services.NewMockAccountService(s.ctrl)
	s.handler = NewAccountGetByIDHandler(s.accountServiceMock)
}

func (s *accountGetByIdHandlerTestSuite) TestRoute() {
	s.Equal("/accounts/{id}", s.handler.Route())
}

func (s *accountGetByIdHandlerTestSuite) TestMethod() {
	s.Equal([]string{http.MethodGet}, s.handler.Method())
}

func (s *accountGetByIdHandlerTestSuite) TestServeHTTP() {
	accountId := 1
	createdAt := time.Now().UTC()

	account := &models.Account{
		BaseModel: models.BaseModel{
			ID:        accountId,
			CreatedAt: createdAt,
		},
		DocumentNumber: "123456789",
	}

	tests := []struct {
		description        string
		responseFindByID   *models.Account
		errFindByID        error
		expectedResponse   string
		expectedStatusCode int
	}{
		{
			description:      "Success",
			responseFindByID: account,
			expectedResponse: fmt.Sprintf(
				`{"id":1,"created_at":"%s","document_number":"123456789"}`,
				createdAt.Format("2006-01-02T15:04:05.000000000Z"),
			),
			expectedStatusCode: http.StatusOK,
		},
		{
			description:        "Account not found",
			errFindByID:        entities.NewItemNotFoundError("Account", accountId),
			expectedResponse:   `{"message":"Account not found","details":["Account","1"]}`,
			expectedStatusCode: http.StatusNotFound,
		},
		{
			description:        "Internal server error",
			errFindByID:        errors.New("Unknown error"),
			expectedResponse:   `{"message":"Error fetching account by ID","details":["Unknown error"]}`,
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		s.Run(test.description, func() {
			request := mux.SetURLVars(
				httptest.NewRequest(http.MethodGet, fmt.Sprintf("/accounts/%d", accountId), nil),
				map[string]string{"id": fmt.Sprintf("%d", accountId)},
			)
			response := httptest.NewRecorder()

			s.accountServiceMock.EXPECT().FindByID(
				request.Context(), accountId,
			).Return(test.responseFindByID, test.errFindByID)

			s.handler.ServeHTTP(response, request)

			s.Equal(test.expectedStatusCode, response.Code)
			s.Equal(test.expectedResponse, response.Body.String())
		})
	}
}
