package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nambelaas/payroll-system-go/internal/handlers"
	"github.com/nambelaas/payroll-system-go/internal/models"
	"github.com/nambelaas/payroll-system-go/pkg"
	"github.com/stretchr/testify/assert"
)

func setupAuthTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/login", handlers.LoginHandler)

	return r
}

func TestLoginSuccess(t *testing.T) {
	db := pkg.ConnectTestDB()
	pkg.DB = db
	// db.AutoMigrate(&models.User{})

	// db.Exec("DELETE From users")

	// insert user ke db
	hashedPassword, _ := pkg.HashPassword("passwordTest123")
	db.FirstOrCreate(&models.User{
		ID:       1,
		Username: "testuser",
		Password: hashedPassword,
		Role:     "employee",
	}, models.User{
		Username: "testuser",
	})

	r := setupAuthTestRouter()

	body := map[string]string{
		"username": "testuser",
		"password": "passwordTest123",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "token")
}

func TestLoginUserNotFound(t *testing.T) {
	db := pkg.ConnectTestDB()
	pkg.DB = db
	db.AutoMigrate(&models.User{})

	r := setupAuthTestRouter()

	body := map[string]string{
		"username": "testuser1",
		"password": "passwordTest1234",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
	assert.Contains(t, resp.Body.String(), "User not found")
}

func TestLoginWrongPassword(t *testing.T) {
	db := pkg.ConnectTestDB()
	pkg.DB = db
	db.AutoMigrate(&models.User{})

	r := setupAuthTestRouter()

	body := map[string]string{
		"username": "testuser",
		"password": "passwordTest1234",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
	assert.Contains(t, resp.Body.String(), "Wrong password")
}
