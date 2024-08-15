package models

import (
	"gorm.io/gorm"
)

type CommentStatus string

const (
	CommentStatusPending  CommentStatus = "pending"
	CommentStatusApproved CommentStatus = "approved"
	CommentStatusSpam     CommentStatus = "spam"
)

type Comment struct {
	PostId uint  `gorm:"index;not null"`  
	AuthorId uint  `gorm:"index;not null"` 
	Body  string `gorm:"type:text;not null"`
	ParentCommentID *uint `gorm:"index"`   
	gorm.Model
}


