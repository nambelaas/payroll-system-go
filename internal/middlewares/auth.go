package middlewares

import (
	"net/http"
	"strings"

	"github.com/nambelaas/payroll-system-go/internal/utils"

	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Missing or invalid Authorization header",
			})
			c.Abort()
			return
		}

		token := strings.TrimPrefix(header, "Bearer ")
		claims, err := utils.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserId)
		c.Set("role", claims.Role)
		c.Next()
	}
}

func OnlyRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		r, _ := c.Get("role")
		if r != role {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Forbidden: " + role + " only",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
