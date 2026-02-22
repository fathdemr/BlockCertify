package models

import "github.com/gofrs/uuid/v5"

const TableDepartment = "departments"

func (Department) TableName() string {
	return TableDepartment
}

type Department struct {
	ID             uuid.UUID `gorm:"primaryKey;column:id"`
	DepartmentName string    `gorm:"column:name"`

	FacultyID uuid.UUID `gorm:"column:faculty_id"`
	Faculty   Faculties `gorm:"foreignKey:FacultyID;references:ID"`

	BaseRecordFields
}
