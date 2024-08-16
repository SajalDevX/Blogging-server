package controllers

import (
	"fmt"
	"io"
	"main-module/initializers"
	"main-module/models"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"

	"github.com/gin-gonic/gin"
)

func UpdateProfile(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	userModel := user.(models.User)

	// Parse form
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Extract profile fields
	name := c.PostForm("name")
	email := c.PostForm("email")
	bio := c.PostForm("bio")
	imageFiles := form.File["files"]

	// Validate name
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Field name"})
		return
	}
	userModel.Name = name

	// Validate and check email
	if email != "" {
		if !isValidEmail(email) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
			return
		}

		var existingUser models.User
		if err := initializers.DB.Where("LOWER(email) = ?", email).First(&existingUser).Error; err == nil && existingUser.ID != userModel.ID {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already in use"})
			return
		}
		userModel.Email = email
	}

	// Update bio if provided
	if bio != "" {
		userModel.Bio = &bio // Assigning the address of bio (converting to *string)
	}

	// Handle image file
	if len(imageFiles) > 0 {
		file := imageFiles[0]

		// Validate file size (example: max 10MB)
		if file.Size > 10*1024*1024 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "File size exceeds 10MB"})
			return
		}

		// Validate file extension (example: only images)
		ext := filepath.Ext(file.Filename)
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" && ext != ".webp" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type"})
			return
		}

		if err := saveFile(file); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image file"})
			return
		}

		imageURL := "/profile-image/" + file.Filename
		userModel.ImageUrl = &imageURL
	}

	// Save updated user
	db := initializers.DB
	if err := db.Save(&userModel).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully", "image_url": userModel.ImageUrl})
}

// isValidEmail validates email format
func isValidEmail(email string) bool {
	// Simple email validation regex
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

func saveFile(fileHeader *multipart.FileHeader) error {
    fmt.Println("Saving file:", fileHeader.Filename)
    dir := filepath.Join("build", "resources", "main", "static", "profile-image")
    fmt.Println("Directory:", dir)

    if err := os.MkdirAll(dir, 0755); err != nil {
        return err
    }

    dstPath := filepath.Join(dir, fileHeader.Filename)
    fmt.Println("File path:", dstPath)

    src, err := fileHeader.Open()
    if err != nil {
        return err
    }
    defer src.Close()

    dst, err := os.Create(dstPath)
    if err != nil {
        return err
    }
    defer dst.Close()

    _, err = io.Copy(dst, src)
    if err != nil {
        return err
    }
    return nil
}
