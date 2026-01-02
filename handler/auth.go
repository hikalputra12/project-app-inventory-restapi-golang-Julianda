package handler

import (
	"app-inventory/service"
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
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	user, err := h.AuthService.Login(email, password)
	if err != nil {
		return
	}

	// cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    "lumos-" + strconv.Itoa(user.ID),
		Path:     "/",
		HttpOnly: true,
	})

	http.Redirect(w, r, "/user/home", http.StatusSeeOther)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
