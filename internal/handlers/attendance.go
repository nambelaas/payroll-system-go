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

func SubmitAttendance(c *gin.Context) {
	userId, _ := c.Get("user_id")

	var emp models.Employee
	if err := pkg.DB.Where("user_id=?", userId).First(&emp).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Employee not found",
		})
		return
	}

	now := time.Now()
	if now.Weekday() == time.Saturday || now.Weekday() == time.Sunday {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Cannot submit attendance on weekends",
		})
		return
	}

	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	var existingAttendance models.Attendance
	err := pkg.DB.Where("employee_id = ? AND date >= ? AND date < ?", emp.ID, startOfDay, endOfDay).First(&existingAttendance).Error
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Already check in for today",
		})
		return
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database error",
		})
		return
	}

	att := models.Attendance{
		EmployeeId: emp.ID,
		Date:       now,
		CheckIn:    now,
	}
	pkg.DB.Create(&att)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Check in success",
		"time":    now.Format("15:04"),
	})
}

func SubmitCheckOut(c *gin.Context) {
	userId, _ := c.Get("user_id")

	var emp models.Employee
	if err := pkg.DB.Where("user_id=?", userId).First(&emp).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Employee not found",
		})
		return
	}

	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	var att models.Attendance
	err := pkg.DB.Where("employee_id = ? and date >= ? and date < ?", emp.ID, startOfDay, endOfDay).First(&att).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "You Haven't checked in today",
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database error",
		})
		return
	}

	if att.CheckOut != nil && !att.CheckOut.IsZero() {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "You have already check out",
		})
		return
	}

	att.CheckOut = &now
	pkg.DB.Save(&att)

	c.JSON(http.StatusOK, gin.H{
		"message": "Check out success",
		"time":    now.Format("15:04"),
	})
}
