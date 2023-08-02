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

type createAccountHandlerTestSuite struct {
	suite.Suite
	ctrl            *gomock.Controller
	authServiceMock *mock_services.MockAccountService
	handler         Handler
}

func TestCreateAccountHandlerTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(createAccountHandlerTestSuite))
}

func (s *createAccountHandlerTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.authServiceMock = mock_services.NewMockAccountService(s.ctrl)
	s.handler = NewAccountCreateHandler(s.authServiceMock)
}

func (s *createAccountHandlerTestSuite) TestRoute() {
	s.Equal("/accounts", s.handler.Route())
}

func (s *createAccountHandlerTestSuite) TestMethod() {
	s.Equal([]string{http.MethodPost}, s.handler.Method())
}

func (s *createAccountHandlerTestSuite) TestServeHTTP() {
	accountId := 1
	account := &models.Account{
		BaseModel: models.BaseModel{
			ID: accountId,
		},
	}

	tests := []struct {
		description             string
		payload                 string
		invalidPayloadErr       bool
		responseCreate          *models.Account
		errCreate               error
		accountAlreadyExistsErr bool
		expectedResponse        string
		expectedStatusCode      int
	}{
		{
			description:        "Success",
			payload:            `{"document_number": "123456789"}`,
			responseCreate:     account,
			expectedStatusCode: http.StatusCreated,
		},
		{
			description:        "Invalid payload",
			payload:            `{"document_number": "123456789"`,
			invalidPayloadErr:  true,
			expectedResponse:   `{"message":"Invalid JSON","details":["unexpected EOF"]}`,
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			description:             "Account already exists",
			payload:                 `{"document_number": "123456789"}`,
			errCreate:               entities.NewAccountAlreadyExistsError("123456789"),
			accountAlreadyExistsErr: true,
			expectedResponse:        `{"message":"Account already exists","details":["123456789"]}`,
			expectedStatusCode:      http.StatusConflict,
		},
		{
			description:        "Internal server error",
			payload:            `{"document_number": "123456789"}`,
			errCreate:          errors.New("Unexpected error"),
			expectedResponse:   `{"message":"Internal server error","details":["Unexpected error"]}`,
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		s.Run(test.description, func() {
			request := httptest.NewRequest(http.MethodPost, "/accounts", bytes.NewReader([]byte(test.payload)))
			response := httptest.NewRecorder()

			if !test.invalidPayloadErr {
				s.authServiceMock.EXPECT().
					Create(request.Context(), models.Account{DocumentNumber: "123456789"}).
					Return(test.responseCreate, test.errCreate)
			}

			s.handler.ServeHTTP(response, request)

			s.Equal(test.expectedStatusCode, response.Code)
			if test.responseCreate == nil {
				s.Equal(test.expectedResponse, response.Body.String())
			} else {
				s.Equal(fmt.Sprint(accountId), response.Header().Get("ETag"))
				s.Equal(fmt.Sprintf("/accounts/%d", accountId), response.Header().Get("Location"))
			}
		})
	}
}
