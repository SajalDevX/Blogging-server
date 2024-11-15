package controllers

import (
	"main-module/initializers"
	"main-module/services"
	"main-module/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(c *gin.Context) {

	var userInput models.User

	if c.BindJSON(&userInput) != nil {
		c.JSON(400, gin.H{"message": "Invalid request"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(userInput.Password), bcrypt.DefaultCost+2)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Error hashing password"})
		return
	}
	user := models.User{
		Name:     userInput.Name,
		Email:    userInput.Email,
		Password: string(hash),
		Role:     userInput.Role,
		Bio:      userInput.Bio,
		ImageUrl: userInput.ImageUrl,
	}
	result := initializers.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create user"})
		return
	}
	accessToken, err := generateToken(user.ID, user.Role, os.Getenv("ACCESS_TOKEN_SECRET"), time.Minute*15)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create access token"})
	}
	refreshToken, err := generateToken(user.ID, user.Role, os.Getenv("REFRESH_TOKEN_SECRET"), time.Hour*24*30)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create refresh token"})
		return
	}

	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie("AccessToken", accessToken, 15*60, "/", "", false, true)
	c.SetCookie("RefreshToken", refreshToken, 30*24*3600, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})

}

func Login(c *gin.Context) {
	var body struct {
		Email    string `json:"email" binding:"required"`
		Passowrd string `json:"password" binding:"required"`
	}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}
	var user models.User
	if err := initializers.DB.First(&user, "email=?", body.Email).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Passowrd)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	accessToken, err := generateToken(user.ID, user.Role, os.Getenv("ACCESS_TOKEN_SECRET"), time.Minute*15)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create access token"})
		return
	}
	refreshToken, err := generateToken(user.ID, user.Role, os.Getenv("REFRESH_TOKEN_SECRET"), time.Hour*24*30)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create refresh token"})
		return
	}

	// Send it back in a cookie
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie("AccessToken", accessToken, 15*60, "/", "", false, true)
	c.SetCookie("RefreshToken", refreshToken, 30*24*3600, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "Logged in successfully"})
}

func generateToken(userID uint, role models.UserRole, secret string, expiration time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  userID,
		"role": role,
		"jti":  uuid.NewString(),
		"exp":  time.Now().Add(expiration).Unix(),
	})
	return token.SignedString([]byte(secret))
}

func Logout(c *gin.Context) {
	tokenString, err := c.Cookie("AccessToken")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No access token found"})
		return
	}

	// Parse the token to extract claims
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("ACCESS_TOKEN_SECRET")), nil
	})
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		return
	}

	// Revoke the token by storing its jti in Redis
	jti := claims["jti"].(string)
	exp := time.Unix(int64(claims["exp"].(float64)), 0) // Token expiration time
	if err := services.RevokeToken(jti, exp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to revoke token"})
		return
	}

	// Clear the token cookies
	c.SetCookie("AccessToken", "", -1, "/", "", false, true)
	c.SetCookie("RefreshToken", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}
