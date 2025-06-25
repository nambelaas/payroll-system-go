package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nambelaas/payroll-system-go/internal/handlers"
	"github.com/nambelaas/payroll-system-go/internal/middlewares"
	"github.com/nambelaas/payroll-system-go/internal/models"
	"github.com/nambelaas/payroll-system-go/pkg"
	"github.com/stretchr/testify/assert"
)

func setupAuditLogRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})
	r.Use(middlewares.LogRequest())

	r.POST("/employee/attendance/checkin", handlers.SubmitAttendance)

	return r
}

func TestLogRequestMiddleware(t *testing.T) {
	db := pkg.ConnectTestDB()
	pkg.DB = db

	db.Exec("Delete from audit_logs")
	db.Exec("Delete from attendances")

	r := setupAuditLogRouter()

	req, _ := http.NewRequest(http.MethodPost, "/employee/attendance/checkin", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)

	var log models.AuditLog
	err := db.Last(&log).Error

	assert.NoError(t, err)
	assert.Equal(t, uint(1), log.UserId)
	assert.Equal(t, "/employee/attendance/checkin", log.Endpoint)
	assert.Equal(t, 201, log.Status)
	assert.Contains(t, log.Result, "Check in success")
	assert.NotEmpty(t, log.RequestId)
}
