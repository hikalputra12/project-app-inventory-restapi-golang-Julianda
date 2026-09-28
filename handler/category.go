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

type CategoryHandler struct {
	service service.CategoryServiceInterface
	logger  *zap.Logger
}

func NewCategoryHandler(service service.CategoryServiceInterface, log *zap.Logger) CategoryHandler {
	return CategoryHandler{
		service: service,
		logger:  log,
	}
}

func (h *CategoryHandler) ListCategory(c *gin.Context) {
	page := utils.StringToInt(c.DefaultQuery("page", "1"), 1)
	limit := utils.StringToInt(c.DefaultQuery("limit", "10"), 10)

	categories, pagination, err := h.service.GetAllCategory(c.Request.Context(), page, limit)
	if err != nil {
		h.logger.Error("Failed to get all categories", zap.Error(err))
		utils.ResponseError(c, http.StatusInternalServerError, "Gagal mengambil daftar kategori: "+err.Error(), nil)
		return
	}

	var response []dto.CategoryListResponse
	for _, item := range categories {
		response = append(response, dto.CategoryListResponse{
			ID:                item.ID,
			Name:              item.Name,
			Rack_inventory_id: item.Rack_inventory_id,
			RackName:          item.RackInventory,
		})
	}
	if response == nil {
		response = []dto.CategoryListResponse{}
	}

	utils.ResponsePagination(c, http.StatusOK, "Berhasil memuat daftar kategori", response, *pagination)
}

func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req dto.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Failed to decode JSON body", zap.Error(err))
		utils.ResponseBadRequest(c, http.StatusBadRequest, "Format JSON tidak valid", nil)
		return
	}

	if fieldErrors, err := utils.ValidateErrors(req); err != nil {
		utils.ResponseBadRequest(c, http.StatusBadRequest, "Validasi input gagal", fieldErrors)
		return
	}

	newCategory := model.Category{
		Name:              req.Name,
		Rack_inventory_id: req.Rack_inventory_id,
	}

	err := h.service.CreateCategory(c.Request.Context(), &newCategory)
	if err != nil {
		h.logger.Error("Failed to create category", zap.Error(err))
		utils.ResponseError(c, http.StatusBadRequest, "Gagal membuat kategori: "+err.Error(), nil)
		return
	}

	utils.ResponseSuccess(c, http.StatusCreated, "Kategori baru berhasil dibuat", map[string]interface{}{
		"id":   newCategory.ID,
		"name": newCategory.Name,
	})
}

func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		utils.ResponseError(c, http.StatusBadRequest, "Format ID kategori tidak valid", nil)
		return
	}

	var req dto.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Failed to decode JSON body", zap.Error(err))
		utils.ResponseBadRequest(c, http.StatusBadRequest, "Format JSON tidak valid", nil)
		return
	}

	if fieldErrors, err := utils.ValidateErrors(req); err != nil {
		utils.ResponseBadRequest(c, http.StatusBadRequest, "Validasi input gagal", fieldErrors)
		return
	}

	updateCat := model.Category{
		Name:              req.Name,
		Rack_inventory_id: req.Rack_inventory_id,
	}

	err = h.service.UpdateCategory(c.Request.Context(), id, &updateCat)
	if err != nil {
		h.logger.Error("Failed to update category", zap.Int("id", id), zap.Error(err))
		utils.ResponseError(c, http.StatusBadRequest, "Gagal memperbarui kategori: "+err.Error(), nil)
		return
	}

	utils.ResponseSuccess(c, http.StatusOK, "Kategori berhasil diperbarui", map[string]int{"id": id})
}

func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		utils.ResponseError(c, http.StatusBadRequest, "Format ID kategori tidak valid", nil)
		return
	}

	err = h.service.DeleteCategory(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("Failed to delete category", zap.Int("id", id), zap.Error(err))
		utils.ResponseError(c, http.StatusBadRequest, "Gagal menghapus kategori: "+err.Error(), nil)
		return
	}

	utils.ResponseSuccess(c, http.StatusOK, "Kategori berhasil dihapus", map[string]int{"id": id})
}

func (h *CategoryHandler) GetCategoryById(c *gin.Context) {
	idStr := c.Param("id")
	categoryID, err := strconv.Atoi(idStr)
	if err != nil || categoryID <= 0 {
		utils.ResponseBadRequest(c, http.StatusBadRequest, "Format ID kategori tidak valid", nil)
		return
	}

	category, err := h.service.GetCategoryById(c.Request.Context(), categoryID)
	if err != nil {
		h.logger.Error("Failed to get category by id", zap.Error(err), zap.Int("id", categoryID))
		utils.ResponseError(c, http.StatusNotFound, "Kategori tidak ditemukan", nil)
		return
	}

	response := dto.CategoryByIdResponse{
		ID:              category.ID,
		Name:            category.Name,
		RackInventory:   category.RackInventory,
		RackInventoryId: category.Rack_inventory_id,
	}
	utils.ResponseSuccess(c, http.StatusOK, "Berhasil mendapatkan data kategori", response)
}
