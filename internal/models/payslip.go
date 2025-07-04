package models

import "time"

type Payslip struct {
	ID               uint `gorm:"primaryKey"`
	EmployeeId       uint `gorm:"index;constraint:OnDelete:CASCADE;"`
	PayrollPeriodId  uint `gorm:"index"`
	AttendanceHours  float64
	ProratedSalary   float64
	OvertimeHours    float64
	OvertimePay      float64
	ReimbursementSum float64
	TotalTakeHome    float64
	CreatedBy        string
	UpdatedBy        *string `gorm:"default:null"`
	CreatedAt        time.Time
	UpdatedAt        *time.Time `gorm:"default:null"`
}
