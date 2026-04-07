package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mikasa/skk/internal/service"
)

const UserIDKey = "user_id"

func AuthMiddleware(authSvc *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			c.Abort()
			return
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			c.Abort()
			return
		}

		claims, err := authSvc.ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			c.Abort()
			return
		}

		c.Set(UserIDKey, claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}

// GetUserID extracts the authenticated user ID from gin.Context.
func GetUserID(c *gin.Context) (int64, bool) {
	id, exists := c.Get(UserIDKey)
	if !exists {
		return 0, false
	}
	userID, ok := id.(int64)
	return userID, ok
}
