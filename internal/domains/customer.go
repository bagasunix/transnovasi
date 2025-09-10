package domains

import "time"

type Customer struct {
	ID       int    `gorm:"primaryKey;autoIncrement"`
	Name     string `gorm:"type:varchar(100);not null"`
	Email    string `gorm:"type:varchar(100);unique;not null"`
	Phone    string `gorm:"type:varchar(14);not null"`
	Address  string
	IsActive int `gorm:"column:is_active"`

	CreatedBy int        `gorm:"column:created_by"`
	CreatedAt *time.Time `gorm:"column:created_at;autoautoCreateTime"`
	UpdatedAt *time.Time `gorm:"column:created_at;autoUpdateTime"`
	DeletedAt *time.Time `gorm:"column:deleted_at"`

	Vehicle []Vehicle `gorm:"foreignKey:CustomerID;references:ID"`
}

func (Customer) TableName() string {
	return "customers"
}
