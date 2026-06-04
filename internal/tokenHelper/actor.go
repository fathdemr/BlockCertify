package tokenHelper

import "github.com/gofrs/uuid/v5"

type Actor struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Role string    `json:"role"`
}
