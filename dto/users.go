package dto

type AssignmentRequest struct {
	Name  string `json:"name" validate:"required,min=3,max=150"`
	Email string `json:"email" validate:"required,min=10"`
	Role  string `json:"role" validate:"required"`
}

type AssignmentResponse struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}
