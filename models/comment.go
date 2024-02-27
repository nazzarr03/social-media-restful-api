package models

import "time"

type Comment struct {
    CommentID      int       `json:"comment_id" gorm:"primaryKey;autoIncrement"`
    ParentCommentID *int     `json:"parent_comment_id"`
    UserID         int       `json:"user_id" gorm:"not null"`
    Content        string   `json:"content" gorm:"not null"`
    ImageURL       *string   `json:"image_url"`
    PostID         int       `json:"post_id" gorm:"not null"`
    Comments       []Comment `json:"comments" gorm:"foreignKey:ParentCommentID"`
    CreatedAt      time.Time `gorm:"default:current_timestamp"`
    UpdatedAt      time.Time `gorm:"default:current_timestamp"`
}