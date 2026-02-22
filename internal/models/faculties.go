package models

import "github.com/gofrs/uuid/v5"

const TableFaculties = "faculties"

func (Faculties) TableName() string {
	return TableFaculties
}

type Faculties struct {
	ID          uuid.UUID `gorm:"primaryKey;column:id"`
	FacultyName string    `gorm:"column:name"`

	UniversityID uuid.UUID    `gorm:"column:university_id"`
	University   Universities `gorm:"foreignKey:UniversityID;references:ID"`

	Departments []Department `gorm:"foreignKey:FacultyID;references:ID"`

	BaseRecordFields
}
