package models

import "github.com/gofrs/uuid/v5"

const StudentTableName = "student"

func (Student) TableName() string {
	return StudentTableName
}

type Student struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID       uuid.UUID `gorm:"type:uuid;not null"`
	UniversityID uuid.UUID `gorm:"type:uuid;not null"`
	FacultyID    uuid.UUID `gorm:"type:uuid;not null"`
	DepartmentID uuid.UUID `gorm:"type:uuid;not null"`

	University Universities `gorm:"foreignKey:UniversityID;"`
	Faculty    Faculties    `gorm:"foreignKey:FacultyID;"`
	Department Department   `gorm:"foreignKey:DepartmentID;"`
	User       User         `gorm:"foreignKey:UserID;"`
	BaseRecordFields
}
