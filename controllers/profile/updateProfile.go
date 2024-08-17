package controllers

import (
	"encoding/json"
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

type UpdateProfileRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Bio   string `json:"bio"`
}

func UpdateProfile(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	userModel := user.(models.User)

	// Extract JSON from the "profile" form field
	profile := c.PostForm("profile")
	if profile == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Profile data is required"})
		return
	}

	var req UpdateProfileRequest
	if err := json.Unmarshal([]byte(profile), &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	// Update name if provided and not empty
	if req.Name != "" {
		userModel.Name = req.Name
	}

	// Update email if provided and not empty
	if req.Email != "" {
		if !isValidEmail(req.Email) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
			return
		}

		var existingUser models.User
		if err := initializers.DB.Where("LOWER(email) = ?", req.Email).First(&existingUser).Error; err == nil && existingUser.ID != userModel.ID {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already in use"})
			return
		}
		userModel.Email = req.Email
	}

	// Update bio if provided and not empty
	if req.Bio != "" {
		userModel.Bio = &req.Bio // Assigning the address of bio (converting to *string)
	}

	imageUpdated:=false
	// Handle image file (optional part)
	imageFile, err := c.FormFile("image")
	if err == nil {
		if err := handleImageFile(imageFile, &userModel); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		imageUpdated = true
	}

	// Save updated user
	if err := initializers.DB.Save(&userModel).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}
	response := gin.H{
		"message":"Profile updated successfully",
	}
	if imageUpdated{
		response["message"] = "Image updated successfully"
	}

	c.JSON(http.StatusOK, response)
}

func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

func handleImageFile(file *multipart.FileHeader, userModel *models.User) error {
	fmt.Println("Saving file:", file.Filename)
	dir := filepath.Join("build", "resources", "main", "static", "profile-image")
	fmt.Println("Directory:", dir)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	dstPath := filepath.Join(dir, file.Filename)
	fmt.Println("File path:", dstPath)

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	imageURL := "/profile-image/" + file.Filename
	userModel.ImageUrl = &imageURL
	return nil
}
