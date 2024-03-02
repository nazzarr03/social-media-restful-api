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
	"github.com/nazzarr03/social-media-restful-api/middleware"
	"github.com/nazzarr03/social-media-restful-api/models"
	"golang.org/x/crypto/bcrypt"
)

type SignUpRequest struct {
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func SignUp(c *fiber.Ctx) {
	var request SignUpRequest

	if err := c.BodyParser(&request); err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(request.Password), 14)

	user := models.User{
		Name:     request.Name,
		Surname:  request.Surname,
		Username: request.Username,
		Email:    request.Email,
		Password: string(hashedPassword),
	}

	database.Db.Create(&user)

	if database.Db.Error != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot create user",
		})
		return
	}

	c.Status(fiber.StatusCreated).JSON(user)
}

func Login(c *fiber.Ctx) {
	var request LoginRequest

	if err := c.BodyParser(&request); err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
		return
	}

	var user models.User

	database.Db.Where("username = ?", request.Username).First(&user)

	if user.UserID == 0 {
		c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(request.Password), 14)

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), hashedPassword)

	if err != nil {
		c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid password",
		})
	}

	secret := os.Getenv("JWT_SECRET")

	accessToken, err := middleware.CreateAccessToken(&user, secret, 24)

	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot create access token",
		})
		return
	}

	refreshToken, err := middleware.CreateRefreshToken(&user, secret, 168)

	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot create refresh token",
		})
		return
	}

	c.JSON(fiber.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func Logout(c *fiber.Ctx) {
	authHeader := c.Get("Authorization")
	parts := strings.Split(authHeader, " ")

	if len(parts) != 2 {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid authorization header",
		})
		return
	}

	token := parts[1]

	result := database.Db.Exec("INSERT INTO blacklists (token) VALUES (?)", token)

	if result.Error != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot blacklist token",
		})
		return
	}

	c.JSON(fiber.Map{
		"message": "Successfully logged out",
	})
}

func UploadProfilePicture(c *fiber.Ctx) {
	file, err := c.FormFile("profile_picture")

	if err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse form",
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

	user := models.User{}

	database.Db.First(&user, "user_id = ?", c.Params("id"))

	if user.UserID == 0 {
		c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
		return
	}

	user.ImageURL = &resp.SecureURL
	database.Db.Save(&user)

	if database.Db.Error != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot save image URL",
		})
		return
	}

	c.JSON(fiber.Map{
		"message":   "Profile Image uploaded successfully",
		"image_url": resp.SecureURL,
	})
}

func UpdateProfilePicture(c *fiber.Ctx) {
	user := models.User{}
	database.Db.First(&user, "user_id = ?", c.Params("id"))

	if user.UserID == 0 {
		c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
		return
	}

	lastPart := path.Base(*user.ImageURL)

	publicID := strings.Split(lastPart, ".")[0]

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

	_, err = cld.Admin.DeleteAssets(ctx, admin.DeleteAssetsParams{
		PublicIDs: []string{publicID},
	})

	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot delete image",
		})
		return
	}

	file, err := c.FormFile("profile_picture")

	if err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse form",
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

	resp, err := cld.Upload.Upload(ctx, tempFilePath, uploader.UploadParams{})

	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot upload image",
		})
		return
	}

	os.Remove(tempFilePath)

	user.ImageURL = &resp.SecureURL

	database.Db.Save(&user)

	if database.Db.Error != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot save image URL",
		})
		return
	}

	c.JSON(fiber.Map{
		"message":   "Profile Image updated successfully",
		"image_url": resp.SecureURL,
	})

}
