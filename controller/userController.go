package controller

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/gofiber/fiber"
	"github.com/nazzarr03/social-media-restful-api/database"
	"github.com/nazzarr03/social-media-restful-api/dto"
	"github.com/nazzarr03/social-media-restful-api/middleware"
	"github.com/nazzarr03/social-media-restful-api/models"
	"github.com/nazzarr03/social-media-restful-api/utils"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(c *fiber.Ctx) {
	var request dto.SignUpRequest

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
	var request dto.LoginRequest

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

	database.Db.Exec("INSERT INTO blacklists (token) VALUES (?)", token)

	if database.Db.Error != nil {
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

	cld, err := utils.ConnectToCloudinary()
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
		return
	}

	if err := utils.DeleteFromCloudinary(cld, publicID); err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
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

	resp, err := utils.UploadToCloudinary(cld, tempFilePath)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
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
