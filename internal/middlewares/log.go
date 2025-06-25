package middlewares

import (
	"bytes"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nambelaas/payroll-system-go/internal/models"
	"github.com/nambelaas/payroll-system-go/internal/utils"
	"github.com/nambelaas/payroll-system-go/pkg"
)

type bodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func LogRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		bw := &bodyWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = bw
		userId, _ := c.Get("user_id")

		c.Next()

		uid, _ := userId.(uint)
		requestId := uuid.New().String()
		endpoint := c.FullPath()
		myIp, _ := utils.GetMyIp()
		status := c.Writer.Status()
		responseBody := bw.body.String()

		auditLog := models.AuditLog{
			RequestId: requestId,
			UserId:    uid,
			Endpoint:  endpoint,
			Status:    status,
			Result:    responseBody,
			IpAddress: myIp,
		}

		pkg.DB.Create(&auditLog)
	}
}
