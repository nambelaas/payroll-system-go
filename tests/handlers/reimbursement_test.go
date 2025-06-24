package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nambelaas/payroll-system-go/internal/handlers"
	"github.com/nambelaas/payroll-system-go/pkg"
	"github.com/stretchr/testify/assert"
)

func setupReimbursementRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})

	r.POST("/employee/reimbursement", handlers.SubmitReimbursement)

	return r
}

func TestSubmitReimbursementSuccess(t *testing.T) {
	db := pkg.ConnectTestDB()
	pkg.DB = db

	r := setupReimbursementRouter()

	body := map[string]any{
		"description": "Beli bensin",
		"amount":      30000.0,
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPost, "/employee/reimbursement", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)
	assert.Contains(t, resp.Body.String(), "Reimbursement submitted")
}

func TestSubmitReimbursementInvalidDescriptionOrAmount(t *testing.T) {
	db := pkg.ConnectTestDB()
	pkg.DB = db

	r := setupReimbursementRouter()

	body := map[string]any{
		"description": "",
		"amount":      0,
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPost, "/employee/reimbursement", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, resp.Body.String(), "Invalid description or amount")
}
