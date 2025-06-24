package seeder

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/nambelaas/payroll-system-go/internal/models"
	"github.com/nambelaas/payroll-system-go/pkg"

	"github.com/brianvoe/gofakeit/v6"
	"golang.org/x/crypto/bcrypt"
)

func SeedUser() {
	db := pkg.DB
	gofakeit.Seed(time.Now().UnixNano())

	// admin
	adminPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	admin := models.User{
		Username: "admin",
		Password: string(adminPassword),
		Role:     "admin",
	}
	db.FirstOrCreate(&admin, models.User{Username: "admin"})
	log.Println("Admin Created: ", admin.Username, "/", admin.Password)

	// employee
	for i := 0; i < 100; i++ {
		name := gofakeit.Name()
		username := fmt.Sprintf("employee%d", i+1)
		password := fmt.Sprintf("password%d", i+1)
		hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		user := models.User{
			Username: username,
			Password: string(hashed),
			Role:     "employee",
		}
		db.FirstOrCreate(&user, models.User{Username: username})

		salary := float64(rand.Intn(5)+5) * 1000000 // gaji random antara 5-10 juta
		employee := models.Employee{
			UserId: user.ID,
			Name:   name,
			Salary: salary,
		}
		db.FirstOrCreate(&employee, models.Employee{UserId: user.ID})
		log.Printf("Employee Created: %s / %s (%.0f)", username, password, salary)
	}
}
