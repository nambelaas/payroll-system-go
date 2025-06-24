package models

import "time"

type PayrollPeriod struct {
	ID        uint      `gorm:"primaryKey"`
	StartDate time.Time `gorm:"not null"`
	EndDate   time.Time `gorm:"not null"`
	CreatedAt time.Time
	IsClosed  bool `gorm:"default:false"` // untuk lock setelah payroll dijalankan
}
