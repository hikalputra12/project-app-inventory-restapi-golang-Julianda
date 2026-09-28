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

type RackHandler struct {
	service service.RackServiceInterface
	logger  *zap.Logger
}

func NewRackHandler(service service.RackServiceInterface, log *zap.Logger) RackHandler {
	return RackHandler{
		service: service,
		logger:  log,
	}
}

func (h *RackHandler) ListRack(c *gin.Context) {
	page := utils.StringToInt(c.DefaultQuery("page", "1"), 1)
	limit := utils.StringToInt(c.DefaultQuery("limit", "10"), 10)

	racks, pagination, err := h.service.GetAllRack(c.Request.Context(), page, limit)
	if err != nil {
		h.logger.Error("Failed to get all racks", zap.Error(err))
		utils.ResponseError(c, http.StatusInternalServerError, "Gagal mengambil daftar rak: "+err.Error(), nil)
		return
	}

	var response []dto.RackListResponse
	for _, item := range racks {
		response = append(response, dto.RackListResponse{
			ID:                     item.ID,
			Name:                   item.Name,
			Warehouse_inventory_id: item.WarehouseInventoryId,
			WarehouseName:          item.WarehouseInventory,
		})
	}
	if response == nil {
		response = []dto.RackListResponse{}
	}

	utils.ResponsePagination(c, http.StatusOK, "Berhasil memuat daftar rak", response, *pagination)
}

func (h *RackHandler) CreateRack(c *gin.Context) {
	var req dto.CreateRackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Failed to decode JSON body", zap.Error(err))
		utils.ResponseBadRequest(c, http.StatusBadRequest, "Format JSON tidak valid", nil)
		return
	}

	if fieldErrors, err := utils.ValidateErrors(req); err != nil {
		utils.ResponseBadRequest(c, http.StatusBadRequest, "Validasi input gagal", fieldErrors)
		return
	}

	newRack := model.Rack{
		Name:                 req.Name,
		WarehouseInventoryId: req.Warehouse_inventory_id,
	}

	err := h.service.CreateRack(c.Request.Context(), &newRack)
	if err != nil {
		h.logger.Error("Failed to create rack", zap.Error(err))
		utils.ResponseError(c, http.StatusBadRequest, "Gagal membuat rak: "+err.Error(), nil)
		return
	}

	utils.ResponseSuccess(c, http.StatusCreated, "Rak baru berhasil dibuat", map[string]interface{}{
		"id":   newRack.ID,
		"name": newRack.Name,
	})
}

func (h *RackHandler) UpdateRack(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		utils.ResponseError(c, http.StatusBadRequest, "Format ID rak tidak valid", nil)
		return
	}

	var req dto.UpdateRackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Failed to decode JSON body", zap.Error(err))
		utils.ResponseBadRequest(c, http.StatusBadRequest, "Format JSON tidak valid", nil)
		return
	}

	if fieldErrors, err := utils.ValidateErrors(req); err != nil {
		utils.ResponseBadRequest(c, http.StatusBadRequest, "Validasi input gagal", fieldErrors)
		return
	}

	updateRack := model.Rack{
		Name:                 req.Name,
		WarehouseInventoryId: req.Warehouse_inventory_id,
	}

	err = h.service.UpdateRack(c.Request.Context(), id, &updateRack)
	if err != nil {
		h.logger.Error("Failed to update rack", zap.Int("id", id), zap.Error(err))
		utils.ResponseError(c, http.StatusBadRequest, "Gagal memperbarui rak: "+err.Error(), nil)
		return
	}

	utils.ResponseSuccess(c, http.StatusOK, "Rak berhasil diperbarui", map[string]int{"id": id})
}

func (h *RackHandler) DeleteRack(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		utils.ResponseError(c, http.StatusBadRequest, "Format ID rak tidak valid", nil)
		return
	}

	err = h.service.DeleteRack(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("Failed to delete rack", zap.Int("id", id), zap.Error(err))
		utils.ResponseError(c, http.StatusBadRequest, "Gagal menghapus rak: "+err.Error(), nil)
		return
	}

	utils.ResponseSuccess(c, http.StatusOK, "Rak berhasil dihapus", map[string]int{"id": id})
}

func (h *RackHandler) GetRackById(c *gin.Context) {
	idStr := c.Param("id")
	rackID, err := strconv.Atoi(idStr)
	if err != nil || rackID <= 0 {
		utils.ResponseBadRequest(c, http.StatusBadRequest, "Format ID rak tidak valid", nil)
		return
	}

	rack, err := h.service.GetRackById(c.Request.Context(), rackID)
	if err != nil {
		h.logger.Error("Failed to get rack by id", zap.Error(err), zap.Int("id", rackID))
		utils.ResponseError(c, http.StatusNotFound, "Rak tidak ditemukan", nil)
		return
	}

	response := dto.RackByIdResponse{
		ID:                   rack.ID,
		Name:                 rack.Name,
		WarehouseInventory:   rack.WarehouseInventory,
		WarehouseInventoryId: rack.WarehouseInventoryId,
	}
	utils.ResponseSuccess(c, http.StatusOK, "Berhasil mendapatkan data rak", response)
}
