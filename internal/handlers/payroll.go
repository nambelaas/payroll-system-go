package handlers

import (
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

	period := models.PayrollPeriod{
		StartDate: startDate,
		EndDate:   endDate,
		IsClosed:  false,
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

	for _, emp := range employees {
		var attendances []models.Attendance
		pkg.DB.Where("employee_id = ? and date between ? and ?", emp.ID, period.StartDate, period.EndDate).Find(&attendances)

		attHours := float64(len(attendances) * 8)
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
			pkg.DB.Save(&r)
		}

		perHour := emp.Salary / totalHoursPerMonth
		proratedSalary := perHour * attHours
		overtimePay := otHours * perHour * 2
		total := proratedSalary + overtimePay + rbSum

		payslip := models.Payslip{
			EmployeeId:       emp.ID,
			PayrollPeriodId:  period.ID,
			AttendanceHours:  attHours,
			ProratedSalary:   proratedSalary,
			OvertimeHours:    otHours,
			OvertimePay:      overtimePay,
			ReimbursementSum: rbSum,
			TotalTakeHome:    total,
		}
		pkg.DB.Create(&payslip)

		period.IsClosed = true
		pkg.DB.Save(&period)

		c.JSON(http.StatusOK, gin.H{
			"error": "Payroll processed successfully",
		})
	}
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
