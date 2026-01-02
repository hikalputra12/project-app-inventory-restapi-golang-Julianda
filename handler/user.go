package handler

import (
	"app-inventory/service"
	"app-inventory/utils"
	"net/http"
	"strconv"

	"go.uber.org/zap"
)

type UserHandler struct {
	service service.UserServiceInterface
	logger  *zap.Logger
}

// constructor
func NewUserHandler(service service.UserServiceInterface, log *zap.Logger) UserHandler {
	return UserHandler{
		service: service,
		logger:  log,
	}
}

func (h *UserHandler) ListUser(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "Invalid page", nil)
		return
	}

	// config limit pagination
	limit := 3

	// Get data users form service all users
	users, pagination, err := h.service.GetAllUser(page, limit)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusInternalServerError, "Failed to fetch assignments: "+err.Error(), nil)
		return
	}

	utils.ResponsePagination(w, http.StatusOK, "success get data", users, *pagination)

}
