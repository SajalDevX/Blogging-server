package posts

import (
	"main-module/initializers"
	"main-module/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreatePost(c *gin.Context) {
	var input struct {
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid input"})
		return
	}
	post := models.Post{
		Title:    input.Title,
		Body:  input.Content,
		AuthorID: user.(models.User).ID, 
	}
	result:=initializers.DB.Create(&post)
	if result.Error!=nil{
		c.JSON(http.StatusInternalServerError,gin.H{"error":"Failed to create Post"})
		return
	}

	c.JSON(http.StatusOK,gin.H{"post":post})

}
