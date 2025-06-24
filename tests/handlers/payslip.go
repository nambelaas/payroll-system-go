package handlers

import (
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

func setupPayslipTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})

	r.GET("/employee/payslip/:payroll_period_id", handlers.GetPayslip)
	return r
}

func TestGetPayslipSuccess(t *testing.T) {
	db := pkg.ConnectTestDB()
	pkg.DB = db

	// Setup employee & user
	db.FirstOrCreate(&models.User{ID: 1, Username: "employee1", Password: "pass", Role: "employee"}, models.User{ID: 1})
	db.FirstOrCreate(&models.Employee{ID: 1, UserId: 1, Salary: 5000000}, models.Employee{
		ID:     1,
		UserId: 1,
	})

	// Setup payroll period
	db.FirstOrCreate(&models.PayrollPeriod{ID: 1, StartDate: time.Now().AddDate(0, 0, -30), EndDate: time.Now()}, models.PayrollPeriod{ID: 1})

	// Setup payslip
	db.FirstOrCreate(&models.Payslip{
		EmployeeId:       1,
		PayrollPeriodId:  1,
		AttendanceHours:  160,
		ProratedSalary:   5000000,
		OvertimeHours:    4,
		OvertimePay:      640000,
		ReimbursementSum: 100000,
		TotalTakeHome:    5740000,
	}, models.Payslip{
		EmployeeId:      1,
		PayrollPeriodId: 1,
	})

	// Setup reimbursement
	db.FirstOrCreate(&models.Reimbursement{
		EmployeeId:  1,
		Amount:      100000,
		Description: "Internet",
		Date:        time.Now(),
		Status:      "approved",
	})

	router := setupPayslipTestRouter()

	req, _ := http.NewRequest(http.MethodGet, "/employee/payslip/1", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), `"prorated_salary":5000000`)
	assert.Contains(t, resp.Body.String(), `"overtime_hours":4`)
	assert.Contains(t, resp.Body.String(), `"description":"Internet"`)
}

func TestGetPayslipNotFound(t *testing.T) {
	db := pkg.ConnectTestDB()
	pkg.DB = db

	db.AutoMigrate(&models.User{}, &models.Employee{}, &models.Payslip{})

	// Setup employee only
	db.FirstOrCreate(&models.User{ID: 1, Username: "employee1", Password: "pass", Role: "employee"}, models.User{
		ID:       1,
		Username: "employee1",
	})
	db.FirstOrCreate(&models.Employee{ID: 1, UserId: 1, Salary: 5000000}, models.Employee{
		ID:     1,
		UserId: 1,
	})

	router := setupPayslipTestRouter()

	req, _ := http.NewRequest(http.MethodGet, "/employee/payslip/99", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)
	assert.Contains(t, resp.Body.String(), "Payslip not found")
}

func TestGetPayslipEmployeeNotFound(t *testing.T) {
	db := pkg.ConnectTestDB()
	pkg.DB = db

	db.AutoMigrate(&models.User{}, &models.Employee{})

	db.FirstOrCreate(&models.User{ID: 1, Username: "ghost", Password: "pass", Role: "employee"}, models.User{
		ID:       1,
		Username: "ghost",
	})

	router := setupPayslipTestRouter()

	req, _ := http.NewRequest(http.MethodGet, "/employee/payslip/1", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)
	assert.Contains(t, resp.Body.String(), "Employee not found")
}
