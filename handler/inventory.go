package handler

import (
	"app-inventory/dto"
	"app-inventory/model"
	"app-inventory/service"
	"app-inventory/utils"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
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

func (h *InventoryHandler) CreateInventory(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateInventoryRequest
	//mengubah json body ke struct
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "Invalid JSON format", nil)
		return
	}
	newInventory := model.Inventory{
		Name:                  req.Name,
		Price:                 req.Price,
		Stock:                 req.Stock,
		Category_inventory_id: req.Category_id,
	}
	err := h.service.CreateInventory(&newInventory)
	if err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  true,
		"message": "Create new inventory succesfully",
	})

}

func (h *InventoryHandler) UpdateInventory(w http.ResponseWriter, r *http.Request) {
	//mengambil id
	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)

	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Invalid ID format (harus angka)", nil)
		return
	}
	var req dto.UpdateInventoryRequest
	//mengubah json body ke struct
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "Invalid JSON format", nil)
		return
	}

	newInventory := model.Inventory{
		Name:                  req.Name,
		Price:                 req.Price,
		Stock:                 req.Stock,
		Category_inventory_id: req.Category_id,
	}
	err = h.service.UpdateInventory(id, &newInventory)
	if err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  true,
		"message": "Update user succesfully",
	})

}

func (h *InventoryHandler) DeleteInventory(w http.ResponseWriter, r *http.Request) {
	//mengambil id
	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)

	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Invalid ID format (harus angka)", nil)
		return
	}

	err = h.service.DeleteInventory(id)
	if err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  true,
		"message": "Delete user succesfully",
	})

}

func (h *InventoryHandler) CheckStock(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "Invalid page", nil)
		return
	}

	// config limit pagination
	limit := 3

	// Get data Inventorys form service all Inventorys
	Inventories, pagination, err := h.service.CheckStock(page, limit)
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
