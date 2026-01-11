package dto

type UserListResponse struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type CreateNewUserRequest struct {
	Name     string `json:"name" validate:"required, min=3"`
	Email    string `json:"email" validate:"required,email"`
	Role_id  int    `json:"role_id" validate:"required,gte=1,lte=3"`
	Password string `json:"password" validate:"required,min=6"`
}

type UpdateUserRequest struct {
	Name     string `json:"name" validate:"required, min=3"`
	Email    string `json:"email" validate:"required,email"`
	Role_id  int    `json:"role_id" validate:"required,gte=1,lte=3"`
	Password string `json:"password" validate:"required,min=6"`
}
type UserByIdResponse struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	RoleID   int    `json:"-"`
	Password string `json:"-"`
}
