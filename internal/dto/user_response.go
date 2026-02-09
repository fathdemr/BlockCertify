package dto

type LoginResponse struct {
	Token     string `json:"token"` // Aliasing for frontend
	TokenType string `json:"tokenType"`
	ExpiresIn int64  `json:"expiresIn"`
	Role      string `json:"role"`
}
