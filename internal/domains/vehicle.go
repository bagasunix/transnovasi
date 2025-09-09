package domains

import "time"

type Vehicle struct {
	ID         int    `gorm:"primaryKey;autoIncrement"`
	CustomerID int    `gorm:"not null"`
	PlateNo    string `gorm:"not null"` // nomor polisi
	Model      string // model/tipe kendaraan
	Brand      string // merk
	Color      string
	Year       int
	IsActive   int

	CreatedBy int        `gorm:"column:created_by"`
	CreatedAt *time.Time `gorm:"column:created_at;autoautoCreateTime"`
	UpdatedAt *time.Time `gorm:"column:created_at;autoUpdateTime"`
	DeletedAt *time.Time `gorm:"column:deleted_at"`
}

func (Vehicle) TableName() string {
	return "vehicles"
}
