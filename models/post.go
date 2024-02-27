package models

import "time"

type Post struct {
	PostID   int    `json:"post_id" gorm:"primaryKey;autoIncrement"`
	UserID   int    `json:"user_id" gorm:"not null"`
	Content  string `json:"content" gorm:"not null"`
	ImageURL *string `json:"image_url"`
	Comments []Comment `json:"comments" gorm:"foreignKey:PostID"`
	CreatedAt time.Time `gorm:"default:current_timestamp"`
	UpdatedAt time.Time `gorm:"default:current_timestamp"`
}