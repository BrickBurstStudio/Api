package middleware

import (
	"github.com/NikSchaefer/go-fiber/handlers"
	"github.com/NikSchaefer/go-fiber/model"
	"github.com/gofiber/fiber/v2"
)

func Authenticated(c *fiber.Ctx) error {
	json := new(model.Session)
	if err := c.BodyParser(json); err != nil {
		return c.JSON(fiber.Map{
			"code": 400,
			"message": "Bad request",

		})
	}
	user, status := handlers.GetUser(json.SessionID)
	if status != 0 {
		return c.JSON(fiber.Map{
			"code": status,
		})
	}
	c.Locals("user", user)
	return c.Next()
}
