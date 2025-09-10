package domains

import "time"

type AuditLog struct {
	ID         int       `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID     string    `json:"user_id" gorm:"column:user_id"`
	CustomerID string    `json:"customer_id" gorm:"column:customer_id"`
	Method     string    `json:"method" gorm:"type:varchar(50);not null"`
	Endpoint   string    `json:"url" gorm:"column:url;not null"`
	StatusCode string    `json:"status_code" gorm:"column:status_code"`
	Request    string    `json:"request" gorm:"column:request"`
	Response   string    `json:"Response" gorm:"column:response"`
	IPAddress  string    `json:"ip_address" gorm:"column:ip_address;type:varchar(45);not null"`
	UserAgent  string    `json:"user_agent" gorm:"column:user_agent;type:text;not null"`
	CreatedAt  time.Time `json:"created_at" gorm:"column:created_at;type:timestamp;default:current_timestamp"`
}

// TableName sets the insert table name for this struct type
func (AuditLog) TableName() string {
	return "audit_logs"
}
