package models

import "gorm.io/gorm"

type PostStatus string

const (
	PostStatusDraft     PostStatus = "draft"
	PostStatusPublished PostStatus = "published"
	PostStatusArchived  PostStatus = "archived"
)

type Post struct {
	Title      string     `gorm:"not null"`
	Body       string     `gorm:"type:text"`
	CategoryID uint       `gorm:"index"`                // Foreign key for the post's category
	Tags       []Tag      `gorm:"many2many:post_tags;"` // Many-to-many relationship with tags
	AuthorID   uint       `gorm:"not null"`             // Foreign key for the author
	Status     PostStatus `gorm:"default:'draft'"`
	gorm.Model
}
type Category struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"unique;not null"`
}
type Tag struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"unique;not null"`
}
