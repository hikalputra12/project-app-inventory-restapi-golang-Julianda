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
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateNewUserRequest
	//mengubah json body ke struct
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "Invalid JSON format", nil)
		return
	}
	newUser := model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Role_id:  req.Role_id,
	}
	err := h.service.CreateUser(&newUser)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "input tidak sesuai format yang di tentukan", nil)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  true,
		"message": "Create new user succesfully",
	})

}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateUserRequest
	//mengubah json body ke struct
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "Invalid JSON format", nil)
		return
	}

	newUser := model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Role_id:  req.Role_id,
	}
	err := h.service.UpdateUser(req.User_id, &newUser)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "input tidak sesuai format yang di tentukan", nil)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  true,
		"message": "Update user succesfully",
	})

}
