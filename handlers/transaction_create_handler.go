package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"pismo/entities"
	"pismo/models"
	"pismo/services"
)

func NewTransactionCreateHandler(
	transactionService services.TransactionService,
) Handler {
	return &transactionCreateHandler{
		transactionService: transactionService,
	}
}

type transactionCreateHandler struct {
	transactionService services.TransactionService
}

func (h *transactionCreateHandler) Route() string {
	return "/transactions"
}

func (h *transactionCreateHandler) Method() []string {
	return []string{http.MethodPost}
}

func (h *transactionCreateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var payload models.Transaction
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		err = entities.NewError("Invalid JSON", []string{err.Error()})

		jsonPayload, _ := json.Marshal(err)
		w.Write(jsonPayload)
		return
	}

	transaction, err := h.transactionService.Create(r.Context(), payload)
	if err != nil {
		switch err.(type) {
		case *entities.ItemNotFoundError:
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusInternalServerError)
			err = entities.NewError("Internal server error", []string{err.Error()})
		}

		payloadJson, _ := json.Marshal(err)
		w.Write(payloadJson)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("ETag", fmt.Sprint(transaction.ID))
}
