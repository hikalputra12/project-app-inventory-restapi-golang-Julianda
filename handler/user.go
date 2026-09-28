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

type UserHandler struct {
	service service.UserServiceInterface
	logger  *zap.Logger
}

func NewUserHandler(service service.UserServiceInterface, log *zap.Logger) UserHandler {
	return UserHandler{
		service: service,
		logger:  log,
	}
}

func (h *UserHandler) ListUser(c *gin.Context) {
	page := utils.StringToInt(c.DefaultQuery("page", "1"), 1)
	limit := utils.StringToInt(c.DefaultQuery("limit", "10"), 10)

	users, pagination, err := h.service.GetAllUser(c.Request.Context(), page, limit)
	if err != nil {
		h.logger.Error("Failed to get list user from service", zap.Error(err))
		utils.ResponseError(c, http.StatusInternalServerError, "Gagal mengambil daftar pengguna: "+err.Error(), nil)
		return
	}

	var response []dto.UserListResponse
	for _, item := range users {
		response = append(response, dto.UserListResponse{
			ID:    item.ID,
			Name:  item.Name,
			Email: item.Email,
			Role:  item.Role,
		})
	}
	if response == nil {
		response = []dto.UserListResponse{}
	}

	utils.ResponsePagination(c, http.StatusOK, "Berhasil mendapatkan daftar pengguna", response, *pagination)
}

func (h *UserHandler) UserById(c *gin.Context) {
	idStr := c.Param("id")
	userID, err := strconv.Atoi(idStr)
	if err != nil || userID <= 0 {
		utils.ResponseBadRequest(c, http.StatusBadRequest, "Format ID pengguna tidak valid", nil)
		return
	}

	user, err := h.service.GetUserById(c.Request.Context(), userID)
	if err != nil {
		h.logger.Error("Failed to get user by id", zap.Error(err), zap.Int("id", userID))
		utils.ResponseError(c, http.StatusNotFound, "Pengguna tidak ditemukan", nil)
		return
	}

	response := dto.UserByIdResponse{
		ID:     user.ID,
		Name:   user.Name,
		Email:  user.Email,
		Role:   user.Role,
		RoleID: user.Role_id,
	}
	utils.ResponseSuccess(c, http.StatusOK, "Berhasil mendapatkan data pengguna", response)
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req dto.CreateNewUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseBadRequest(c, http.StatusBadRequest, "Format JSON tidak valid", nil)
		return
	}

	if fieldErrors, err := utils.ValidateErrors(req); err != nil {
		utils.ResponseBadRequest(c, http.StatusBadRequest, "Validasi input gagal", fieldErrors)
		return
	}

	newUser := model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Role_id:  req.Role_id,
	}

	err := h.service.CreateUser(c.Request.Context(), &newUser)
	if err != nil {
		h.logger.Error("Failed to create user on service", zap.Error(err))
		utils.ResponseError(c, http.StatusBadRequest, "Gagal membuat pengguna baru: "+err.Error(), nil)
		return
	}

	utils.ResponseSuccess(c, http.StatusCreated, "Pengguna baru berhasil dibuat", map[string]interface{}{
		"id":    newUser.ID,
		"name":  newUser.Name,
		"email": newUser.Email,
	})
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		utils.ResponseError(c, http.StatusBadRequest, "Format ID pengguna tidak valid", nil)
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Failed to decode JSON body", zap.Error(err))
		utils.ResponseBadRequest(c, http.StatusBadRequest, "Format JSON tidak valid", nil)
		return
	}

	if fieldErrors, err := utils.ValidateErrors(req); err != nil {
		utils.ResponseBadRequest(c, http.StatusBadRequest, "Validasi input gagal", fieldErrors)
		return
	}

	updateUser := model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Role_id:  req.Role_id,
	}

	err = h.service.UpdateUser(c.Request.Context(), id, &updateUser)
	if err != nil {
		h.logger.Error("Failed to update user on service", zap.Int("id", id), zap.Error(err))
		utils.ResponseError(c, http.StatusBadRequest, "Gagal memperbarui pengguna: "+err.Error(), nil)
		return
	}

	utils.ResponseSuccess(c, http.StatusOK, "Pengguna berhasil diperbarui", map[string]int{"id": id})
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		utils.ResponseError(c, http.StatusBadRequest, "Format ID pengguna tidak valid", nil)
		return
	}

	err = h.service.DeleteUser(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("Failed to delete user on service", zap.Int("id", id), zap.Error(err))
		utils.ResponseError(c, http.StatusBadRequest, "Gagal menghapus pengguna: "+err.Error(), nil)
		return
	}

	utils.ResponseSuccess(c, http.StatusOK, "Pengguna berhasil dinonaktifkan/dihapus", map[string]int{"id": id})
}
