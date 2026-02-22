package models

import "time"

type BaseRecordFields struct {
	CreatedAt        time.Time `gorm:"column:created_at"`      // Oluşturulma Tarihi
	CreatedBy        uint64    `gorm:"column:created_by"`      // Oluşturan Kullanıcı ID'si
	CreatedByName    string    `gorm:"column:created_by_name"` // Oluşturan Kullanıcı Adı
	CreatedIpAddress string    `gorm:"column:created_ip"`      // Oluşturan kullanıcının ip adresi
	UpdatedAt        time.Time `gorm:"column:updated_at"`      // Güncellenme Tarihi
	UpdatedBy        uint64    `gorm:"column:updated_by"`      // Güncelleyen Kullanıcı ID'si
	UpdatedByName    string    `gorm:"column:updated_by_name"` // Güncelleyen Kullanıcı Adı
	UpdatedIp        string    `gorm:"column:updated_ip"`      // Güncelleyen Kullanıcı IP'si
}
