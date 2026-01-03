package dto

type LoginRequest struct {
	// Perhatikan tag 'validate' di belakang
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}
