package models

import "time"

type BaseRecordFields struct {
	CreatedAt        time.Time `json:"created_at,omitempty"`         // Oluşturulma Tarihi
	CreatedBy        uint64    `json:"created_by,omitempty"`         // Oluşturan Kullanıcı ID'si
	CreatedByName    string    `json:"created_by_name,omitempty"`    // Oluşturan Kullanıcı Adı
	CreatedIpAddress string    `json:"created_ip_address,omitempty"` // Oluşturan kullanıcının ip adresi
	UpdatedAt        time.Time `json:"updated_at,omitempty"`         // Güncellenme Tarihi
	UpdatedBy        uint64    `json:"updated_by,omitempty"`         // Güncelleyen Kullanıcı ID'si
	UpdatedByName    string    `json:"updated_by_name,omitempty"`    // Güncelleyen Kullanıcı Adı
	UpdatedIpAddress string    `json:"updated_ip_address,omitempty"` // Güncelleyen kullanıcının ip adresi
}

func (u *BaseRecordFields) SetCreatedUser(id uint64, name string, ipAddress string) {
	u.CreatedAt = time.Now()
	u.CreatedBy = id
	u.CreatedByName = name
	u.CreatedIpAddress = ipAddress
}

func (u *BaseRecordFields) SetUpdatedUser(id uint64, name string, ipAddress string) {
	u.UpdatedAt = time.Now()
	u.UpdatedBy = id
	u.UpdatedByName = name
	u.UpdatedIpAddress = ipAddress
}
