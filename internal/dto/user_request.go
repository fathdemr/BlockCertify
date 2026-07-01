package dto

type RegisterRequest struct {
	FirstName string `json:"firstName" validate:"required,min=2"`
	LastName  string `json:"lastName" validate:"required,min=2"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type UpdateMeRequest struct {
	FirstName string `json:"firstName" validate:"omitempty,min=2"`
	LastName  string `json:"lastName" validate:"omitempty,min=2"`
	Password  string `json:"password" validate:"omitempty,min=8"`
}
