package handlers

import (
	"math"
	"net/http"
	"time"

	"github.com/nambelaas/payroll-system-go/internal/models"
	"github.com/nambelaas/payroll-system-go/pkg"

	"github.com/gin-gonic/gin"
)

type PayrollPeriodRequest struct {
	StartDate string `json:"start_date" binding:"required"` // format: yyyy-mm-dd
	EndDate   string `json:"end_date" binding:"required"`   // format: yyyy-mm-dd
}

func CreatePayrollPeriod(c *gin.Context) {
	userId, _ := c.Get("user_id")
	var req PayrollPeriodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input format",
		})
		return
	}

	startDate, err1 := time.Parse("2006-01-02", req.StartDate)
	endDate, err2 := time.Parse("2006-01-02", req.EndDate)

	if err1 != nil || err2 != nil || endDate.Before(startDate) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid date range",
		})
		return
	}

	var user models.User
	err := pkg.DB.First(&user, userId).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User not found",
		})
		return
	}

	period := models.PayrollPeriod{
		StartDate: startDate,
		EndDate:   endDate,
		IsClosed:  false,
		CreatedBy: user.Username,
	}
	if err := pkg.DB.Create(&period).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create payroll period",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Payroll period created",
		"id":      period.ID,
	})
}

type RunPayrollRequest struct {
	PayrollPeriodID uint `json:"payroll_period_id"`
}

func RunPayroll(c *gin.Context) {
	userId, _ := c.Get("user_id")
	var req RunPayrollRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Request",
		})
		return
	}

	var period models.PayrollPeriod
	if err := pkg.DB.First(&period, req.PayrollPeriodID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Payroll period not found",
		})
		return
	}

	if period.IsClosed {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Payroll already processed for this period",
		})
		return
	}

	var employees []models.Employee
	pkg.DB.Find(&employees)

	workdays := countWeekdays(period.StartDate, period.EndDate)
	totalHoursPerMonth := float64(workdays * 8)

	var user models.User
	err := pkg.DB.First(&user, userId).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User not found",
		})
		return
	}

	now := time.Now()
	updateTime := time.Date(now.Year(), now.Month(), now.Day(), now.Day(), now.Minute(), now.Second(), 0, now.Location())

	for _, emp := range employees {
		attHours := 0.0
		var attendances []models.Attendance
		pkg.DB.Where("employee_id = ? and date between ? and ?", emp.ID, period.StartDate, period.EndDate).Find(&attendances)

		for _, att := range attendances {
			if att.CheckOut != nil && !att.CheckOut.IsZero() {
				durasi := att.CheckOut.Sub(att.CheckIn).Hours()
				attHours += durasi
			}
		}

		var overtimes []models.Overtime
		pkg.DB.Where("employee_id = ? and date between ? and ?", emp.ID, period.StartDate, period.EndDate).Find(&overtimes)

		otHours := 0.0
		for _, ot := range overtimes {
			otHours += ot.Hours
		}

		var reimbursements []models.Reimbursement
		pkg.DB.Where("employee_id = ? and date between ? and ? and status = ?", emp.ID, period.StartDate, period.EndDate, "pending").Find(&reimbursements)

		rbSum := 0.0
		for _, r := range reimbursements {
			rbSum += r.Amount
			r.Status = "approved"
			r.UpdatedBy = &user.Username
			r.UpdatedAt = &updateTime
			pkg.DB.Save(&r)
		}

		perHour := emp.Salary / totalHoursPerMonth
		proratedSalary := math.Round(perHour*attHours*100) / 100
		overtimePay := math.Round(otHours*perHour*2*100) / 100
		rbSum = math.Round(rbSum*100) / 100
		total := math.Round((proratedSalary+overtimePay+rbSum)*100) / 100

		payslip := models.Payslip{
			EmployeeId:       emp.ID,
			PayrollPeriodId:  period.ID,
			AttendanceHours:  attHours,
			ProratedSalary:   proratedSalary,
			OvertimeHours:    otHours,
			OvertimePay:      overtimePay,
			ReimbursementSum: rbSum,
			TotalTakeHome:    total,
			CreatedBy:        user.Username,
		}
		pkg.DB.Create(&payslip)
	}

	period.IsClosed = true
	period.UpdatedBy = &user.Username
	period.UpdatedAt = &updateTime
	pkg.DB.Save(&period)

	c.JSON(http.StatusOK, gin.H{
		"message": "Payroll processed successfully",
	})
}

func countWeekdays(start, end time.Time) int {
	count := 0
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		if d.Weekday() != time.Saturday && d.Weekday() != time.Sunday {
			count++
		}
	}
	return count
}
