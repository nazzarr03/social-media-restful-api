package models

import (
	"time"
)

type Like struct {
	LikeID    int       `json:"like_id" gorm:"primaryKey;autoIncrement"`
	UserID    int       `json:"user_id"`
	PostID    int       `json:"post_id"`
	CommentID *int      `json:"comment_id"`
	IsLiked   bool      `json:"is_liked" gorm:"default:false"`
	CreatedAt time.Time `gorm:"default:current_timestamp"`
	UpdatedAt time.Time `gorm:"default:current_timestamp"`
}
