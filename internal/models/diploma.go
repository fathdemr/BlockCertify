package models

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

const TableDiploma = "diploma"

func (Diploma) TableName() string {
	return TableDiploma
}

type Diploma struct {
	ID          uuid.UUID `gorm:"primary_key;type:uuid"`
	PublicID    string    `gorm:"uniqueIndex;not null"`
	Hash        string
	ArweaveTxID string
	ArweaveURL  string
	PolygonTxID string
	PolygonURL  string
	Owner       string
	Timestamp   time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time

	MetaData DiplomaMetaData `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
