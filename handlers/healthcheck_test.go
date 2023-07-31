package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"
)

type healthcheckHandlerTestSuite struct {
	suite.Suite
	handler Handler
}

func TestHealthCheckHandlerTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(healthcheckHandlerTestSuite))
}

func (s *healthcheckHandlerTestSuite) SetupTest() {
	s.handler = NewHealthcheckHandler()
}

func (s *healthcheckHandlerTestSuite) TestRoute() {
	s.Equal("/", s.handler.Route())
}

func (s *healthcheckHandlerTestSuite) TestMethod() {
	s.Equal([]string{http.MethodGet}, s.handler.Method())
}

func (s *healthcheckHandlerTestSuite) TestServeHTTP() {
	s.T().Setenv("SERVICE_NAME", "hello-world")
	s.T().Setenv("VERSION", "1.0.0")

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()

	s.handler.ServeHTTP(response, request)

	s.Equal(http.StatusOK, response.Code)
	s.Equal(response.Body.String(), `{"service_name":"hello-world","version":"1.0.0"}`)
}
