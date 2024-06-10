package controller

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/gofiber/fiber"
	"github.com/nazzarr03/social-media-restful-api/database"
	"github.com/nazzarr03/social-media-restful-api/dto"
	"github.com/nazzarr03/social-media-restful-api/models"
	"github.com/nazzarr03/social-media-restful-api/utils"
)

func CreatePost(c *fiber.Ctx) {
	var request dto.CreatePostRequest

	post := models.Post{}
	user := models.User{}

	if err := c.BodyParser(&request); err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
		return
	}

	if request.Content == "" {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Content cannot be empty",
		})
		return
	}

	file, err := c.FormFile("post_picture")
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
		return
	}

	tempFilePath := fmt.Sprintf("./uploads/%s", file.Filename)

	if err := c.SaveFile(file, tempFilePath); err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot save file",
		})
		return
	}

	cld, err := utils.ConnectToCloudinary()
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
		return
	}

	resp, err := utils.UploadToCloudinary(cld, tempFilePath)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
		return
	}

	os.Remove(tempFilePath)

	post.ImageURL = &resp.SecureURL

	database.Db.First(&user, "user_id = ?", c.Params("id"))

	if user.UserID == 0 {
		c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
		return
	}

	post.Content = request.Content
	post.UserID = user.UserID

	database.Db.Create(&post)

	if database.Db.Error != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot create post",
		})
		return
	}

	c.Status(fiber.StatusCreated).JSON(post)
}

func EditPost(c *fiber.Ctx) {
	var request dto.EditPostRequest

	post := models.Post{}
	user := models.User{}

	if err := c.BodyParser(&request); err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
		return
	}

	database.Db.First(&user, "user_id = ?", c.Params("userid"))

	if user.UserID == 0 {
		c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
		return
	}

	database.Db.First(&post, "post_id = ?", c.Params("postid"))

	if post.PostID == 0 {
		c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Post not found",
		})
		return
	}

	if request.Content == "" {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Content cannot be empty",
		})
		return
	}

	file, err := c.FormFile("post_picture")

	if err == nil {
		cld, err := utils.ConnectToCloudinary()
		if err != nil {
			c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
			return
		}

		if post.ImageURL != nil {
			lastPart := path.Base(*post.ImageURL)

			publicID := strings.Split(lastPart, ".")[0]

			if err := utils.DeleteFromCloudinary(cld, publicID); err != nil {
				c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": err.Error(),
				})
				return
			}
		}

		tempFilePath := fmt.Sprintf("./uploads/%s", file.Filename)

		if err := c.SaveFile(file, tempFilePath); err != nil {
			c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Cannot save file",
			})
			return
		}

		resp, err := utils.UploadToCloudinary(cld, tempFilePath)
		if err != nil {
			c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
			return
		}

		os.Remove(tempFilePath)

		post.ImageURL = &resp.SecureURL
	}

	post.Content = request.Content

	database.Db.Save(&post)

	if database.Db.Error != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot save post",
		})
		return
	}

	c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Post updated successfully",
		"post":    post,
	})
}

func DeletePost(c *fiber.Ctx) {
	post := models.Post{}
	user := models.User{}

	database.Db.First(&user, "user_id = ?", c.Params("userid"))

	if user.UserID == 0 {
		c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
		return
	}

	database.Db.First(&post, "post_id = ?", c.Params("postid"))

	if post.PostID == 0 {
		c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Post not found",
		})
		return
	}

	if post.ImageURL != nil {
		cld, err := utils.ConnectToCloudinary()
		if err != nil {
			c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
			return
		}

		lastPart := path.Base(*post.ImageURL)

		publicID := strings.Split(lastPart, ".")[0]

		if err := utils.DeleteFromCloudinary(cld, publicID); err != nil {
			c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
			return
		}
	}

	database.Db.Delete(&post)

	if database.Db.Error != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot delete post",
		})
		return
	}

	c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Post deleted successfully",
	})
}

func GetPosts(c *fiber.Ctx) {
	post := models.Post{}
	user := models.User{}

	database.Db.Find(&user, "user_id = ?", c.Params("userid"))

	if user.UserID == 0 {
		c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
		return
	}

	database.Db.Find(&post, "post_id = ?", c.Params("postid"))

	if post.PostID == 0 {
		c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Post not found",
		})
		return
	}

	c.Status(fiber.StatusOK).JSON(post)
}
