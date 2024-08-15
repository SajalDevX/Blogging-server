package models

import "gorm.io/gorm"

type User struct {
	Name     string `gorm:"not null"`
	Email    string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	Role     UserRole `gorm:"not null"`
	Bio      *string `gorm:"type:text"`
	ImageUrl *string `gorm:"type:text"`
	gorm.Model
}

type UserRole string

const (
	UserRoleAdmin  UserRole = "admin"
	UserRoleEditor UserRole = "editor"
	UserRoleAuthor UserRole = "author"
	UserRoleViewer UserRole = "viewer"
)
