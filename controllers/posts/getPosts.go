package posts

import (
	"main-module/initializers"
	"main-module/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetPost(c *gin.Context) {
	id := c.Param("id")

	postID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	var post models.Post

	// Query the database for the post with the given ID, preloading related data
	if result := initializers.DB.Preload("Tags").Preload("Category").Preload("Author").First(&post, postID); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"post": post})
}
