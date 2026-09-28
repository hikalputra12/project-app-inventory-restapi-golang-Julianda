package handler

import (
	"app-inventory/dto"
	"app-inventory/service"
	"app-inventory/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthHandler struct {
	AuthService    service.AuthServiceInterface
	SessionService service.SessionServiceInterface
	Log            *zap.Logger
}

func NewAuthHandler(authHandler service.AuthServiceInterface, sessionHandler service.SessionServiceInterface, log *zap.Logger) AuthHandler {
	return AuthHandler{
		AuthService:    authHandler,
		SessionService: sessionHandler,
		Log:            log,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseBadRequest(c, http.StatusBadRequest, "Format JSON tidak valid", nil)
		return
	}

	if fieldErrors, err := utils.ValidateErrors(req); err != nil {
		utils.ResponseBadRequest(c, http.StatusBadRequest, "Validasi gagal", fieldErrors)
		return
	}

	loginResp, err := h.AuthService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		utils.ResponseError(c, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	// Set HTTP-only Cookie with JWT token as fallback for web browsers
	c.SetCookie("session", loginResp.AccessToken, 24*3600, "/", "", false, true)

	utils.ResponseSuccess(c, http.StatusOK, "Login berhasil", loginResp)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	// Invalidate browser cookie
	c.SetCookie("session", "", -1, "/", "", false, true)

	utils.ResponseSuccess(c, http.StatusOK, "Logout berhasil", nil)
}
