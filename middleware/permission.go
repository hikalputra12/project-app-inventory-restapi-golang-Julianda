package middleware

import (
	"app-inventory/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RequirePermission verifies role permission in Gin
func (m *CustomMiddleware) RequirePermission(code string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, ok := utils.GetClaimsFromGin(c)
		if !ok || claims == nil || claims.UserID <= 0 {
			utils.ResponseError(c, http.StatusUnauthorized, "Pengguna tidak terautentikasi", nil)
			c.Abort()
			return
		}

		// Super Admin bypass
		if claims.Role == "super_admin" || claims.Role == "Super Admin" || claims.RoleID == 1 {
			c.Next()
			return
		}

		allowed, err := m.Service.Permission.Allowed(c.Request.Context(), claims.UserID, code)
		if err != nil {
			m.Log.Error("Failed to check user permission",
				zap.Int("user_id", claims.UserID),
				zap.String("permission_code", code),
				zap.Error(err),
			)
			utils.ResponseError(c, http.StatusInternalServerError, "Gagal memeriksa hak akses pengguna", nil)
			c.Abort()
			return
		}

		if !allowed {
			utils.ResponseError(c, http.StatusForbidden, "Akses ditolak: Anda tidak memiliki izin untuk tindakan ini ("+code+")", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}
