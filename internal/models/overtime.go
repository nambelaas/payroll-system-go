package models

import "time"

type Overtime struct {
	ID         uint      `gorm:"primaryKey"`
	EmployeeId uint      `gorm:"index;constraint:OnDelete:CASCADE;"`
	Date       time.Time `gorm:"index"`
	StartTime  time.Time
	EndTime    time.Time
	Hours      float64
	CreatedAt  time.Time
}
