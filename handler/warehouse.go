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

type WarehouseHandler struct {
	service service.WarehouseServiceInterface
	logger  *zap.Logger
}

// constructor
func NewWarehouseHandler(service service.WarehouseServiceInterface, log *zap.Logger) WarehouseHandler {
	return WarehouseHandler{
		service: service,
		logger:  log,
	}
}

func (h *WarehouseHandler) ListWarehouse(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "Invalid page", nil)
		return
	}

	// config limit pagination
	limit := 3

	// Get data Warehouses form service all Warehouses
	Warehouse, pagination, err := h.service.GetAllWarehouse(page, limit)
	if err != nil {
		h.logger.Error("failed get list warehouse on service",
			zap.Error(err),
		)
		utils.ResponseBadRequest(w, http.StatusBadRequest, "Failed to fetch Warehouse: "+err.Error(), nil)
		return
	}
	var response []dto.WarehouseListResponse
	for _, item := range Warehouse {
		response = append(response, dto.WarehouseListResponse{
			Name: item.Name,
		})

	}
	utils.ResponsePagination(w, http.StatusOK, "success get data", response, *pagination)

}

func (h *WarehouseHandler) CreateWarehouse(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateWarehouseRequest
	//mengubah json body ke struct
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("Gagal decode JSON body", zap.Error(err))
		utils.ResponseBadRequest(w, http.StatusBadRequest, "Invalid JSON format", nil)
		return
	}

	newWarehouse := model.Warehouse{
		Name:     req.Name,
		Location: req.Location,
	}
	err := h.service.CreateWarehouse(&newWarehouse)
	if err != nil {
		h.logger.Error("failed create warehouse on service",
			zap.Error(err),
		)
		utils.ResponseError(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  true,
		"message": "Create new Warehouse succesfully",
	})

}

func (h *WarehouseHandler) UpdateWarehouse(w http.ResponseWriter, r *http.Request) {
	//mengambil id
	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)

	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Invalid ID format (harus angka)", nil)
		return
	}
	var req dto.UpdateWarehouseRequest
	//mengubah json body ke struct
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("Gagal decode JSON body", zap.Error(err))
		utils.ResponseBadRequest(w, http.StatusBadRequest, "Invalid JSON format", nil)
		return
	}

	newWarehouse := model.Warehouse{
		Name:     req.Name,
		Location: req.Location,
	}
	err = h.service.UpdateWarehouse(id, &newWarehouse)
	if err != nil {
		h.logger.Error("failed update warehouse on service",
			zap.Error(err),
		)
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

func (h *WarehouseHandler) DeleteWarehouse(w http.ResponseWriter, r *http.Request) {
	//mengambil id
	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)

	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Invalid ID format (harus angka)", nil)
		return
	}
	err = h.service.DeleteWarehouse(id)
	if err != nil {
		h.logger.Error("failed delete warehouse on service",
			zap.String("warehouse_id", idStr),
			zap.Error(err),
		)
		utils.ResponseError(w, http.StatusBadRequest, err.Error(), nil)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  true,
		"message": "Delete user succesfully",
	})

}

// get warehouse by id
func (h *WarehouseHandler) GetWarehouseById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	WarehouseID, _ := strconv.Atoi(id)

	// Get data Racks form service all Racks
	warehouse, err := h.service.GetWarehouseById(WarehouseID)
	if err != nil {
		h.logger.Error("failed get warehouse by id on service",
			zap.Error(err),
		)
		utils.ResponseBadRequest(w, http.StatusBadRequest, "Failed to fetch assignments: "+err.Error(), nil)
		return
	}
	var response dto.WarehouseByIdResponse
	response = dto.WarehouseByIdResponse{
		Name:     warehouse.Name,
		Location: warehouse.Location,
	}
	utils.ResponseSuccess(w, http.StatusOK, "success get data", response)
}
