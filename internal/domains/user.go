package domains

import "time"

type User struct {
	ID        int    `gorm:"primaryKey;autoIncrement"`
	Name      string `gorm:"type:varchar(100);not null"`
	Email     string `gorm:"type:varchar(100);unique;not null"`
	Password  string `gorm:"type:varchar(100);not null"`
	RoleID    string `gorm:"not null"`
	IsActive  int
	CreatedAt *time.Time `gorm:"column:created_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at"`
	UpdatedAt *time.Time
}

func (User) TableName() string {
	return "users"
}
