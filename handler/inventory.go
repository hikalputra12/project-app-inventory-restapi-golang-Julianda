package handler

import (
	"app-inventory/dto"
	"app-inventory/service"
	"app-inventory/utils"
	"net/http"
	"strconv"

	"go.uber.org/zap"
)

type InventoryHandler struct {
	service service.InventoryServiceInterface
	logger  *zap.Logger
}

// constructor
func NewInventoryHandler(service service.InventoryServiceInterface, log *zap.Logger) InventoryHandler {
	return InventoryHandler{
		service: service,
		logger:  log,
	}
}

func (h *InventoryHandler) ListInventory(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "Invalid page", nil)
		return
	}

	// config limit pagination
	limit := 3

	// Get data Inventorys form service all Inventorys
	Inventories, pagination, err := h.service.GetAllInventory(page, limit)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusInternalServerError, "Failed to fetch inventory: "+err.Error(), nil)
		return
	}
	var response []dto.InventoryListResponse
	for _, item := range Inventories {
		response = append(response, dto.InventoryListResponse{
			Name:      item.Name,
			Price:     item.Price,
			Stock:     item.Stock,
			Category:  item.Category,
			Rack:      item.Rack,
			Warehouse: item.Warehouse,
		})

	}
	utils.ResponsePagination(w, http.StatusOK, "success get data", response, *pagination)

}
