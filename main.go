package main

import (
	"main-module/controllers"
	"main-module/initializers"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
	initializers.SyncDatabase()
}

func main() {
	r:=gin.Default()
	r.POST("/signup",controllers.SignUp)
	r.POST("/login",controllers.Login)
	r.Run()
}
