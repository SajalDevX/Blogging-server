package controllers

import (
	"main-module/initializers"
	"main-module/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(c *gin.Context) {

	var userInput models.User

	if c.BindJSON(&userInput) != nil {
		c.JSON(400, gin.H{"message": "Invalid request"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(userInput.Password), 10)

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
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  user.ID,
		"role": user.Role,
		"exp":  time.Now().Add(time.Hour * 24830).Unix(),
	})
	secret := os.Getenv("SECRET")
	tokenString, err := token.SignedString([]byte(secret))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create token"})
		return
	}
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "/", "", false, true)

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
	result := initializers.DB.First(&user, "email=?", body.Email)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email or password"})
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  user.ID,
		"role": user.Role,
		"exp":  time.Now().Add(time.Hour * 24 * 30).Unix(),
	})
	secret := os.Getenv("SECRET")
	tokenString, err := token.SignedString([]byte(secret))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create token"})
		return
	}

	// Send it back in a cookie
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logged in successfully"})
}
