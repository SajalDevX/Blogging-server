package models

import "gorm.io/gorm"

type PostStatus string

const (
	PostStatusDraft     PostStatus = "draft"
	PostStatusPublished PostStatus = "published"
	PostStatusArchived  PostStatus = "archived"
)

type Post struct {
	gorm.Model
	Title      string     `gorm:"not null"`
	Body       string     `gorm:"type:text"`
	CategoryID uint       `gorm:"index"`               // Foreign key for the post's category
	Category   Category   `gorm:"foreignKey:CategoryID"` // Include this to preload the Category
	Tags       []Tag      `gorm:"many2many:post_tags;"`
	AuthorID   uint       `gorm:"not null"`            // Foreign key for the author
	Author     User       `gorm:"foreignKey:AuthorID"` // Include this to preload the Author
	Status     PostStatus `gorm:"default:'draft'"`
}

type Category struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"unique;not null"`
}
type Tag struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"unique;not null"`
}
