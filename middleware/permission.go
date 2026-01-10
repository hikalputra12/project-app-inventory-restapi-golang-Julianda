package middleware

import (
	"app-inventory/model"
	"net/http"
)

func (middlewareCostume *MiddlewareCostume) RequirePermission(code string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, _ := r.Cookie("session")
			getSessionID := session.Value
			sessionID := &model.Session{
				SessionID: getSessionID,
			}
			userID, err := middlewareCostume.Service.SessionService.GetUserIDBySession(sessionID)

			allowed, err := middlewareCostume.Service.Permission.Allowed(userID, code)
			if err != nil {
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}

			if !allowed {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
