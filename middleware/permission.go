package middleware

import (
	"net/http"
	"strconv"
	"strings" // <--- JANGAN LUPA IMPORT INI
)

func (middlewareCostume *MiddlewareCostume) RequirePermission(code string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// 1. CEK APAKAH COOKIE ADA
			cookie, err := r.Cookie("session")
			if err != nil {
				// Jika error (cookie tidak ada), return 401 Unauthorized
				// JANGAN lanjut ke bawah, langsung return
				http.Error(w, "Unauthorized: Silakan login terlebih dahulu", http.StatusUnauthorized)
				return
			}

			// 2. BERSIHKAN PREFIX "lumos-"
			// Cookie Anda formatnya: "lumos-12", kita butuh angka "12" saja.
			cleanValue := strings.TrimPrefix(cookie.Value, "lumos-")

			// 3. KONVERSI KE INT
			userID, err := strconv.Atoi(cleanValue)
			if err != nil {
				// Jika gagal convert (misal cookie dirusak user), return Bad Request
				http.Error(w, "Bad Request: Session invalid", http.StatusBadRequest)
				return
			}

			// 4. CEK PERMISSION KE SERVICE
			allowed, err := middlewareCostume.Service.Permission.Allowed(userID, code)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			if !allowed {
				http.Error(w, "Forbidden: Anda tidak punya akses", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
