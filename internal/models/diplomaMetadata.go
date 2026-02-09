package models

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

const TableName = "diploma_metadata"

func (DiplomaMetaData) TableName() string {
	return TableName
}

type DiplomaMetaData struct {
	ID             uuid.UUID `json:"id" gorm:"primary_key;type:uuid" binding:"required"`
	DiplomaID      uuid.UUID `json:"diplomaId" gorm:"uniqueIndex;type:uuid" binding:"required"`
	FirstName      string
	LastName       string
	Email          string
	University     string
	Faculty        string
	Department     string
	GraduationYear int
	StudentNumber  string
	Nationality    string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
