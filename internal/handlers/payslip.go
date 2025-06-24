package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nambelaas/payroll-system-go/internal/models"
	"github.com/nambelaas/payroll-system-go/pkg"
)

func GetPayslip(c *gin.Context) {
	userId, _ := c.Get("user_id")

	var emp models.Employee
	if err := pkg.DB.Where("user_id=?", userId).First(&emp).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Employee not found",
		})
		return
	}

	pidStr := c.Param("payroll_period_id")
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid period Id",
		})
		return
	}

	var ps models.Payslip
	if err := pkg.DB.Where("employee_id=? and payroll_period_id=?", emp.ID, pid).First(&ps).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Payslip not found",
		})
		return
	}

	var period models.PayrollPeriod
	if err := pkg.DB.First(&period, pid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Payroll period not found",
		})
		return
	}

	var reimbursements []models.Reimbursement
	pkg.DB.Where("employee_id=? and date between ? and ? and status =?", emp.ID, period.StartDate, period.EndDate, "approved").Find(&reimbursements)

	rbList := []gin.H{}
	for _, rb := range reimbursements {
		rbList = append(rbList, gin.H{
			"description": rb.Description,
			"amount":      rb.Amount,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"employee_id":       emp.ID,
		"payroll_period_id": pid,
		"attendance_hours":  ps.AttendanceHours,
		"prorated_salary":   ps.ProratedSalary,
		"overtime_hours":    ps.OvertimeHours,
		"overtime_pay":      ps.OvertimePay,
		"reimbursements":    rbList,
		"reimbursement_sum": ps.ReimbursementSum,
		"total_take_home":   ps.TotalTakeHome,
	})
}

func GetPayslipSummary(c *gin.Context) {
	pidStr := c.Param("payroll_period_id")
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid period Id",
		})
		return
	}

	var payslips []models.Payslip
	if err := pkg.DB.Where("payroll_period_id=?", pid).Find(&payslips).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get payslips",
		})
		return
	}

	summary := []gin.H{}
	total := 0.0

	for _, ps := range payslips {
		total += ps.TotalTakeHome
		summary = append(summary, gin.H{
			"employee_id":       ps.EmployeeId,
			"prorated_salary":   ps.ProratedSalary,
			"overtime_pay":      ps.OvertimePay,
			"reimbursement_sum": ps.ReimbursementSum,
			"total_take_home":   ps.TotalTakeHome,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"payroll_period_id":   pid,
		"total_take_home_all": total,
		"employees":           summary,
	})
}
