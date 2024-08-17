package main

import (
	"github.com/gin-gonic/gin"
	auth "main-module/controllers/auth"
	profile "main-module/controllers/profile"
	post "main-module/controllers/posts"
	"main-module/initializers"
	"main-module/middleware"
	"main-module/models"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
	initializers.SyncDatabase()
}

func main() {
	r := gin.Default()
	//Authentication routes
	r.POST("/signup", auth.SignUp)
	r.POST("/login", auth.Login)

	//Profile routes
	r.GET("/profile", middleware.RequireAuth, middleware.RoleMiddleware(models.UserRoleAdmin, models.UserRoleAuthor, models.UserRoleEditor, models.UserRoleViewer), profile.GetProfile)
	r.PUT("/profile/update", middleware.RequireAuth, middleware.RoleMiddleware(models.UserRoleAdmin, models.UserRoleAuthor, models.UserRoleEditor, models.UserRoleViewer), profile.UpdateProfile)
	r.Static("/profile-image", "./build/resources/main/static/profile-image")

	//Post routes
	r.POST("/post/create",middleware.RequireAuth,middleware.RoleMiddleware(models.UserRoleAdmin, models.UserRoleAuthor, models.UserRoleEditor),post.CreatePost)
	r.GET("/post/:id", middleware.RequireAuth,post.GetPost)
	r.Run()
}
