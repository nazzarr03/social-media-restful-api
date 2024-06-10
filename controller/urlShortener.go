package controller

import (
	"math/rand"
	"net/url"
	"time"

	"github.com/gofiber/fiber"
	"github.com/nazzarr03/social-media-restful-api/database"
	"github.com/nazzarr03/social-media-restful-api/models"
)

func CreateShortURL(c *fiber.Ctx) {
	longURL := c.FormValue("long_url")

	if longURL == "" {
		c.Status(400).JSON(fiber.Map{
			"error": "Long URL is required",
		})
		return
	}

	_, err := url.ParseRequestURI(longURL)

	if err != nil {
		c.Status(400).JSON(fiber.Map{
			"error": "Invalid URL",
		})
		return
	}

	shortURL := models.ShortURL{
		LongURL:  longURL,
		ShortKey: generateShortKey(),
	}

	database.Db.Create(&shortURL)

	if database.Db.Error != nil {
		c.Status(500).JSON(fiber.Map{
			"error": "Failed to create short URL",
		})
		return
	}

	c.Status(201).JSON(shortURL)
}

func generateShortKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	const keyLength = 6

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	shortKey := make([]byte, keyLength)

	for i := range shortKey {
		shortKey[i] = charset[r.Intn(len(charset))]
	}

	return string(shortKey)
}

func RedirectShortURL(c *fiber.Ctx) {
	shortURL := models.ShortURL{}

	database.Db.First(&shortURL, "short_key = ?", c.Params("shortkey"))

	if shortURL.ID == 0 {
		c.Status(404).JSON(fiber.Map{
			"error": "Short URL not found",
		})
		return
	}

	c.Redirect(shortURL.LongURL, 301)
}
