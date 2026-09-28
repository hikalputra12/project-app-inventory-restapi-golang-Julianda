package utils

import (
	"app-inventory/model"
	"context"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const (
	UserClaimsKey contextKey = "user_claims"
)

type JWTClaims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	RoleID int    `json:"role_id"`
	jwt.RegisteredClaims
}

// GenerateJWT creates a new signed JWT token for a user
func GenerateJWT(user *model.User, secret string, expiry time.Duration) (string, error) {
	if secret == "" {
		secret = "default-secret-key-change-in-prod"
	}
	if expiry <= 0 {
		expiry = 24 * time.Hour
	}

	claims := JWTClaims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RoleID: user.Role_id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Subject:   user.Email,
			Issuer:    "inventory-api",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ValidateJWT validates the token string and returns parsed claims
func ValidateJWT(tokenStr, secret string) (*JWTClaims, error) {
	if secret == "" {
		secret = "default-secret-key-change-in-prod"
	}

	token, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid or expired token")
	}

	return claims, nil
}

// SetUserClaimsToContext stores user claims into standard context
func SetUserClaimsToContext(ctx context.Context, claims *JWTClaims) context.Context {
	return context.WithValue(ctx, UserClaimsKey, claims)
}

// GetClaimsFromContext retrieves user claims from standard context
func GetClaimsFromContext(ctx context.Context) (*JWTClaims, bool) {
	claims, ok := ctx.Value(UserClaimsKey).(*JWTClaims)
	return claims, ok
}

// GetUserIDFromContext retrieves user ID from standard context safely
func GetUserIDFromContext(ctx context.Context) (int, bool) {
	claims, ok := GetClaimsFromContext(ctx)
	if !ok || claims == nil {
		return 0, false
	}
	return claims.UserID, true
}

// GetClaimsFromGin retrieves user claims from Gin context
func GetClaimsFromGin(c *gin.Context) (*JWTClaims, bool) {
	val, exists := c.Get(string(UserClaimsKey))
	if exists {
		if claims, ok := val.(*JWTClaims); ok {
			return claims, true
		}
	}
	return GetClaimsFromContext(c.Request.Context())
}

// GetUserIDFromGin retrieves user ID from Gin context safely
func GetUserIDFromGin(c *gin.Context) (int, bool) {
	claims, ok := GetClaimsFromGin(c)
	if !ok || claims == nil {
		return 0, false
	}
	return claims.UserID, true
}
