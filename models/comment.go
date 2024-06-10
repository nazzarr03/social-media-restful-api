package models

import "time"

type Comment struct {
	CommentID       int       `json:"comment_id" gorm:"primaryKey;autoIncrement"`
	PostID          int       `json:"post_id"`
	UserID          int       `json:"user_id"`
	Content         string    `json:"content"`
	ImageURL        *string   `json:"image_url"`
	ParentCommentID *int      `json:"parent_comment"`
	Comments        []Comment `json:"comments" gorm:"foreignKey:ParentCommentID"`
	Likes           []Like    `json:"likes" gorm:"foreignKey:CommentID"`
	CreatedAt       time.Time `gorm:"default:current_timestamp"`
	UpdatedAt       time.Time `gorm:"default:current_timestamp"`
}
