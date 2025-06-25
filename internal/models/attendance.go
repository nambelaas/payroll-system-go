package models

import "time"

type Attendance struct {
	ID         uint      `gorm:"primaryKey"`
	EmployeeId uint      `gorm:"index;constraint:OnDelete:CASCADE;"`
	Date       time.Time `gorm:"index"`
	CheckIn    time.Time
	CheckOut   *time.Time `gorm:"default:null"`
	CreatedAt  time.Time
	UpdatedAt  *time.Time `gorm:"default:null"`
	CreatedBy  string
	UpdatedBy  *string `gorm:"default:null"`
}
