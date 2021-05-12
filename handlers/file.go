package handlers

import (
	"github.com/NikSchaefer/go-fiber/database"
	"github.com/NikSchaefer/go-fiber/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type File from model.File

func GetFile(c *fiber.Ctx) error {
	json := new(File)

	if err := c.BodyParser(json); err != nil {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "Bad request",
		})
	}

	db := database.DB

	found := File{}
	query := File{Url: json.Url}
	err := db.First(&found, &query).Error

	return c.JSON(fiber.Map{
		"code":    200,
		"message": "success",
		"data":    found,
	})
}

func UpdateFile(c *fiber.Ctx) error {
	json := new(File)
	if err := c.BodyParser(json); err != nil {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "Bad request",
		})
	}

	db := database.DB

	found := File{}
	query := File{Url: json.Url}
	err := db.First(&found, &query).Error

	if err == gorm.ErrRecordNotFound {
		new := File{
			Url: json.Url,
		}
		db.Create(&new)
		return c.JSON(fiber.Map{
			"code":    200,
			"message": "success",
			"data":    new,
		})
	}

	found.Url = json.Url
	db.Save(&found)

	return c.JSON(fiber.Map{
		"code":    200,
		"message": "success",
		"data":    found,
	})
}