package dto

// LoginRequest defines credentials for authentication
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// LoginResponse defines the JWT response structure
type LoginResponse struct {
	AccessToken string      `json:"access_token"`
	TokenType   string      `json:"token_type"`
	User        UserProfile `json:"user"`
}

// UserProfile defines basic user info returned on authentication
type UserProfile struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	RoleID int    `json:"role_id"`
}
