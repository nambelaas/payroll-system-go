package models

import "time"

type Reimbursement struct {
	ID          uint    `gorm:"primaryKey"`
	EmployeeId  uint    `gorm:"index;constraint:OnDelete:CASCADE;"`
	Description string  `gorm:"type:text"`
	Amount      float64 `gorm:"not null"`
	Status      string  `gorm:"default:'pending'"`
	Date        time.Time
	CreatedBy   string
	UpdatedBy   *string `gorm:"default:null"`
	CreatedAt   time.Time
	UpdatedAt   *time.Time `gorm:"default:null"`
}
