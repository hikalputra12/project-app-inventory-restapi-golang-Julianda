package handler

import (
	"app-inventory/dto"
	"app-inventory/model"
	"app-inventory/service"
	"app-inventory/utils"
	"encoding/json"
	"net/http"
	"strconv"

	"go.uber.org/zap"
)

type TransactionHandler struct {
	service service.TransactionServiceInterface
	logger  *zap.Logger
}

// constructor
func NewTransactionHandler(service service.TransactionServiceInterface, log *zap.Logger) TransactionHandler {
	return TransactionHandler{
		service: service,
		logger:  log,
	}
}

func (h *TransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateTransactionRequest
	//mengubah json body ke struct
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "Invalid JSON format", nil)
		return
	}
	//pengambilan id melalui cookie
	cookie, _ := r.Cookie("session")

	// 4. Konversi ke Integer (jika ID Anda berupa angka)
	userId, _ := strconv.Atoi(cookie.Value)
	newTransaction := model.Transaction{
		UserId:      userId,
		InventoryId: req.InventoryId,
		Quantity:    req.Quantity,
	}
	err := h.service.CreateTransaction(&newTransaction)
	if err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  true,
		"message": "Create new Transaction succesfully",
	})

}
