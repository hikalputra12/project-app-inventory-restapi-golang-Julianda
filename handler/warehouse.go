package handler

import (
	"app-inventory/dto"
	"app-inventory/model"
	"app-inventory/service"
	"app-inventory/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type WarehouseHandler struct {
	service service.WarehouseServiceInterface
	logger  *zap.Logger
}

func NewWarehouseHandler(service service.WarehouseServiceInterface, log *zap.Logger) WarehouseHandler {
	return WarehouseHandler{
		service: service,
		logger:  log,
	}
}

func (h *WarehouseHandler) ListWarehouse(c *gin.Context) {
	page := utils.StringToInt(c.DefaultQuery("page", "1"), 1)
	limit := utils.StringToInt(c.DefaultQuery("limit", "10"), 10)

	warehouses, pagination, err := h.service.GetAllWarehouse(c.Request.Context(), page, limit)
	if err != nil {
		h.logger.Error("Failed to get warehouse list", zap.Error(err))
		utils.ResponseError(c, http.StatusInternalServerError, "Gagal mengambil daftar gudang: "+err.Error(), nil)
		return
	}

	var response []dto.WarehouseListResponse
	for _, item := range warehouses {
		response = append(response, dto.WarehouseListResponse{
			ID:       item.ID,
			Name:     item.Name,
			Location: item.Location,
		})
	}
	if response == nil {
		response = []dto.WarehouseListResponse{}
	}

	utils.ResponsePagination(c, http.StatusOK, "Berhasil memuat daftar gudang", response, *pagination)
}

func (h *WarehouseHandler) CreateWarehouse(c *gin.Context) {
	var req dto.CreateWarehouseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Failed to decode JSON body", zap.Error(err))
		utils.ResponseBadRequest(c, http.StatusBadRequest, "Format JSON tidak valid", nil)
		return
	}

	if fieldErrors, err := utils.ValidateErrors(req); err != nil {
		utils.ResponseBadRequest(c, http.StatusBadRequest, "Validasi input gagal", fieldErrors)
		return
	}

	newWarehouse := model.Warehouse{
		Name:     req.Name,
		Location: req.Location,
	}

	err := h.service.CreateWarehouse(c.Request.Context(), &newWarehouse)
	if err != nil {
		h.logger.Error("Failed to create warehouse", zap.Error(err))
		utils.ResponseError(c, http.StatusBadRequest, "Gagal membuat gudang: "+err.Error(), nil)
		return
	}

	utils.ResponseSuccess(c, http.StatusCreated, "Gudang baru berhasil dibuat", map[string]interface{}{
		"id":       newWarehouse.ID,
		"name":     newWarehouse.Name,
		"location": newWarehouse.Location,
	})
}

func (h *WarehouseHandler) UpdateWarehouse(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		utils.ResponseError(c, http.StatusBadRequest, "Format ID gudang tidak valid", nil)
		return
	}

	var req dto.UpdateWarehouseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Failed to decode JSON body", zap.Error(err))
		utils.ResponseBadRequest(c, http.StatusBadRequest, "Format JSON tidak valid", nil)
		return
	}

	if fieldErrors, err := utils.ValidateErrors(req); err != nil {
		utils.ResponseBadRequest(c, http.StatusBadRequest, "Validasi input gagal", fieldErrors)
		return
	}

	updateWarehouse := model.Warehouse{
		Name:     req.Name,
		Location: req.Location,
	}

	err = h.service.UpdateWarehouse(c.Request.Context(), id, &updateWarehouse)
	if err != nil {
		h.logger.Error("Failed to update warehouse", zap.Int("id", id), zap.Error(err))
		utils.ResponseError(c, http.StatusBadRequest, "Gagal memperbarui gudang: "+err.Error(), nil)
		return
	}

	utils.ResponseSuccess(c, http.StatusOK, "Gudang berhasil diperbarui", map[string]int{"id": id})
}

func (h *WarehouseHandler) DeleteWarehouse(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		utils.ResponseError(c, http.StatusBadRequest, "Format ID gudang tidak valid", nil)
		return
	}

	err = h.service.DeleteWarehouse(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("Failed to delete warehouse", zap.Int("id", id), zap.Error(err))
		utils.ResponseError(c, http.StatusBadRequest, "Gagal menghapus gudang: "+err.Error(), nil)
		return
	}

	utils.ResponseSuccess(c, http.StatusOK, "Gudang berhasil dihapus", map[string]int{"id": id})
}

func (h *WarehouseHandler) GetWarehouseById(c *gin.Context) {
	idStr := c.Param("id")
	warehouseID, err := strconv.Atoi(idStr)
	if err != nil || warehouseID <= 0 {
		utils.ResponseBadRequest(c, http.StatusBadRequest, "Format ID gudang tidak valid", nil)
		return
	}

	warehouse, err := h.service.GetWarehouseById(c.Request.Context(), warehouseID)
	if err != nil {
		h.logger.Error("Failed to get warehouse by id", zap.Error(err), zap.Int("id", warehouseID))
		utils.ResponseError(c, http.StatusNotFound, "Gudang tidak ditemukan", nil)
		return
	}

	response := dto.WarehouseByIdResponse{
		ID:       warehouse.ID,
		Name:     warehouse.Name,
		Location: warehouse.Location,
	}
	utils.ResponseSuccess(c, http.StatusOK, "Berhasil mendapatkan data gudang", response)
}
