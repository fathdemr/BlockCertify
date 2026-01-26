package dto

type RegisterRequest struct {
	FirstName   string `json:"firstName" validate:"required,min=2"`
	LastName    string `json:"lastName" validate:"required,min=2"`
	Email       string `json:"email" validate:"required,email"`
	Institution string `json:"institution" validate:"required,min=2"`
	Password    string `json:"password" validate:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}
