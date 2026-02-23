package dto

import "github.com/gofrs/uuid/v5"

type DepartmentRequest struct {
	FacultyID uuid.UUID `json:"id"`
}
