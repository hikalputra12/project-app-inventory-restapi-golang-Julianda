package middleware

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

// Logging adalah middleware untuk mencatat setiap request HTTP
func Logging(log *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Jalankan request
			next.ServeHTTP(w, r)

			// Catat log setelah request selesai
			duration := time.Since(start)
			log.Info("HTTP Request",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Duration("duration", duration),
			)
		})
	}
}
