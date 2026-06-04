package models

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type Log struct {
	ID             uuid.UUID `json:"id,omitempty" gorm:"primaryKey"`
	ActionTime     time.Time `json:"action_time,omitempty" gorm:"index:IX_logs_responsestatusactiontime_V1,priority:2;index:IX_logs_responsestatusactiontime_V2,priority:2"`
	UserType       string    `json:"user_type,omitempty" gorm:"type:nvarchar(50)"`
	UserID         uint64    `json:"user_id,omitempty" gorm:"index:IX_logs_responsestatusactiontime_V1,priority:3"`
	IpAddress      string    `json:"ip_address,omitempty" gorm:"type:nvarchar(50)"`
	HTTPMethod     string    `json:"http_method,omitempty" gorm:"type:nvarchar(50);index:IX_actiontime_V1,priority:1"`
	URLPath        string    `json:"url_path,omitempty" gorm:"type:nvarchar(250);index:IX_logs_url_path,priority:1;index:IX_actiontime_V1,priority:2"`
	ResponseStatus int       `json:"response_status,omitempty" gorm:"index:IX_logs_responsestatusactiontime_V1,priority:1;index:IX_logs_responsestatusactiontime_V2,priority:1;;index:IX_actiontime_V1,priority:3"`
	RequestBody    string    `json:"request_body,omitempty" gorm:""`
	ResponseBody   string    `json:"response_body,omitempty" gorm:""`
	Latency        int64     `json:"latency,omitempty"`
}
