package initializers

import (
	"main-module/models"
)

func SyncDatabase() {
	DB.AutoMigrate(
		&models.Post{}, 
		&models.Category{}, 
		&models.Tag{},
	)
}