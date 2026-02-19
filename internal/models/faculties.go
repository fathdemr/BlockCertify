package models

import "github.com/gofrs/uuid/v5"

const TableFaculties = "faculties"

func (Faculties) TableName() string {
	return TableFaculties
}

type Faculties struct {
	ID           uuid.UUID
	UniversityID uuid.UUID
	FacultyName  string

	BaseRecordFields
}
