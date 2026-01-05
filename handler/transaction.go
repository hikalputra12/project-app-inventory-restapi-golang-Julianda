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

func (h *TransactionHandler) ListTransaction(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "Invalid page", nil)
		return
	}

	// config limit pagination
	limit := 10

	// Get data Transactions form service all Transactions
	Transaction, pagination, err := h.service.GetAllTransaction(page, limit)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusInternalServerError, "Failed to fetch Transaction: "+err.Error(), nil)
		return
	}
	var response []dto.TransactionListResponse
	for _, item := range Transaction {
		response = append(response, dto.TransactionListResponse{
			Name:     item.Name,
			Quantity: item.Quantity,
			Price:    item.Price,
		})

	}
	utils.ResponsePagination(w, http.StatusOK, "success get data", response, *pagination)

}

func (h *TransactionHandler) UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateTransactionRequest
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
		UserId:   userId,
		Quantity: req.Quantity,
	}
	err := h.service.UpdateTransaction(req.SalesItemID, &newTransaction)
	if err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  true,
		"message": "Update transaction succesfully",
	})

}

func (h *TransactionHandler) DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	var req dto.DeleteTransactionRequest
	//mengubah json body ke struct
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "Invalid JSON format", nil)
		return
	}

	err := h.service.DeleteTransaction(req.SalesItemID)
	if err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  true,
		"message": "Delete user succesfully with ID:" + strconv.Itoa(req.SalesItemID),
	})

}
