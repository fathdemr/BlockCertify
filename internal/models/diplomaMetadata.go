package models

import (
	"time"

	"github.com/google/uuid"
)

const TableName = "diploma_metadata"

func (DiplomaMetaData) TableName() string {
	return TableName
}

type DiplomaMetaData struct {
	ID             uuid.UUID `json:"id" gorm:"primary_key" binding:"required"`
	DiplomaID      uuid.UUID `json:"diplomaId" gorm:"uniqueIndex" binding:"required"`
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
