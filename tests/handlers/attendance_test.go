package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nambelaas/payroll-system-go/internal/handlers"
	"github.com/nambelaas/payroll-system-go/internal/models"
	"github.com/nambelaas/payroll-system-go/pkg"
	"github.com/stretchr/testify/assert"
)

func setupAttendanceRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})

	r.POST("/employee/attendance/checkin", handlers.SubmitAttendance)
	r.POST("/employee/attendance/checkout", handlers.SubmitCheckOut)
	return r
}

func TestSubmitCheckInSuccess(t *testing.T) {
	db := pkg.ConnectTestDB()
	pkg.DB = db

	db.FirstOrCreate(&models.Employee{
		ID:     1,
		UserId: 1,
		Name:   "testing",
		Salary: 5000000,
	}, models.Employee{
		ID: 1,
	})

	db.Exec("delete from attendances")

	r := setupAttendanceRouter()

	req, _ := http.NewRequest(http.MethodPost, "/employee/attendance/checkin", nil)
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)
	assert.Contains(t, resp.Body.String(), "Check in success")
}

func TestSubmitCheckInAlreadySubmitted(t *testing.T) {
	db := pkg.ConnectTestDB()
	pkg.DB = db

	router := setupAttendanceRouter()

	req, _ := http.NewRequest(http.MethodPost, "/employee/attendance/checkin", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, resp.Body.String(), "Already check in for today")
}

func TestSubmitCheckOutSuccess(t *testing.T) {
	db := pkg.ConnectTestDB()
	pkg.DB = db

	db.FirstOrCreate(&models.Employee{
        ID:     1,
        UserId: 1,
        Name:   "testing",
        Salary: 5000000,
    }, models.Employee{
        ID: 1,
    })

	db.Exec("delete from attendances")

	r := setupAttendanceRouter()

	reqCheckIn, _ := http.NewRequest(http.MethodPost, "/employee/attendance/checkin", nil)
    respCheckIn := httptest.NewRecorder()
    r.ServeHTTP(respCheckIn, reqCheckIn)
    assert.Equal(t, http.StatusCreated, respCheckIn.Code)

	req, _ := http.NewRequest(http.MethodPost, "/employee/attendance/checkout", nil)
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "Check out success")
}
