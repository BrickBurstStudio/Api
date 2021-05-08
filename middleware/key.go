package middleware

import (
	"os"

	"github.com/gofiber/fiber/v2"
)

func KeyCheck(c *fiber.Ctx) error {
	type KeyChecker struct {
		Key string `json:"key"`
	}
	json := new(KeyChecker)
	err := c.BodyParser(json)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	if json.Key != os.Getenv("KEY") {
		return c.Status(fiber.StatusBadRequest).SendString("Incorrect Key")
	}
	
	return c.Next()
}
