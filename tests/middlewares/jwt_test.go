package middlewares

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nambelaas/payroll-system-go/internal/handlers"
	"github.com/nambelaas/payroll-system-go/internal/middlewares"
	"github.com/nambelaas/payroll-system-go/internal/utils"
	"github.com/nambelaas/payroll-system-go/pkg"
	"github.com/stretchr/testify/assert"
)

func setupJwtRouter(token string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.Use(func(ctx *gin.Context) {
		if token != "" {
			ctx.Request.Header.Set("Authorization", token)
		}
	})

	r.Use(middlewares.JWTAuthMiddleware())
	r.POST("/employee/attendance/checkin", handlers.SubmitAttendance)

	return r
}

func TestAuthWithValidToken(t *testing.T) {
	db := pkg.ConnectTestDB()
	pkg.DB = db

	db.Exec("Delete from attendances")

	body := map[string]string{
		"username": "testuser",
		"password": "passwordTest123",
	}
	jsonBody, _ := json.Marshal(body)

	token, _ := utils.GenerateJWT(1, "employee")
	r := setupJwtRouter("Bearer " + token)

	req, _ := http.NewRequest(http.MethodPost, "/employee/attendance/checkin", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)
	assert.Contains(t, resp.Body.String(), "Check in success")
}

func TestAuthWithNoToken(t *testing.T) {
	body := map[string]string{
		"username": "testuser",
		"password": "passwordTest123",
	}
	jsonBody, _ := json.Marshal(body)

	r := setupJwtRouter("")

	req, _ := http.NewRequest(http.MethodPost, "/employee/attendance/checkin", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
	assert.Contains(t, resp.Body.String(), "Missing or invalid Authorization header")
}

func TestAuthWithInvalidToken(t *testing.T) {
	body := map[string]string{
		"username": "testuser",
		"password": "passwordTest123",
	}
	jsonBody, _ := json.Marshal(body)

	r := setupJwtRouter("Bearer invalidToken")

	req, _ := http.NewRequest(http.MethodPost, "/employee/attendance/checkin", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
	assert.Contains(t, resp.Body.String(), "Unauthorized")
}

func setupOnlyRoleRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	return r
}

func TestOnlyRoleSuccess(t *testing.T){
	db := pkg.ConnectTestDB()
	pkg.DB = db

	db.Exec("Delete from attendances")

	body := map[string]string{
		"username": "testuser",
		"password": "passwordTest123",
	}
	jsonBody, _ := json.Marshal(body)

	r := setupOnlyRoleRouter()

	r.Use(func(ctx *gin.Context) {
		ctx.Set("role", "employee")
		ctx.Set("user_id", 1)
		ctx.Next()
	})

	r.Use(middlewares.OnlyRole("employee"))
	
	r.POST("/employee/attendance/checkin", handlers.SubmitAttendance)

	req, _ := http.NewRequest(http.MethodPost, "/employee/attendance/checkin", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)
	assert.Contains(t, resp.Body.String(), "Check in success")
}

func TestOnlyRoleForbidden(t *testing.T){
	db := pkg.ConnectTestDB()
	pkg.DB = db

	db.Exec("Delete from attendances")

	body := map[string]string{
		"username": "testuser",
		"password": "passwordTest123",
	}
	jsonBody, _ := json.Marshal(body)

	r := setupOnlyRoleRouter()

	r.Use(func(ctx *gin.Context) {
		ctx.Set("role", "employee")
		ctx.Set("user_id", 1)
		ctx.Next()
	})

	r.Use(middlewares.OnlyRole("admin"))

	r.POST("/employee/attendance/checkin", handlers.SubmitAttendance)

	req, _ := http.NewRequest(http.MethodPost, "/employee/attendance/checkin", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusForbidden, resp.Code)
	assert.Contains(t, resp.Body.String(), "Forbidden: admin only")
}