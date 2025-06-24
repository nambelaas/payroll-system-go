package pkg

import (
	"log"

	"github.com/nambelaas/payroll-system-go/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	cred := "host=localhost user=postgres password=postgres dbname=payrolldb port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(cred), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// auto migrate
	db.AutoMigrate(
		&models.User{},
		&models.Employee{},
		&models.PayrollPeriod{},
		&models.Attendance{},
		&models.Overtime{},
		&models.Reimbursement{},
		&models.Payslip{},
	)

	DB = db
}

func ConnectTestDB() *gorm.DB {
	cred := "host=localhost user=postgres password=postgres dbname=payrolldb_test port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(cred), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to test DB")
	}

	db.AutoMigrate(&models.User{}, &models.Employee{}, &models.Attendance{}, &models.PayrollPeriod{}, &models.Reimbursement{}, &models.Overtime{}, &models.Payslip{})
	return db
}
