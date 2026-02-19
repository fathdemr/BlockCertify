package models

import "github.com/gofrs/uuid/v5"

const AdminTableName = "admin"

func (Admin) TableName() string {
	return AdminTableName
}

type Admin struct {
	ID           uuid.UUID
	UserID       uuid.UUID
	UniversityID uuid.UUID

	University Universities `gorm:"foreignKey:UniversityID;"`
	User       User         `gorm:"foreignKey:UserID;"`
	BaseRecordFields
}
