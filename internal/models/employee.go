package models

type Employee struct {
	ID     uint `gorm:"primaryKey"`
	UserId uint `gorm:"uniqueIndex;constraint:OnDelete:CASCADE;"`
	Name   string
	Salary float64
}
