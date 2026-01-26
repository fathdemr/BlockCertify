package models

import "time"

type User struct {
	ID          string    `json:"id" gorm:"primaryKey" binding:"required"`
	FirstName   string    `json:"firstName" binding:"required"`
	LastName    string    `json:"lastName" binding:"required"`
	Email       string    `json:"email" gorm:"uniqueIndex" binding:"required"`
	Password    string    `json:"password" binding:"required"`
	Institution string    `json:"institutionName" binding:"required"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
