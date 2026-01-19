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
			//set uuid baru
			newSession := utils.NewUUID()
			//exted cookie
			NewsessionID := &model.Session{
				SessionID: newSession,
			}
			err = middlewareCostume.Service.SessionService.ExtendSession(NewsessionID)
			if err != nil {
				return
			}
			//lakukan update cookie
			http.SetCookie(w, &http.Cookie{
				Name:     "session",
				Value:    newSession,
				Path:     "/",
				MaxAge:   24 * 60 * 60,
				HttpOnly: true,
			})

			next.ServeHTTP(w, r)
		})
	}
}
