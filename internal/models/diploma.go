package models

import (
	"time"

	"github.com/google/uuid"
)

const TableDiploma = "diploma"

func (Diploma) TableName() string {
	return TableDiploma
}

type Diploma struct {
	ID          uuid.UUID `gorm:"primary_key"`
	PublicID    string    `gorm:"uniqueIndex;not null"`
	Hash        string
	ArweaveTxID string
	Owner       string
	Timestamp   time.Time

	MetaData DiplomaMetaData `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
