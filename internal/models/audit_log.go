package models

import "time"

type AuditLog struct {
	ID        uint   `gorm:"primaryKey"`
	RequestId string `gorm:"unique"`
	UserId    uint   `gorm:"index;constraint:OnDelete:CASCADE;"`
	Endpoint  string
	Status    int
	Result    string `gorm:"type:text"`
	IpAddress string
	CreatedAt time.Time
	UpdatedAt *time.Time `gorm:"default:null"`
}
