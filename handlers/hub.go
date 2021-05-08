package handlers

import (
	"github.com/NikSchaefer/go-fiber/database"
	"github.com/NikSchaefer/go-fiber/model"
	"github.com/gofiber/fiber/v2"
	guuid "github.com/google/uuid"
	"gorm.io/gorm"
)

type Hub model.Hub

func CreateScript(c *fiber.Ctx) error {
	db := database.DB
	
	json := new(Hub)
	if err := c.BodyParser(json); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	new := Hub{
		ID:       guuid.New(),
		Name: 	  json.Name,
		Value: 	  json.Value,
	}

	found := Hub{}
	query := Hub{Value: json.Value}
	err := db.First(&found, &query).Error

	if err != gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusBadRequest).SendString("Script Already Exists")
	}

	err = db.Create(&new).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Creation Error")
	}

	return c.Status(fiber.StatusOK).JSON(new)
}

func GetScripts(c *fiber.Ctx) error {
	db := database.DB
	Scripts := []Hub{}
	db.Model(&model.Hub{}).Order("ID asc").Limit(100).Find(&Scripts)
	return c.Status(fiber.StatusOK).JSON(Scripts)
}

func DeleteScript(c *fiber.Ctx) error {
	db := database.DB
	json := new(Hub)

	if err := c.BodyParser(json); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	found := Hub{}
	query := Hub{Name: json.Name}
	err := db.First(&found, &query).Error
	if err == gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusNotFound).SendString("Key not Found")
	}

	db.Delete(&found)
	return c.SendStatus(fiber.StatusOK)
}
