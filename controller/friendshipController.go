package controller

import (
	"github.com/gofiber/fiber"
	"github.com/nazzarr03/social-media-restful-api/database"
	"github.com/nazzarr03/social-media-restful-api/models"
)

func AddFriend(c *fiber.Ctx) {
	user := models.User{}
	friend := models.User{}

	database.Db.First(&user, "user_id = ?", c.Params("userid"))

	if user.UserID == 0 {
		c.Status(404).JSON(fiber.Map{
			"message": "User not found",
		})
		return
	}

	database.Db.First(&friend, "user_id = ?", c.Params("friendid"))

	if friend.UserID == 0 {
		c.Status(404).JSON(fiber.Map{
			"message": "Friend not found",
		})
		return
	}


	friendship := models.Friendship{
		UserID:   user.UserID,
		FriendID: friend.UserID,
		IsActive: true,
	}

	database.Db.Create(&friendship)

	if database.Db.Error != nil {
		c.Status(500).JSON(fiber.Map{
			"message": "Error creating friendship",
		})
		return
	}

	c.Status(200).JSON(fiber.Map{
		"message": "Friendship created successfully",
	})
}

func RemoveFriend(c *fiber.Ctx) {
	user := models.User{}
	friend := models.User{}

	database.Db.First(&user, "user_id = ?", c.Params("userid"))

	if user.UserID == 0 {
		c.Status(404).JSON(fiber.Map{
			"message": "User not found",
		})
		return
	}

	database.Db.First(&friend, "user_id = ?", c.Params("friendid"))

	if friend.UserID == 0 {
		c.Status(404).JSON(fiber.Map{
			"message": "Friend not found",
		})
		return
	}

	friendship := models.Friendship{}

	database.Db.Where("user_id = ? AND friend_id = ?", user.UserID, friend.UserID).First(&friendship)

	if friendship.FriendshipID == 0 {
		c.Status(404).JSON(fiber.Map{
			"message": "Friendship not found",
		})
		return
	}

	database.Db.Delete(&friendship)

	if database.Db.Error != nil {
		c.Status(500).JSON(fiber.Map{
			"message": "Error deleting friendship",
		})
		return
	}

	c.Status(200).JSON(fiber.Map{
		"message": "Friendship deleted successfully",
	})

}