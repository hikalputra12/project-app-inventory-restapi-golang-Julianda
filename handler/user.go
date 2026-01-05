package handler

import (
	"app-inventory/dto"
	"app-inventory/model"
	"app-inventory/service"
	"app-inventory/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
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

	var response []dto.UserListResponse
	for _, item := range users {
		response = append(response, dto.UserListResponse{
			Name:  item.Name,
			Email: item.Email,
			Role:  item.Role,
		})

	}
	utils.ResponsePagination(w, http.StatusOK, "success get data", response, *pagination)

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
	// --- DEBUGGING START ---
	fmt.Println("\n--- DEBUG INFO ---")
	fmt.Printf("1. URL Asli yg masuk: %s\n", r.URL.Path)

	// Cek apa yang ditangkap chi
	rctx := chi.RouteContext(r.Context())
	if rctx != nil {
		fmt.Printf("2. Chi Pattern: %s\n", rctx.RoutePattern())
		fmt.Printf("3. Chi URL Params Keys: %v\n", rctx.URLParams.Keys)
		fmt.Printf("4. Chi URL Params Values: %v\n", rctx.URLParams.Values)
	} else {
		fmt.Println("2. Chi Context NIL (Bahaya!)")
	}
	// --- DEBUGGING END ---
	//mengambil id
	idStr := chi.URLParam(r, "id")
	// DEBUG: Cek apa isi idStr
	fmt.Printf("DEBUG ID: '%s'\n", idStr)
	id, err := strconv.Atoi(idStr)

	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Invalid ID format (harus angka)", nil)
		return
	}
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
	err = h.service.UpdateUser(id, &newUser)
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

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	//mengambil id
	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)

	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Invalid ID format (harus angka)", nil)
		return
	}

	err = h.service.DeleteUser(id)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "input tidak sesuai format yang di tentukan", nil)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  true,
		"message": "Delete user succesfully ",
	})

}
