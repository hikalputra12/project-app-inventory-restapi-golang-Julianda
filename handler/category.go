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

type CategoryHandler struct {
	service service.CategoryServiceInterface
	logger  *zap.Logger
}

// constructor
func NewCategoryHandler(service service.CategoryServiceInterface, log *zap.Logger) CategoryHandler {
	return CategoryHandler{
		service: service,
		logger:  log,
	}
}

func (h *CategoryHandler) ListCategory(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "Invalid page", nil)
		return
	}

	// config limit pagination
	limit := 3

	// Get data Categorys form service all Categorys
	Categories, pagination, err := h.service.GetAllCategory(page, limit)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusInternalServerError, "Failed to fetch Category: "+err.Error(), nil)
		return
	}
	var response []dto.CategoryListResponse
	for _, item := range Categories {
		response = append(response, dto.CategoryListResponse{
			Name:              item.Name,
			Rack_inventory_id: item.Rack_inventory_id,
		})

	}
	utils.ResponsePagination(w, http.StatusOK, "success get data", response, *pagination)

}

func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateCategoryRequest
	//mengubah json body ke struct
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "Invalid JSON format", nil)
		return
	}
	newCategory := model.Category{
		Name:              req.Name,
		Rack_inventory_id: req.Rack_inventory_id,
	}
	err := h.service.CreateCategory(&newCategory)
	if err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  true,
		"message": "Create new Category succesfully",
	})

}

func (h *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateCategoryRequest
	//mengubah json body ke struct
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "Invalid JSON format", nil)
		return
	}

	newCategory := model.Category{
		Name:              req.Name,
		Rack_inventory_id: req.Rack_inventory_id,
	}
	err := h.service.UpdateCategory(req.Category_Inventory_id, &newCategory)
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

func (h *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	var req dto.DeleteCategoryRequest
	//mengubah json body ke struct
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "Invalid JSON format", nil)
		return
	}

	err := h.service.DeleteCategory(req.Category_Inventory_id)
	if err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  true,
		"message": "Delete user succesfully with ID:" + strconv.Itoa(req.Category_Inventory_id),
	})

}
