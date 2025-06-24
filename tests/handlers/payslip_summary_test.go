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

func setupSummaryRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(999))
		c.Next()
	})

	r.GET("/admin/payslip/summary/:payroll_period_id", handlers.GetPayslipSummary)
	return r
}

func TestGetPayslipSummarySuccess(t *testing.T) {
	db := pkg.ConnectTestDB()
	pkg.DB = db

	hashedPassword, _ := pkg.HashPassword("passwordTest123")
	db.FirstOrCreate(&models.User{ID: 3, Username: "testing3", Password: hashedPassword, Role: "employee"}, models.User{
		ID: 3,
	})

	// Payroll period
	db.FirstOrCreate(&models.PayrollPeriod{ID: 1, StartDate: time.Now().AddDate(0, -1, 0), EndDate: time.Now()}, models.PayrollPeriod{ID: 1})

	// Employee dan payslip
	db.FirstOrCreate(&models.Employee{ID: 1, UserId: 1, Salary: 5000000}, models.Employee{ID: 1, UserId: 1})
	db.FirstOrCreate(&models.Employee{ID: 3, UserId: 3, Salary: 6000000}, models.Employee{ID: 3, UserId: 3})

	db.Exec("DELETE FROM payslips")

	db.FirstOrCreate(&models.Payslip{
		PayrollPeriodId:  1,
		EmployeeId:       1,
		ProratedSalary:   5000000,
		OvertimePay:      300000,
		ReimbursementSum: 100000,
		TotalTakeHome:    5400000,
	}, models.Payslip{
		EmployeeId:      1,
		PayrollPeriodId: 1,
	})

	db.FirstOrCreate(&models.Payslip{
		PayrollPeriodId:  1,
		EmployeeId:       3,
		ProratedSalary:   6000000,
		OvertimePay:      0,
		ReimbursementSum: 0,
		TotalTakeHome:    6000000,
	}, models.Payslip{
		EmployeeId:      3,
		PayrollPeriodId: 1,
	})

	router := setupSummaryRouter()
	req, _ := http.NewRequest(http.MethodGet, "/admin/payslip/summary/1", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), `"total_take_home_all":11400000`)
	assert.Contains(t, resp.Body.String(), `"employee_id":1`)
	assert.Contains(t, resp.Body.String(), `"employee_id":3`)
}

func TestGetPayslipSummaryNoData(t *testing.T) {
	db := pkg.ConnectTestDB()
	pkg.DB = db

	db.AutoMigrate(&models.PayrollPeriod{}, &models.Payslip{})

	db.FirstOrCreate(&models.PayrollPeriod{ID: 9, StartDate: time.Now().AddDate(0, -1, 0), EndDate: time.Now()}, models.PayrollPeriod{
		ID: 9,
	})

	router := setupSummaryRouter()
	req, _ := http.NewRequest(http.MethodGet, "/admin/payslip/summary/9", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), `"employees":[]`)
}
