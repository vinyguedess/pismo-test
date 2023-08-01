package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"pismo/entities"
	"pismo/models"
	"pismo/services"
)

type createAccountHandler struct {
	accountService services.AccountService
}

func NewCreateAccountHandler(accountService services.AccountService) Handler {
	return &createAccountHandler{
		accountService: accountService,
	}
}

func (h *createAccountHandler) Route() string {
	return "/accounts"
}

func (h *createAccountHandler) Method() []string {
	return []string{http.MethodPost}
}

func (h *createAccountHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var payload models.Account
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		err = entities.NewError("Invalid JSON", []string{err.Error()})

		jsonPayload, _ := json.Marshal(err)
		w.Write(jsonPayload)
		return
	}

	account, err := h.accountService.Create(r.Context(), payload)
	if err != nil {
		switch err.(type) {
		case *entities.AccountAlreadyExistsError:
			w.WriteHeader(http.StatusConflict)
		default:
			w.WriteHeader(http.StatusInternalServerError)
			err = entities.NewError("Internal server error", []string{err.Error()})
		}

		jsonPayload, _ := json.Marshal(err)
		w.Write(jsonPayload)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("ETag", fmt.Sprint(account.ID))
	w.Header().Set("Location", fmt.Sprintf("/accounts/%d", account.ID))
}
