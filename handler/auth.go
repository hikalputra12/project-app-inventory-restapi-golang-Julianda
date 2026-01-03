package handler

import (
	"app-inventory/dto"
	"app-inventory/service"
	"app-inventory/utils"
	"encoding/json"
	"net/http"
	"strconv"

	"go.uber.org/zap"
)

type AuthHandler struct {
	AuthService service.AuthServiceInterface
	Log         *zap.Logger
}

func NewAuthHandler(authHendler service.AuthServiceInterface, log *zap.Logger) AuthHandler {
	return AuthHandler{
		AuthService: authHendler,
		Log:         log,
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {

	var req dto.LoginRequest

	//mengubah json body ke struct
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "Invalid JSON format", nil)
		return
	}

	validationErrors, err := utils.ValidateErrors(req)
	if err != nil {

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  false,
			"message": "Validation failed",
			"errors":  validationErrors,
		})
		return
	}
	user, err := h.AuthService.Login(req.Email, req.Password)
	if err != nil {
		utils.ResponseError(w, http.StatusUnauthorized, "Email atau password salah", nil)
		return
	}
	// cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    strconv.Itoa(user.ID),
		Path:     "/",
		HttpOnly: true,
	})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  true,
		"message": "Login successful",
	})

}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Clear Cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	// Return JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  true,
		"message": "Logout successful",
	})
}
