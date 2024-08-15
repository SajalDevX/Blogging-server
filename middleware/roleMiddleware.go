package middleware

import (
	"main-module/models"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func RoleMiddleware(allowedRoles ...models.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("Authorization")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No token found"})
			c.Abort()
			return
		}
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SECRET")), nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}
		userRoleStr := claims["role"].(string)
		userRole := models.UserRole(userRoleStr)
		for _, role := range allowedRoles {
			if role == userRole {
				c.Next()
				return
			}
		}
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to access this route"})
		c.Abort()
	}
}
