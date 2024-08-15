package initializers

import (
	"main-module/models"
)

func SyncDatabase() {
	DB.AutoMigrate(&models.User{})
}