package middleware

import (
	"fmt"
	"main-module/initializers"
	"main-module/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func RequireAuth(c *gin.Context) {
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "No authorization token provided",
		})
		c.Abort()
		return
	}
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Return the secret for validation
		return []byte(os.Getenv("SECRET")), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		exp, expOk := claims["exp"].(float64)
		if !expOk || float64(time.Now().Unix()) > exp {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token expired",
			})
			c.Abort()
			return
		}
		var user models.User
		sub, subOk := claims["sub"].(float64)
		if !subOk {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token subject"})
			c.Abort()
			return
		}
		initializers.DB.First(&user, uint(sub))
		if user.ID == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			c.Abort()
			return
		}
		c.Set("user", user)
		c.Next()
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}
}
