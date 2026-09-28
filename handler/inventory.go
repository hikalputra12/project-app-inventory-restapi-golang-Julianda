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

type InventoryHandler struct {
	service service.InventoryServiceInterface
	logger  *zap.Logger
}

func NewInventoryHandler(service service.InventoryServiceInterface, log *zap.Logger) InventoryHandler {
	return InventoryHandler{
		service: service,
		logger:  log,
	}
}

func (h *InventoryHandler) ListInventory(c *gin.Context) {
	page := utils.StringToInt(c.DefaultQuery("page", "1"), 1)
	limit := utils.StringToInt(c.DefaultQuery("limit", "10"), 10)

	inventories, pagination, err := h.service.GetAllInventory(c.Request.Context(), page, limit)
	if err != nil {
		h.logger.Error("Failed to get all inventory on service", zap.Error(err))
		utils.ResponseError(c, http.StatusInternalServerError, "Gagal mengambil daftar inventaris: "+err.Error(), nil)
		return
	}

	var response []dto.InventoryListResponse
	for _, item := range inventories {
		response = append(response, dto.InventoryListResponse{
			ID:                  item.ID,
			Name:                item.Name,
			Price:               item.Price,
			Stock:               item.Stock,
			Category:            item.Category,
			Rack:                item.Rack,
			Warehouse:           item.Warehouse,
			CategoryInventoryID: item.Category_inventory_id,
		})
	}
	if response == nil {
		response = []dto.InventoryListResponse{}
	}

	utils.ResponsePagination(c, http.StatusOK, "Berhasil memuat data inventaris", response, *pagination)
}

func (h *InventoryHandler) CreateInventory(c *gin.Context) {
	var req dto.CreateInventoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Failed to decode JSON body", zap.Error(err))
		utils.ResponseBadRequest(c, http.StatusBadRequest, "Format JSON tidak valid", nil)
		return
	}

	if fieldErrors, err := utils.ValidateErrors(req); err != nil {
		utils.ResponseBadRequest(c, http.StatusBadRequest, "Validasi input gagal", fieldErrors)
		return
	}

	newInventory := model.Inventory{
		Name:                  req.Name,
		Price:                 req.Price,
		Stock:                 req.Stock,
		Category_inventory_id: req.Category_id,
	}

	err := h.service.CreateInventory(c.Request.Context(), &newInventory)
	if err != nil {
		h.logger.Error("Failed to create inventory", zap.Error(err))
		utils.ResponseError(c, http.StatusInternalServerError, "Gagal membuat inventaris: "+err.Error(), nil)
		return
	}

	utils.ResponseSuccess(c, http.StatusCreated, "Inventaris baru berhasil dibuat", map[string]interface{}{
		"id":    newInventory.ID,
		"name":  newInventory.Name,
		"stock": newInventory.Stock,
		"price": newInventory.Price,
	})
}

func (h *InventoryHandler) UpdateInventory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		utils.ResponseError(c, http.StatusBadRequest, "Format ID inventaris tidak valid", nil)
		return
	}

	var req dto.UpdateInventoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Failed to decode JSON body", zap.Error(err))
		utils.ResponseBadRequest(c, http.StatusBadRequest, "Format JSON tidak valid", nil)
		return
	}

	if fieldErrors, err := utils.ValidateErrors(req); err != nil {
		utils.ResponseBadRequest(c, http.StatusBadRequest, "Validasi input gagal", fieldErrors)
		return
	}

	err = h.service.UpdateInventory(c.Request.Context(), id, &req)
	if err != nil {
		h.logger.Error("Failed to update inventory", zap.Int("id", id), zap.Error(err))
		utils.ResponseError(c, http.StatusBadRequest, "Gagal memperbarui inventaris: "+err.Error(), nil)
		return
	}

	utils.ResponseSuccess(c, http.StatusOK, "Inventaris berhasil diperbarui", map[string]int{"id": id})
}

func (h *InventoryHandler) DeleteInventory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		utils.ResponseError(c, http.StatusBadRequest, "Format ID inventaris tidak valid", nil)
		return
	}

	err = h.service.DeleteInventory(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("Failed to delete inventory", zap.Int("id", id), zap.Error(err))
		utils.ResponseError(c, http.StatusBadRequest, "Gagal menghapus inventaris: "+err.Error(), nil)
		return
	}

	utils.ResponseSuccess(c, http.StatusOK, "Inventaris berhasil dihapus", map[string]int{"id": id})
}

func (h *InventoryHandler) CheckStock(c *gin.Context) {
	page := utils.StringToInt(c.DefaultQuery("page", "1"), 1)
	limit := utils.StringToInt(c.DefaultQuery("limit", "10"), 10)

	inventories, pagination, err := h.service.CheckStock(c.Request.Context(), page, limit)
	if err != nil {
		h.logger.Error("Failed to check low stock", zap.Error(err))
		utils.ResponseError(c, http.StatusInternalServerError, "Gagal memeriksa stok menipis: "+err.Error(), nil)
		return
	}

	var response []dto.InventoryListResponse
	for _, item := range inventories {
		response = append(response, dto.InventoryListResponse{
			ID:                  item.ID,
			Name:                item.Name,
			Price:               item.Price,
			Stock:               item.Stock,
			Category:            item.Category,
			Rack:                item.Rack,
			Warehouse:           item.Warehouse,
			CategoryInventoryID: item.Category_inventory_id,
		})
	}
	if response == nil {
		response = []dto.InventoryListResponse{}
	}

	utils.ResponsePagination(c, http.StatusOK, "Berhasil memuat data stok menipis", response, *pagination)
}

func (h *InventoryHandler) GetInventoryById(c *gin.Context) {
	idStr := c.Param("id")
	inventoryID, err := strconv.Atoi(idStr)
	if err != nil || inventoryID <= 0 {
		utils.ResponseBadRequest(c, http.StatusBadRequest, "Format ID inventaris tidak valid", nil)
		return
	}

	inventory, err := h.service.GetInventoryById(c.Request.Context(), inventoryID)
	if err != nil {
		h.logger.Error("Failed to get inventory by id", zap.Error(err), zap.Int("id", inventoryID))
		utils.ResponseError(c, http.StatusNotFound, "Inventaris tidak ditemukan", nil)
		return
	}

	response := dto.InventoryByIdResponse{
		ID:                  inventory.ID,
		Name:                inventory.Name,
		Price:               inventory.Price,
		Stock:               inventory.Stock,
		CategoryInventoryID: inventory.Category_inventory_id,
		Category:            inventory.Category,
		Rack:                inventory.Rack,
		Warehouse:           inventory.Warehouse,
	}
	utils.ResponseSuccess(c, http.StatusOK, "Berhasil mendapatkan data inventaris", response)
}
