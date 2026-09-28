package middleware

import (
	"app-inventory/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTAuth middleware verifies JWT Bearer token or session cookie in Gin
func (m *CustomMiddleware) JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := ""

		// 1. Check Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
				tokenString = parts[1]
			}
		}

		// 2. Fallback to Cookie
		if tokenString == "" {
			if cookie, err := c.Cookie("session"); err == nil && cookie != "" {
				tokenString = cookie
			}
		}

		if tokenString == "" {
			utils.ResponseError(c, http.StatusUnauthorized, "Token autentikasi tidak ditemukan. Harap login terlebih dahulu.", nil)
			c.Abort()
			return
		}

		// 3. Validate JWT
		claims, err := utils.ValidateJWT(tokenString, m.JWTConfig.Secret)
		if err != nil {
			utils.ResponseError(c, http.StatusUnauthorized, "Sesi tidak valid atau telah kadaluarsa. Silakan login kembali.", nil)
			c.Abort()
			return
		}

		// 4. Store claims in Gin context and request context
		c.Set(string(utils.UserClaimsKey), claims)
		ctx := utils.SetUserClaimsToContext(c.Request.Context(), claims)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

// ValidAndExtendSession alias for backwards compatibility
func (m *CustomMiddleware) ValidAndExtendSession() gin.HandlerFunc {
	return m.JWTAuth()
}
