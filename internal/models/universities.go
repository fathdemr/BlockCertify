package models

import "github.com/gofrs/uuid/v5"

const tableName = "universities"

func (Universities) TableName() string {
	return tableName
}

type Universities struct {
	ID      uuid.UUID `gorm:"primaryKey;column:id"`
	Name    string    `gorm:"column:name"`
	YokCode string    `gorm:"column:yok_code"`

	Faculties []Faculties `gorm:"foreignKey:UniversityID;references:ID"`
	BaseRecordFields
}
