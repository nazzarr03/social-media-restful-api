package middleware

import (
	"github.com/gofiber/fiber"
	"strings"
)

func Authentication(secret string) fiber.Handler {
	return func(c *fiber.Ctx) {
		authHeader := c.Get("Authorization")
		t := strings.Split(authHeader, " ")
		if len(t) == 2 {
			authToken := t[1]
			authorized, _ := IsAuthorized(authToken, secret)

			if authorized {
				userID, err := ExtractIDFromToken(authToken, secret)

				if err != nil {
					c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
						"error": "Cannot extract user id from token",
					})
					c.Next()
					return
				}

				c.Set("user-id", userID)
				c.Next()
				return
			}

			c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Not authorized",
			})
			c.Next()
			return
		}

		c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Not authorized",
		})
		c.Next()
	}
}
