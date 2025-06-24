package handlers

import (
	"net/http"
	"github.com/nambelaas/payroll-system-go/internal/models"
	"github.com/nambelaas/payroll-system-go/pkg"
	"time"

	"github.com/gin-gonic/gin"
)

type ReimbursementRequest struct {
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
}

func SubmitReimbursement(c *gin.Context) {
	userId, _ := c.Get("user_id")

	var emp models.Employee
	if err := pkg.DB.Where("user_id=?", userId).First(&emp).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Employee not found",
		})
		return
	}

	var req ReimbursementRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Amount <= 0 || req.Description == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid description or amount",
		})
		return
	}

	rb := models.Reimbursement{
		EmployeeId:  emp.ID,
		Description: req.Description,
		Amount:      req.Amount,
		Status:      "pending",
		Date:        time.Now(),
	}

	if err := pkg.DB.Create(&rb).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to submit reimbursement",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Reimbursement submitted",
		"id":      rb.ID,
	})
}
