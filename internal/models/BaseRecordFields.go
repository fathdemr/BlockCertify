package models

import "time"

type BaseRecordFields struct {
	CreatedAt        time.Time `gorm:"column:created_at"`         // Oluşturulma Tarihi
	CreatedBy        uint64    `gorm:"column:created_by"`         // Oluşturan Kullanıcı ID'si
	CreatedByName    string    `gorm:"column:created_by_name"`    // Oluşturan Kullanıcı Adı
	CreatedIpAddress string    `gorm:"column:created_ip"`         // Oluşturan kullanıcının ip adresi
	UpdatedAt        time.Time `gorm:"column:updated_at"`         // Güncellenme Tarihi
	UpdatedBy        uint64    `gorm:"column:updated_by"`         // Güncelleyen Kullanıcı ID'si
	UpdatedByName    string    `gorm:"column:updated_by_name"`    // Güncelleyen Kullanıcı Adı
	UpdatedIp        string    `gorm:"updated_ip"`                // Güncelleyen kullanıcının ip
	UpdatedIPAddress string    `gorm:"column:updated_ip_address"` // Güncelleyen kullanıcının ip adresi
}

func (u *BaseRecordFields) SetCreatedUser(id uint64, name string, ipAddress string) {
	u.CreatedAt = time.Now()
	u.CreatedBy = id
	u.CreatedByName = name
	u.CreatedIpAddress = ipAddress
}

func (u *BaseRecordFields) SetUpdatedUser(id uint64, name string, ip string, ipAddress string) {
	u.UpdatedAt = time.Now()
	u.UpdatedBy = id
	u.UpdatedByName = name
	u.UpdatedIp = ip
	u.UpdatedIPAddress = ipAddress
}
