package models

import "github.com/gofrs/uuid/v5"

const TableDepartment = "departments"

func (Department) TableName() string {
	return TableDepartment
}

type Department struct {
	ID             uuid.UUID
	FacultyID      uuid.UUID
	DepartmentName string
}
