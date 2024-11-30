package dto

type CreateUserRequest struct {
	Name     string `json:"name" validate:"required,min=3,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=100"`
}

type UpdateUserRequest struct {
	Name     string  `json:"name,omitempty" validate:"omitempty,min=3,max=100"`
	Email    string  `json:"email,omitempty" validate:"omitempty,email"`
	Password string  `json:"password,omitempty" validate:"omitempty,min=6,max=100"`
	Phone    *string `json:"phone,omitempty" validate:"omitempty,min=3,max=100"`
	Role     string  `json:"role,omitempty" validate:"omitempty,oneof=admin customer staff"`
}
