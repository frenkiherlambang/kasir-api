package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"kasir-api/internal/domain"
	"kasir-api/internal/usecase"
)

type TransactionHandler struct {
	uc *usecase.TransactionUsecase
}

func NewTransactionHandler(uc *usecase.TransactionUsecase) *TransactionHandler {
	return &TransactionHandler{uc: uc}
}

// HandleCheckout handles POST /api/checkout. Body: {"items": [{"product_id": 1, "quantity": 2}, ...]}.
func (h *TransactionHandler) HandleCheckout(w http.ResponseWriter, r *http.Request) {
	var req domain.CheckoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if len(req.Items) == 0 {
		writeError(w, http.StatusBadRequest, "items required")
		return
	}
	tx, err := h.uc.Checkout(req.Items, false)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, tx)
}
