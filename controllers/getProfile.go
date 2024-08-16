package controllers

import (
	"main-module/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetProfile(c *gin.Context) {
	user, exist := c.Get("user")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
	}
	c.JSON(http.StatusOK,gin.H{
		"id":       user.(models.User).ID,
		"name":     user.(models.User).Name,
		"email":    user.(models.User).Email,
		"role":     user.(models.User).Role,
		"bio":      user.(models.User).Bio,
		"imageUrl": user.(models.User).ImageUrl,
	})
}
