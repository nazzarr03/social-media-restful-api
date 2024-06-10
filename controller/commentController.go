package controller

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber"
	"github.com/nazzarr03/social-media-restful-api/database"
	"github.com/nazzarr03/social-media-restful-api/dto"
	"github.com/nazzarr03/social-media-restful-api/models"
	"github.com/nazzarr03/social-media-restful-api/utils"
)

func AddCommentToPost(c *fiber.Ctx) {
	var request dto.AddCommentRequest

	user := models.User{}
	post := models.Post{}

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

	file, err := c.FormFile("comment_picture")

	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot save file",
		})
		return
	}

	if file != nil {
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

		request.ImageURL = &resp.SecureURL
	}

	comment := models.Comment{
		Content:  request.Content,
		ImageURL: request.ImageURL,
		PostID:   post.PostID,
		UserID:   user.UserID,
	}

	database.Db.Create(&comment)

	if database.Db.Error != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot create comment",
		})
		return
	}

	c.Status(fiber.StatusCreated).JSON(comment)
}

func AddCommentToComment(c *fiber.Ctx) {
	var request dto.AddCommentRequest

	user := models.User{}
	post := models.Post{}
	comment := models.Comment{}

	if err := c.BodyParser(&request); err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
		return
	}

	database.Db.First(&comment, "comment_id = ?", c.Params("commentid"))

	if comment.CommentID == 0 {
		c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Parent comment not found",
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

	file, err := c.FormFile("comment_picture")
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot save file",
		})
		return
	}

	if file != nil {
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

		request.ImageURL = &resp.SecureURL
	}

	comment.Comments = append(comment.Comments, models.Comment{
		Content:         request.Content,
		ImageURL:        request.ImageURL,
		PostID:          post.PostID,
		UserID:          user.UserID,
		ParentCommentID: &comment.CommentID,
	})

	database.Db.Create(&comment)

	database.Db.Preload("Comments", "comment_id != ?", comment.CommentID).First(&comment, "comment_id = ?", comment.CommentID)

	if database.Db.Error != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot create comment",
		})
		return
	}

	c.Status(fiber.StatusCreated).JSON(comment)
}
