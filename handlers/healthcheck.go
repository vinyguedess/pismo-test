package handlers

import (
	"encoding/json"
	"net/http"
	"os"
)

type healthcheckHandler struct{}

func NewHealthcheckHandler() Handler {
	return &healthcheckHandler{}
}

func (h *healthcheckHandler) Route() string {
	return "/"
}

func (h *healthcheckHandler) Method() []string {
	return []string{http.MethodGet}
}

func (h *healthcheckHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	payload := map[string]string{
		"service_name": os.Getenv("SERVICE_NAME"),
		"version":      os.Getenv("VERSION"),
	}

	payloadJson, _ := json.Marshal(payload)
	w.Write(payloadJson)
}
