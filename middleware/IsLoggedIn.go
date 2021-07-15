package middleware

import (
	"example.com/app/domain"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

func IsLoggedIn(c *fiber.Ctx) error {
	token := c.Get("Authorization")

	var auth domain.Authentication
	_, loggedIn, err := auth.IsLoggedIn(token)

	if err != nil || loggedIn == false {
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("Unauthorized user")})
	}

	err = c.Next()

	if err != nil {
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("Unauthorized user")})
	}

	return nil
}
