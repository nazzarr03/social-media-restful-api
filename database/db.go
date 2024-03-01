package database

import (
	"fmt"
	"github.com/nazzarr03/social-media-restful-api/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Db *gorm.DB

const createSql = `
    CREATE TABLE IF NOT EXISTS blacklists (
        token VARCHAR(255) PRIMARY KEY,
        added_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );`

func init() {
	ConnectDB()
}

func ConnectDB() {
	var err error
	dsn := "host=localhost user=postgres password=password dbname=socialmedia port=5432 sslmode=disable"
	Db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	if err := Db.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{}, &models.Like{}, &models.Friendship{}); err != nil {
		panic("failed to migrate database")
	}

	result := Db.Exec(createSql)

	if result.Error != nil {
		fmt.Println("Table already exists")
	}

	fmt.Println("Database connected successfully!")

}
