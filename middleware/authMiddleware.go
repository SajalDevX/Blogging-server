package middleware

import (
	"context"
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
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Return the secret for validation
		return []byte(os.Getenv("ACCESS_TOKEN_SECRET")), nil
	})
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired access token"})
		c.Abort()
		return
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		c.Abort()
		return
	}
	if exp, ok := claims["exp"].(float64); !ok || float64(time.Now().Unix()) > exp {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has expired"})
		c.Abort()
		return
	}
	jti := claims["jti"].(string)
	if isTokenRevoked(jti) { // Implement this function to check for revoked tokens
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has been revoked"})
		c.Abort()
		return
	}
	sub, ok := claims["sub"].(float64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token subject"})
		c.Abort()
		return
	}
	var user models.User
	if err := initializers.DB.First(&user, uint(sub)).Error; err != nil || user.ID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		c.Abort()
		return
	}
	c.Set("user", user)
	c.Next()
}
func isTokenRevoked(jti string) bool {
	ctx := context.Background()
		result,err:=initializers.RedisClient.Get(ctx,jti).Result()
		if err!=nil{
			if err.Error() == "redis: nil" {
				return false
			}		}
		return result == "revoked"
}
