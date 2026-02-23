package dto

import "github.com/gofrs/uuid/v5"

type FacultiesRequest struct {
	UniversityID uuid.UUID `json:"id"`
}
