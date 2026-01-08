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

type RackHandler struct {
	service service.RackServiceInterface
	logger  *zap.Logger
}

// constructor
func NewRackHandler(service service.RackServiceInterface, log *zap.Logger) RackHandler {
	return RackHandler{
		service: service,
		logger:  log,
	}
}

func (h *RackHandler) ListRack(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "Invalid page", nil)
		return
	}

	// config limit pagination
	limit := 3

	// Get data Racks form service all Racks
	rack, pagination, err := h.service.GetAllRack(page, limit)
	if err != nil {
		h.logger.Error("failed get all rack on service",
			zap.Error(err),
		)
		utils.ResponseBadRequest(w, http.StatusInternalServerError, "Failed to fetch Rack: "+err.Error(), nil)
		return
	}
	var response []dto.RackListResponse
	for _, item := range rack {
		response = append(response, dto.RackListResponse{
			Name:                   item.Name,
			Warehouse_inventory_id: item.Warehouse_inventory_id,
		})

	}
	utils.ResponsePagination(w, http.StatusOK, "success get data", response, *pagination)

}

func (h *RackHandler) CreateRack(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateRackRequest
	//mengubah json body ke struct
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("Gagal decode JSON body", zap.Error(err))
		utils.ResponseBadRequest(w, http.StatusBadRequest, "Invalid JSON format", nil)
		return
	}
	newRack := model.Rack{
		Name:                   req.Name,
		Warehouse_inventory_id: req.Warehouse_inventory_id,
	}
	err := h.service.CreateRack(&newRack)
	if err != nil {
		h.logger.Error("failed create rack on service",
			zap.Error(err),
		)
		utils.ResponseError(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  true,
		"message": "Create new Rack succesfully",
	})

}

func (h *RackHandler) UpdateRack(w http.ResponseWriter, r *http.Request) {
	//mengambil id
	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)

	if err != nil {
		h.logger.Error("failed update rack on service",
			zap.String("user_id", idStr),
			zap.Error(err),
		)
		utils.ResponseError(w, http.StatusBadRequest, "Invalid ID format (harus angka)", nil)
		return
	}
	var req dto.UpdateRackRequest
	//mengubah json body ke struct
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("Gagal decode JSON body", zap.Error(err))
		utils.ResponseBadRequest(w, http.StatusBadRequest, "Invalid JSON format", nil)
		return
	}

	newRack := model.Rack{
		Name:                   req.Name,
		Warehouse_inventory_id: req.Warehouse_inventory_id,
	}
	err = h.service.UpdateRack(id, &newRack)
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

func (h *RackHandler) DeleteRack(w http.ResponseWriter, r *http.Request) {
	//mengambil id
	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)

	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Invalid ID format (harus angka)", nil)
		return
	}
	err = h.service.DeleteRack(id)
	if err != nil {
		h.logger.Error("failed delete rack on service",
			zap.String("user_id", idStr),
			zap.Error(err),
		)
		utils.ResponseError(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  true,
		"message": "Delete rack succesfully ",
	})

}
