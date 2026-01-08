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
		h.logger.Error("failed gewt all category on service")
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
		h.logger.Warn("Gagal decode JSON body", zap.Error(err))
		utils.ResponseBadRequest(w, http.StatusBadRequest, "Invalid JSON format", nil)
		return
	}
	newCategory := model.Category{
		Name:              req.Name,
		Rack_inventory_id: req.Rack_inventory_id,
	}
	err := h.service.CreateCategory(&newCategory)
	if err != nil {
		h.logger.Error("failed create category on service")
		utils.ResponseError(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  true,
		"message": "Create new Category succesfully",
	})
	utils.ResponseSuccess(w, http.StatusOK, "user ceated", nil)

}

func (h *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	//mengambil id
	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)

	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Invalid ID format (harus angka)", nil)
		return
	}
	var req dto.UpdateCategoryRequest
	//mengubah json body ke struct
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("Gagal decode JSON body", zap.Error(err))
		utils.ResponseBadRequest(w, http.StatusBadRequest, "Invalid JSON format", nil)
		return
	}

	newCategory := model.Category{
		Name:              req.Name,
		Rack_inventory_id: req.Rack_inventory_id,
	}
	err = h.service.UpdateCategory(id, &newCategory)
	if err != nil {
		h.logger.Error("failed update category on service")
		utils.ResponseError(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  true,
		"message": "Update category succesfully",
	})
	utils.ResponseSuccess(w, http.StatusOK, "new category created", nil)
}

func (h *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Invalid ID format (harus angka)", nil)
		return
	}

	err = h.service.DeleteCategory(id)
	if err != nil {
		h.logger.Error("failed delete category on service",
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
		"message": "Delete category succesfully ",
	})
	utils.ResponseSuccess(w, http.StatusOK, "delete category success", nil)

}
