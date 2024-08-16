package main

import (
	"main-module/controllers"
	"main-module/initializers"
	"main-module/middleware"
	"main-module/models"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
	initializers.SyncDatabase()
}

func main() {
	r := gin.Default()
	r.POST("/signup", controllers.SignUp)
	r.POST("/login", controllers.Login)
	r.GET("/profile", middleware.RequireAuth, middleware.RoleMiddleware(models.UserRoleAdmin, models.UserRoleAuthor, models.UserRoleEditor, models.UserRoleViewer), controllers.GetProfile)
	r.PUT("/profile/update", middleware.RequireAuth, middleware.RoleMiddleware(models.UserRoleAdmin, models.UserRoleAuthor, models.UserRoleEditor, models.UserRoleViewer), controllers.UpdateProfile)
	r.Static("/profile-image", "./build/resources/main/static/profile-image")

	r.Run()
}
