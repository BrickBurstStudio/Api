package handlers

import (
	"github.com/gofiber/fiber/v2"
)



func GetFile(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"code":    200,
		"message": "success",
		"data":    "found",
	})
}
