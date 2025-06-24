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

func setupOvertimeRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})

	r.POST("/employee/overtime", handlers.SubmitOvertime)

	return r
}

func TestSubmitOvertimeSuccess(t *testing.T) {
	db := pkg.ConnectTestDB()
	pkg.DB = db

	r := setupOvertimeRouter()

	body := map[string]any{
		"date":       "2025-06-24",
		"start_time": "18:00",
		"end_time":   "20:00",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPost, "/employee/overtime", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)
	assert.Contains(t, resp.Body.String(), "Overtime submitted")
}

func TestSubmitOvertimeExceed3Hour(t *testing.T) {
	db := pkg.ConnectTestDB()
	pkg.DB = db

	db.Exec("DELETE FROM overtimes")
	
	r := setupOvertimeRouter()

	body := map[string]any{
		"date":       "2025-06-24",
		"start_time": "18:00",
		"end_time":   "22:00",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPost, "/employee/overtime", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, resp.Body.String(), "Overtime cannot exceed 3 hours")
}

func TestSubmitOvertimeForTomorrow(t *testing.T) {
	db := pkg.ConnectTestDB()
	pkg.DB = db

	r := setupOvertimeRouter()

	body := map[string]any{
		"date":       "2025-06-25",
		"start_time": "18:00",
		"end_time":   "20:00",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPost, "/employee/overtime", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, resp.Body.String(), "Cannot submit overtime for future")
}

func TestSubmitOvertimeInvalidTime(t *testing.T) {
	db := pkg.ConnectTestDB()
	pkg.DB = db

	r := setupOvertimeRouter()

	body := map[string]any{
		"date":       "2025-06-24",
		"start_time": "18:00",
		"end_time":   "17:00",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPost, "/employee/overtime", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, resp.Body.String(), "End time must be after start time")
}
