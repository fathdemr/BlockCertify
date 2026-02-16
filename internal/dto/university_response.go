package dto

import "github.com/gofrs/uuid/v5"

type UniversitiesResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
