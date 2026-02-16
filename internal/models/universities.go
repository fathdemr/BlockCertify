package models

import "github.com/gofrs/uuid/v5"

const tableName = "universities"

func (Universities) TableName() string {
	return tableName
}

type Universities struct {
	ID      uuid.UUID
	Name    string
	YokCode string
	BaseRecordFields
}
