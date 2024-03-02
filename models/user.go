package models

import "time"

type User struct {
	UserID    int       `json:"user_id" gorm:"primaryKey;autoIncrement"`
	Name      string    `json:"name" gorm:"not null"`
	Surname   string    `json:"surname" gorm:"not null"`
	Username  string    `json:"username" gorm:"unique;not null"`
	Password  string    `json:"-" gorm:"not null"`
	Email     string    `json:"email" gorm:"not null"`
	ImageURL  *string   `json:"image_url"`
	Posts     []Post    `json:"posts" gorm:"foreignKey:UserID"`
	Friends   []User    `json:"friends" gorm:"many2many:friendships;foreignKey:UserID;joinForeignKey:UserID;References:UserID;JoinReferences:FriendID"`
	CreatedAt time.Time `gorm:"default:current_timestamp"`
	UpdatedAt time.Time `gorm:"default:current_timestamp"`
}
