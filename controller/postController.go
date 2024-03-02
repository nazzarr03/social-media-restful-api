package controller

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/admin"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/gofiber/fiber"
	"github.com/nazzarr03/social-media-restful-api/database"
	"github.com/nazzarr03/social-media-restful-api/models"
)

type CreateAndEditPostRequest struct {
	Content  string  `json:"content"`
	ImageURL *string `json:"image_url"`
}

func CreatePost(c *fiber.Ctx) {
	var request CreateAndEditPostRequest

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

		post.ImageURL = &resp.SecureURL
	}

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
	var request CreateAndEditPostRequest

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

		if post.ImageURL != nil {
			lastPart := path.Base(*post.ImageURL)

			publicID := strings.Split(lastPart, ".")[0]

			_, err = cld.Admin.DeleteAssets(ctx, admin.DeleteAssetsParams{
				PublicIDs: []string{publicID},
			})

			if err != nil {
				c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Cannot delete image",
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

		resp, err := cld.Upload.Upload(ctx, tempFilePath, uploader.UploadParams{})

		if err != nil {
			c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Cannot upload image",
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

		lastPart := path.Base(*post.ImageURL)

		publicID := strings.Split(lastPart, ".")[0]

		_, err = cld.Admin.DeleteAssets(ctx, admin.DeleteAssetsParams{
			PublicIDs: []string{publicID},
		})

		if err != nil {
			c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Cannot delete image",
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