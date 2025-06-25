package models

import "time"

type PayrollPeriod struct {
	ID        uint      `gorm:"primaryKey"`
	StartDate time.Time `gorm:"not null"`
	EndDate   time.Time `gorm:"not null"`
	IsClosed  bool      `gorm:"default:false"` // untuk lock setelah payroll dijalankan
	CreatedBy string
	UpdatedBy *string `gorm:"default:null"`
	CreatedAt time.Time
	UpdatedAt *time.Time `gorm:"default:null"`
}
