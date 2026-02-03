package handler

import (
	"app-inventory/dto"
	"app-inventory/model"
	"app-inventory/service"
	"app-inventory/utils"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

type AuthHandler struct {
	AuthService    service.AuthServiceInterface
	SessionService service.SessionServiceInterface
	Log            *zap.Logger
}

func NewAuthHandler(authHendler service.AuthServiceInterface, sessionHandler service.SessionServiceInterface, log *zap.Logger) AuthHandler {
	return AuthHandler{
		AuthService:    authHendler,
		SessionService: sessionHandler,
		Log:            log,
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

	//pembuatan uuid session token
	uuidToken := utils.NewUUID()
	session := &model.Session{
		SessionID: uuidToken,
		UserID:    user.ID,
	}
	err = h.SessionService.CreateSession(session)
	if err != nil {
		h.Log.Error("failed create session on service",
			zap.Error(err),
		)
		utils.ResponseError(w, http.StatusBadRequest, "failed to proccess login session", nil)
		return
	}
	expiryTime := 24 * 60 * 60
	// cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    uuidToken,
		Path:     "/",
		MaxAge:   expiryTime,
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
	cookie, err := r.Cookie("session")
	if err != nil {
		// Jika cookie tidak ada, anggap saja user sudah logout
		utils.ResponseError(w, http.StatusUnauthorized, "User tidak terautentikasi", nil)
		return
	}
	//pembuatan revoke saat logout
	sessionID := cookie.Value
	revoke := &model.Session{
		SessionID: sessionID,
	}

	err = h.SessionService.RevokeSession(revoke)
	if err != nil {
		h.Log.Error("failed revoke session on service",
			zap.Error(err),
		)
		utils.ResponseError(w, http.StatusBadRequest, "failed to proccess logout session", nil)
		return
	}
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
