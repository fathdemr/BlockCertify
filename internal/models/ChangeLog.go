package models

import (
	"time"

	"gorm.io/gorm"
)

type ChangeLog struct {
	ID            uint64    `json:"id" gorm:"primaryKey"`
	UpdatedBy     uint64    `json:"updated_by" gorm:"type:nvarchar(255)"`
	UpdatedType   string    `json:"updated_type" gorm:"type:nvarchar(50)"`
	UpdatedByName string    `json:"updated_by_name" gorm:"type:nvarchar(250)"`
	TableName     string    `json:"table_name" gorm:"type:nvarchar(50)"`
	RecordID      uint64    `json:"record_id"`
	Type          string    `json:"type" gorm:"type:nvarchar(50)"`
	FieldName     string    `json:"field_name" gorm:"type:nvarchar(255)"`
	OldValue      string    `json:"old_value"`
	NewValue      string    `json:"new_value"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type ChangeLogType string

const (
	Create ChangeLogType = "create"
	Update ChangeLogType = "update"
	Delete ChangeLogType = "delete"
)

func (s ChangeLogType) String() string {
	switch s {
	case Create:
		return "create"
	case Update:
		return "update"
	case Delete:
		return "delete"
	}
	return "unknown"
}

func (c *ChangeLog) BeforeCreate(tx *gorm.DB) {
	c.UpdatedAt = time.Now()
}

func (c *ChangeLog) BeforeSave(tx *gorm.DB) {
	c.UpdatedAt = time.Now()
}
