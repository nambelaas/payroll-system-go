package main

import (
	"github.com/nambelaas/payroll-system-go/internal/handlers"
	"github.com/nambelaas/payroll-system-go/internal/middlewares"
	"github.com/nambelaas/payroll-system-go/pkg"

	"github.com/gin-gonic/gin"
)

func main() {
	pkg.ConnectDB()
	r := gin.Default()

	r.POST("/login", handlers.LoginHandler)

	// protected := r.Group("/")
	// protected.Use(middlewares.JWTAuthMiddleware())

	// protected.GET("/whoami", func(c *gin.Context) {
	// 	role, _ := c.Get("role")
	// 	userId, _ := c.Get("user_id")
	// 	c.JSON(200, gin.H{
	// 		"role":    role,
	// 		"user_id": userId,
	// 	})
	// })

	// admin
	admin := r.Group("/admin")
	admin.Use(middlewares.JWTAuthMiddleware(), middlewares.OnlyRole("admin"), middlewares.LogRequest())
	admin.POST("/payroll-period", handlers.CreatePayrollPeriod)
	admin.POST("/payroll/run", handlers.RunPayroll)
	admin.GET("/payslip/summary/:payroll_period_id", handlers.GetPayslipSummary)

	// employee
	employee := r.Group("/employee")
	employee.Use(middlewares.JWTAuthMiddleware(), middlewares.OnlyRole("employee"), middlewares.LogRequest())
	employee.POST("/attendance/checkin", handlers.SubmitAttendance)
	employee.POST("/attendance/checkout", handlers.SubmitCheckOut)
	employee.POST("/overtime", handlers.SubmitOvertime)
	employee.POST("/reimbursement", handlers.SubmitReimbursement)
	employee.GET("/payslip/:payroll_period_id", handlers.GetPayslip)

	r.Run(":8080")
}
