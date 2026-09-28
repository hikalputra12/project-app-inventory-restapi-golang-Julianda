package handler

import (
	"app-inventory/dto"
	"app-inventory/model"
	"app-inventory/service"
	"app-inventory/utils"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type TransactionHandler struct {
	service service.TransactionServiceInterface
	logger  *zap.Logger
}

func NewTransactionHandler(service service.TransactionServiceInterface, log *zap.Logger) TransactionHandler {
	return TransactionHandler{
		service: service,
		logger:  log,
	}
}

func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	var req dto.CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Failed to decode JSON body", zap.Error(err))
		utils.ResponseBadRequest(c, http.StatusBadRequest, "Format JSON tidak valid", nil)
		return
	}

	if fieldErrors, err := utils.ValidateErrors(req); err != nil {
		utils.ResponseBadRequest(c, http.StatusBadRequest, "Validasi input gagal", fieldErrors)
		return
	}

	// Retrieve authenticated user ID securely from Gin context
	userID, ok := utils.GetUserIDFromGin(c)
	if !ok || userID <= 0 {
		utils.ResponseError(c, http.StatusUnauthorized, "Pengguna tidak terautentikasi", nil)
		return
	}

	newTransaction := model.Transaction{
		UserId:      userID,
		InventoryId: req.InventoryId,
		Quantity:    req.Quantity,
	}

	err := h.service.CreateTransaction(c.Request.Context(), &newTransaction)
	if err != nil {
		h.logger.Error("Failed to create transaction", zap.Error(err))
		utils.ResponseError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	resp := dto.TransactionByIdResponse{
		ID:          newTransaction.ID,
		UserID:      newTransaction.UserId,
		InventoryID: newTransaction.InventoryId,
		Name:        newTransaction.Name,
		Quantity:    newTransaction.Quantity,
		Price:       newTransaction.Price,
		TotalPrice:  newTransaction.Quantity * newTransaction.Price,
		CreatedAt:   newTransaction.CreatedAt.Format(time.RFC3339),
	}

	utils.ResponseSuccess(c, http.StatusCreated, "Transaksi berhasil dibuat", resp)
}

func (h *TransactionHandler) ListTransaction(c *gin.Context) {
	page := utils.StringToInt(c.DefaultQuery("page", "1"), 1)
	limit := utils.StringToInt(c.DefaultQuery("limit", "10"), 10)

	transactions, pagination, err := h.service.GetAllTransaction(c.Request.Context(), page, limit)
	if err != nil {
		h.logger.Error("Failed to fetch transaction list", zap.Error(err))
		utils.ResponseError(c, http.StatusInternalServerError, "Gagal memuat daftar transaksi: "+err.Error(), nil)
		return
	}

	var response []dto.TransactionListResponse
	for _, item := range transactions {
		response = append(response, dto.TransactionListResponse{
			ID:          item.ID,
			UserID:      item.UserId,
			InventoryID: item.InventoryId,
			Name:        item.Name,
			Quantity:    item.Quantity,
			Price:       item.Price,
			TotalPrice:  item.Quantity * item.Price,
			CreatedAt:   item.CreatedAt.Format(time.RFC3339),
		})
	}
	if response == nil {
		response = []dto.TransactionListResponse{}
	}

	utils.ResponsePagination(c, http.StatusOK, "Berhasil memuat daftar transaksi", response, *pagination)
}

func (h *TransactionHandler) UpdateTransaction(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		utils.ResponseError(c, http.StatusBadRequest, "Format ID transaksi tidak valid", nil)
		return
	}

	var req dto.UpdateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseBadRequest(c, http.StatusBadRequest, "Format JSON tidak valid", nil)
		return
	}

	if fieldErrors, err := utils.ValidateErrors(req); err != nil {
		utils.ResponseBadRequest(c, http.StatusBadRequest, "Validasi input gagal", fieldErrors)
		return
	}

	newTransaction := model.Transaction{
		Quantity: req.Quantity,
	}

	err = h.service.UpdateTransaction(c.Request.Context(), id, &newTransaction)
	if err != nil {
		h.logger.Error("Failed to update transaction", zap.Int("id", id), zap.Error(err))
		utils.ResponseError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.ResponseSuccess(c, http.StatusOK, "Transaksi berhasil diperbarui", map[string]int{"id": id, "quantity": req.Quantity})
}

func (h *TransactionHandler) DeleteTransaction(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		utils.ResponseError(c, http.StatusBadRequest, "Format ID transaksi tidak valid", nil)
		return
	}

	err = h.service.DeleteTransaction(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("Failed to delete transaction", zap.Int("id", id), zap.Error(err))
		utils.ResponseError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.ResponseSuccess(c, http.StatusOK, fmt.Sprintf("Transaksi ID %d berhasil dibatalkan dan stok dikembalikan", id), nil)
}

func (h *TransactionHandler) GetTransactionById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		utils.ResponseBadRequest(c, http.StatusBadRequest, "Format ID transaksi tidak valid", nil)
		return
	}

	transaction, err := h.service.GetTransactionById(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("Failed to fetch transaction by id", zap.Error(err), zap.Int("id", id))
		utils.ResponseError(c, http.StatusNotFound, "Transaksi tidak ditemukan", nil)
		return
	}

	response := dto.TransactionByIdResponse{
		ID:          transaction.ID,
		UserID:      transaction.UserId,
		InventoryID: transaction.InventoryId,
		Name:        transaction.Name,
		Quantity:    transaction.Quantity,
		Price:       transaction.Price,
		TotalPrice:  transaction.Quantity * transaction.Price,
		CreatedAt:   transaction.CreatedAt.Format(time.RFC3339),
	}

	utils.ResponseSuccess(c, http.StatusOK, "Berhasil mendapatkan data transaksi", response)
}
