package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"pismo/entities"
	"pismo/services"
)

func NewAccountGetByIDHandler(accountService services.AccountService) Handler {
	return &accountGetByIDHandler{accountService: accountService}
}

type accountGetByIDHandler struct {
	accountService services.AccountService
}

func (h *accountGetByIDHandler) Route() string {
	return "/accounts/{id}"
}

func (h *accountGetByIDHandler) Method() []string {
	return []string{http.MethodGet}
}

func (h *accountGetByIDHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	accountId, _ := strconv.Atoi(params["id"])
	account, err := h.accountService.FindByID(r.Context(), accountId)
	if err != nil {
		switch err.(type) {
		case *entities.ItemNotFoundError:
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusInternalServerError)
			err = entities.NewError("Error fetching account by ID", []string{err.Error()})
		}

		payloadJson, _ := json.Marshal(err)
		w.Write(payloadJson)
		return
	}

	payloadJson, _ := json.Marshal(account)
	w.Write(payloadJson)
}
