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

	found := Key{}
	query := Key{IP: ip}
	err := db.First(&found, &query).Error

	if err != gorm.ErrRecordNotFound {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "A key with that ip address has already been created",
		})
	}

	new := Key{
		IP:      ip,
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
	return c.JSON(fiber.Map{
		"code":    200,
		"message": "success",
		"data":    new,
	})
}

// func GetKeys(c *fiber.Ctx) error {
// 	db := database.DB
// 	Keys := []Key{}
// 	db.Model(&model.Key{}).Order("ID asc").Limit(100).Find(&Keys)
// 	return c.Status(fiber.StatusOK).JSON(Keys)
// }

func GetKeyById(c *fiber.Ctx) error {
	ip := c.IP()

	db := database.DB

	found := Key{}
	query := Key{IP: ip}
	err := db.First(&found, &query).Error

	if err == gorm.ErrRecordNotFound {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "Key not found",
		})
	}

	return c.JSON(fiber.Map{
		"code":    200,
		"message": "success",
		"data":    found,
	})
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
	return c.JSON(fiber.Map{
		"code":    200,
		"message": "success",
	})
}

func UpdateKey(c *fiber.Ctx) error {
	json := new(Key)

	if err := c.BodyParser(json); err != nil {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "Bad request",
		})
	}

	ip := c.IP()
	db := database.DB

	found := Key{}
	query := Key{
		IP: ip,
	}
	err := db.First(&found, &query).Error
	if err == gorm.ErrRecordNotFound {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "Product Not Found",
		})
	}
	json.IP = ip
	if json.Check1 != false {
		found.Check1 = json.Check1
	} else if json.Check2 != false {
		found.Check2 = json.Check2
	} else if json.Check3 != false {
		found.Check3 = json.Check3
	} else if json.Check4 != false {
		found.Check4 = json.Check4
	} else if json.Check5 != false {
		found.Check5 = json.Check5
	}

	db.Save(&found)
	return c.JSON(fiber.Map{
		"code":    200,
		"message": "success",
	})
}
