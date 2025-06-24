package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/nambelaas/payroll-system-go/internal/models"
	"github.com/nambelaas/payroll-system-go/pkg"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

type OvertimeRequest struct {
	Date      string `json:"date"`       // format: yyyy-mm-dd
	StartTime string `json:"start_time"` // "18:00"
	EndTime   string `json:"end_time"`   // "20:00"
}

func SubmitOvertime(c *gin.Context) {
	userId, _ := c.Get("user_id")

	var emp models.Employee
	if err := pkg.DB.Where("user_id=?", userId).First(&emp).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Employee not found",
		})
		return
	}

	var req OvertimeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input",
		})
		return
	}

	date, err1 := time.Parse("2006-01-02", req.Date)
	startTime, err2 := time.Parse("15:04", req.StartTime)
	endTime, err3 := time.Parse("15:04", req.EndTime)
	if err1 != nil || err2 != nil || err3 != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid date or time format",
		})
		return
	}

	start := time.Date(date.Year(), date.Month(), date.Day(), startTime.Hour(), startTime.Minute(), 0, 0, time.Local)
	end := time.Date(date.Year(), date.Month(), date.Day(), endTime.Hour(), endTime.Minute(), 0, 0, time.Local)

	if end.Before(start) || end.Equal(start) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "End time must be after start time",
		})
		return
	}

	// cek durasi overtime
	duration := end.Sub(start).Hours()
	if duration > 3 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Overtime cannot exceed 3 hours",
		})
		return
	}

	// cek apakah overtime bukan untuk besok
	if date.After(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Cannot submit overtime for future",
		})
		return
	}

	// cek apakah overtime setelah jam 5
	today := time.Now()
	if date.Year() == today.Year() && date.YearDay() == today.YearDay() {
		workEnd := time.Date(today.Year(), today.Month(), today.Day(), 17, 0, 0, 0, time.Local)
		if time.Now().Before(workEnd) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Can only submit overtime after 5PM",
			})
			return
		}
	}

	// cek sudah ambil overtime pada tanggal tersebut
	var existingOvertime models.Overtime
	err := pkg.DB.Where("employee_id=? and date = ?", emp.ID, date).First(&existingOvertime).Error
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Already submitted overtime on that date",
		})
		return
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database error",
		})
		return
	}

	ot := models.Overtime{
		EmployeeId: emp.ID,
		Date:       date,
		StartTime:  start,
		EndTime:    end,
		Hours:      duration,
	}
	pkg.DB.Create(&ot)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Overtime submitted",
	})
}
