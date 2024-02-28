package controller

import (
	"context"
	"fmt"
	"os"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/gofiber/fiber"
	"github.com/nazzarr03/social-media-restful-api/database"
	"github.com/nazzarr03/social-media-restful-api/models"
)

type AddCommentRequest struct {
	Content      string  `json:"content"`
	ImageURL     *string `json:"image_url"`
}

func AddCommentToPost(c *fiber.Ctx) {
	var request AddCommentRequest

	user := models.User{}
	post := models.Post{}
	comment := models.Comment{}

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

	if err == nil {
		tempFilePath := fmt.Sprintf("./uploads/%s", file.Filename)

		if err := c.SaveFile(file, tempFilePath); err != nil {
			c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Cannot save file",
			})
			return
		}

		cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
		apiKey := os.Getenv("CLOUDINARY_API_KEY")
		apiSecret := os.Getenv("CLOUDINARY_API_SECRET")

		cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)

		if err != nil {
			c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Cannot connect to cloudinary",
			})
			return
		}

		var ctx = context.Background()

		resp, err := cld.Upload.Upload(ctx, tempFilePath, uploader.UploadParams{})

		if err != nil {
			c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Cannot upload image",
			})
			return
		}

		os.Remove(tempFilePath)

		comment.ImageURL = &resp.SecureURL
	}

	comment.Content = request.Content
	comment.PostID = post.PostID
	comment.UserID = user.UserID

	database.Db.Create(&comment)

	if database.Db.Error != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot create comment",
		})
		return
	}

	c.Status(fiber.StatusCreated).JSON(comment)
}
