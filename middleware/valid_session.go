package middleware

import (
	"app-inventory/model"
	"app-inventory/utils"
	"net/http"
)

func (middlewareCostume *MiddlewareCostume) ValidAndExtendSession() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, _ := r.Cookie("session")
			getSessionID := session.Value
			sessionID := &model.Session{
				SessionID: getSessionID,
			}
			valid, err := middlewareCostume.Service.SessionService.IsValid(sessionID)
			if err != nil {
				return
			}

			if !valid {
				utils.ResponseError(w, http.StatusUnauthorized, "Sesi sudah kadaluarsa atau tidak valid", nil)
				return
			}
			err = middlewareCostume.Service.SessionService.ExtendSession(sessionID)
			if err != nil {
				utils.ResponseError(w, http.StatusInternalServerError, "Terjadi kesalahan sistem saat memproses sesi", nil)
				return
			}
			//lakukan update cookie
			http.SetCookie(w, &http.Cookie{
				Name:     "session_token",
				Value:    getSessionID,
				Path:     "/",
				MaxAge:   24 * 60 * 60,
				HttpOnly: true,
			})

			next.ServeHTTP(w, r)
		})
	}
}
