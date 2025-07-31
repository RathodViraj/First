package middleware

import (
	"First/service"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func GinMiddleware(authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.FullPath() == "/login" || c.FullPath() == "/register" || c.FullPath() == "/" || c.FullPath() == "/home" || c.FullPath() == "/ws" {
			c.Next()
			return
		}

		authHead := c.GetHeader("Authorization")
		if authHead == "" || !strings.HasPrefix(authHead, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Authorization header required"})
			return
		}

		tokenstr := strings.TrimPrefix(authHead, "Bearer ")
		userID, role, err := authService.ValidateToken(tokenstr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
			return
		}

		c.Set("userID", userID)
		c.Set("role", role)
		c.Next()
	}
}
