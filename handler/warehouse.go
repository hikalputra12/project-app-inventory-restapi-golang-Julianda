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
		utils.ResponseBadRequest(w, http.StatusInternalServerError, "Failed to fetch Warehouse: "+err.Error(), nil)
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
		utils.ResponseBadRequest(w, http.StatusBadRequest, "Invalid JSON format", nil)
		return
	}
	//pengambilan id melalui cookie
	cookie, _ := r.Cookie("session")

	// 4. Konversi ke Integer (jika ID Anda berupa angka)
	user_id, _ := strconv.Atoi(cookie.Value)

	newWarehouse := model.Warehouse{
		Name:    req.Name,
		User_id: user_id,
	}
	err := h.service.CreateWarehouse(&newWarehouse)
	if err != nil {
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
		utils.ResponseBadRequest(w, http.StatusBadRequest, "Invalid JSON format", nil)
		return
	}
	//pengambilan id melalui cookie
	cookie, _ := r.Cookie("session")

	// 4. Konversi ke Integer (jika ID Anda berupa angka)
	userId, _ := strconv.Atoi(cookie.Value)

	newWarehouse := model.Warehouse{
		Name:    req.Name,
		User_id: userId,
	}
	err = h.service.UpdateWarehouse(id, &newWarehouse)
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
