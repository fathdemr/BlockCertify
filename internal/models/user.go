package models

import "time"

const TableUser = "user"

func (User) TableName() string {
	return TableUser
}

type User struct {
	ID          string `json:"id" gorm:"primaryKey"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Email       string `gorm:"uniqueIndex;not null"`
	Password    string
	Institution string
	Role        string `json:"role"`
	Wallet
	CreatedAt time.Time
	UpdatedAt time.Time
}
