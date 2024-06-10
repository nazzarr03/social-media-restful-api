package models

import "time"

type User struct {
	UserID    int       `json:"user_id" gorm:"primaryKey;autoIncrement"`
	Name      string    `json:"name"`
	Surname   string    `json:"surname"`
	Username  string    `json:"username" gorm:"unique"`
	Password  string    `json:"-"`
	Email     string    `json:"email"`
	ImageURL  *string   `json:"image_url"`
	Posts     []Post    `json:"posts" gorm:"foreignKey:UserID"`
	Friends   []User    `json:"friends" gorm:"many2many:friendships;foreignKey:UserID;joinForeignKey:UserID;References:UserID;JoinReferences:FriendID"`
	CreatedAt time.Time `gorm:"default:current_timestamp"`
	UpdatedAt time.Time `gorm:"default:current_timestamp"`
}
