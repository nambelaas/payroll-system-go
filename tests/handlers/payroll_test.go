package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nambelaas/payroll-system-go/internal/handlers"
	"github.com/nambelaas/payroll-system-go/internal/models"
	"github.com/nambelaas/payroll-system-go/pkg"
	"github.com/stretchr/testify/assert"
)

func setupPayrollRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(2))
		c.Next()
	})

	r.POST("/admin/payroll-period", handlers.CreatePayrollPeriod)
	r.POST("/admin/payroll/run", handlers.RunPayroll)
	return r
}

func TestCreatePayrollPeriodSuccess(t *testing.T) {
	db := pkg.ConnectTestDB()
	pkg.DB = db

	hashedPassword, _ := pkg.HashPassword("passwordTest123")
	db.FirstOrCreate(&models.User{
		ID:       2,
		Username: "testingAdmin",
		Password: hashedPassword,
		Role:     "admin",
	}, models.User{
		ID: 2,
	})

	r := setupPayrollRouter()

	body := map[string]string{
		"start_date": "2025-05-26",
		"end_date":   "2025-06-26",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPost, "/admin/payroll-period", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)
	assert.Contains(t, resp.Body.String(), "Payroll period created")
}

func TestCreatePayrollPeriodInvalidDate(t *testing.T) {
	db := pkg.ConnectTestDB()
	pkg.DB = db

	hashedPassword, _ := pkg.HashPassword("passwordTest123")
	db.FirstOrCreate(&models.User{
		ID:       2,
		Username: "testingAdmin",
		Password: hashedPassword,
		Role:     "admin",
	}, models.User{
		ID: 2,
	})

	r := setupPayrollRouter()

	body := map[string]string{
		"start_date": "2025-07-26",
		"end_date":   "2025-06-26",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPost, "/admin/payroll-period", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, resp.Body.String(), "Invalid date range")
}

func TestRunPayrollSuccess(t *testing.T) {
	db := pkg.ConnectTestDB()
	pkg.DB = db

	db.AutoMigrate(&models.Employee{}, &models.PayrollPeriod{}, &models.Attendance{}, &models.Overtime{}, &models.Reimbursement{}, &models.Payslip{})

	period := models.PayrollPeriod{
		ID:        1,
		StartDate: time.Now().AddDate(0, 0, -5),
		EndDate:   time.Now(),
	}
	db.FirstOrCreate(&period, models.PayrollPeriod{
		ID: 1,
	})

	emp := models.Employee{
		ID:     1,
		UserId: 1,
		Salary: 5000000,
	}
	db.FirstOrCreate(&emp, models.Employee{
		ID: 1,
	})

	for i := 0; i < 3; i++ {
		date := time.Now().AddDate(0, 0, -i)
		db.FirstOrCreate(&models.Attendance{
			EmployeeId: emp.ID,
			Date:       date,
			CheckIn:    time.Date(0, 0, 0, 9, 0, 0, 0, time.Local),
			CheckOut:   ptrTime(time.Date(0, 0, 0, 17, 0, 0, 0, time.Local)),
		}, models.Attendance{
			EmployeeId: emp.ID,
			Date:       date,
		})
	}

	db.FirstOrCreate(&models.Overtime{
		EmployeeId: emp.ID,
		Date:       time.Now(),
		Hours:      2,
	}, models.Overtime{
		EmployeeId: emp.ID,
		Date:       time.Now(),
	})

	db.FirstOrCreate(&models.Reimbursement{
		EmployeeId:  emp.ID,
		Amount:      150000,
		Description: "Transport",
		Date:        time.Now(),
		Status:      "pending",
	}, models.Reimbursement{
		EmployeeId: emp.ID,
		Date:       time.Now(),
	})

	db.Model(&models.PayrollPeriod{}).Where("id = ?", 1).Update("is_closed", false)

	body := map[string]interface{}{"payroll_period_id": 1}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPost, "/admin/payroll/run", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router := setupPayrollRouter()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "Payroll processed successfully")

	var payslip models.Payslip
	err := db.Where("employee_id = ?", emp.ID).First(&payslip).Error
	assert.Nil(t, err)
	assert.Greater(t, payslip.TotalTakeHome, 0.0)
}

func ptrTime(t time.Time) *time.Time {
	return &t
}

func TestRunPayrollAlreadyProcessed(t *testing.T) {
	db := pkg.ConnectTestDB()
	pkg.DB = db

	period := models.PayrollPeriod{
		ID:        2,
		StartDate: time.Now().AddDate(0, 0, -10),
		EndDate:   time.Now(),
		IsClosed:  true,
	}
	db.FirstOrCreate(&period, models.PayrollPeriod{
		ID: 2,
	})

	body := map[string]interface{}{"payroll_period_id": 1}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPost, "/admin/payroll/run", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router := setupPayrollRouter()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, resp.Body.String(), "Payroll already processed for this period")
}

func TestRunPayrollPeriodNotFound(t *testing.T) {
	db := pkg.ConnectTestDB()
	pkg.DB = db

	body := map[string]interface{}{"payroll_period_id": 999}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPost, "/admin/payroll/run", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router := setupPayrollRouter()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)
	assert.Contains(t, resp.Body.String(), "not found")
}
