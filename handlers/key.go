package handlers

import (
	"github.com/NikSchaefer/go-fiber/database"
	"github.com/NikSchaefer/go-fiber/model"
	"github.com/gofiber/fiber/v2"
	guuid "github.com/google/uuid"
	"gorm.io/gorm"
)

type Key model.Key

func CreateKey(c *fiber.Ctx) error {
	db := database.DB
	ip := c.IP()

	json := new(Key)
	if err := c.BodyParser(json); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	found := Key{}
	query := Key{IP: ip}
	err := db.First(&found, &query).Error

	if err != gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusNotFound).SendString("A key with that ip address has already been created")
	}
	
	new := Key{
		IP:ip,
		ID:      guuid.New(),
		Expires: SessionExpires(1),
	}

	db.Create(&new)

	c.Cookie(&fiber.Cookie{
		Name:     "key",
		Expires:  SessionExpires(1),
		Value:    new.ID.String(),
		HTTPOnly: false,
	})
	return c.Status(fiber.StatusOK).JSON(new)
}

func GetKeys(c *fiber.Ctx) error {
	db := database.DB
	Keys := []Key{}
	db.Model(&model.Key{}).Order("ID asc").Limit(100).Find(&Keys)
	return c.Status(fiber.StatusOK).JSON(Keys)
}

func GetKeyById(c *fiber.Ctx) error {
	ip := c.IP()

	db := database.DB
	json := new(Key)

	if err := c.BodyParser(json); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	found := Key{}
	query := Key{IP: ip}
	err := db.First(&found, &query).Error

	if err == gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusNotFound).SendString("Key not Found")
	}

	db.Save(&found)

	return c.Status(fiber.StatusOK).JSON(found)
}

func DeleteKey(c *fiber.Ctx) error {
	db := database.DB
	json := new(Key)

	if err := c.BodyParser(json); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	found := Key{}
	query := Key{ID: json.ID}
	err := db.First(&found, &query).Error
	if err == gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusNotFound).SendString("Key not Found")
	}

	db.Delete(&found)
	return c.SendStatus(fiber.StatusOK)
}
