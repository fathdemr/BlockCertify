package models

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

const TableUser = "user"

func (User) TableName() string {
	return TableUser
}

type User struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `gorm:"uniqueIndex;not null"`
	Password  string
	Role      UserRole
	Wallet
	CreatedAt time.Time
	UpdatedAt time.Time
}
