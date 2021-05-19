package handlers

import (
	"strings"
	"time"

	"github.com/NikSchaefer/go-fiber/database"
	"github.com/NikSchaefer/go-fiber/model"
	"github.com/gofiber/fiber/v2"
	guuid "github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User model.User
type Session model.Session
type Product model.Product

func GetUser(session guuid.UUID) (User, int) {
	db := database.DB
	query := Session{SessionID: session}
	found := Session{}
	err := db.First(&found, &query).Error
	if err == gorm.ErrRecordNotFound {
		return User{}, fiber.StatusNotFound
	}
	user := User{}
	usrQuery := User{ID: found.UserRefer}
	err = db.First(&user, &usrQuery).Error
	if err == gorm.ErrRecordNotFound {
		return User{}, fiber.StatusNotFound
	}
	return user, 0
}

func Login(c *fiber.Ctx) error {
	type LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	db := database.DB
	json := new(LoginRequest)
	if err := c.BodyParser(json); err != nil {
		return c.JSON(fiber.Map{
			"code": 400,
			"message": "Bad request",
		})
	}

	found := User{}
	query := User{Username: strings.ToLower(json.Username)}
	err := db.First(&found, &query).Error
	if err == gorm.ErrRecordNotFound {
		return c.JSON(fiber.Map{
			"code": 400,
			"message": "Invalid Username",
		})
	}
	if !comparePasswords(found.Password, []byte(json.Password)) {
		return c.JSON(fiber.Map{
			"code": 400,
			"message": "Incorrect Password",
		})
	}
	session := Session{UserRefer: found.ID, Expires: SessionExpires(7), SessionID: guuid.New()}
	db.Create(&session)
	return c.JSON(fiber.Map{
		"code": 200,
		"message": "success",
		"data": session,
	})
}

func Logout(c *fiber.Ctx) error {
	db := database.DB
	json := new(Session)
	if err := c.BodyParser(json); err != nil {
		return c.JSON(fiber.Map{
			"code": 400,
			"message": "Bad request",
		})
	}
	session := Session{}
	query := Session{SessionID: json.SessionID}
	err := db.First(&session, &query).Error
	if err == gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusBadRequest).SendString("Session Not Found")
	}
	db.Delete(&session)
	c.ClearCookie("sessionid")
	return c.JSON(fiber.Map{
		"code": 200,
		"message": "success",
	})
}

func CreateUser(c *fiber.Ctx) error {
	db := database.DB
	
	json := new(User)
	print(json)
	if err := c.BodyParser(json); err != nil {
		return c.JSON(fiber.Map{/*  */
			"code": 400,/*  */
			"message": "Bad /*  */request",
			"debug":json,
		})
	}
	password := hashAndSalt([]byte(json.Password))

	new := User{
		Username: strings.ToLower(json.Username),
		Discord:  json.Discord,
		Password: password,
		Email:    json.Email,
		ID:       guuid.New(),
	}
	found := User{}
	query := User{Username: strings.ToLower(json.Username)}
	err := db.First(&found, &query).Error
	if err != gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusBadRequest).SendString("User Already Exists")
	}
	db.Create(&new)
	session := Session{UserRefer: new.ID, SessionID: guuid.New()}
	err = db.Create(&session).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Creation Error")
	}
	return c.JSON(fiber.Map{
		"code": 200,
		"message": "success",
		"data": session,
	})
}

func GetUserInfo(c *fiber.Ctx) error {
	user := c.Locals("user").(User)
	return c.JSON(fiber.Map{
		"code": 200,
		"message": "success",
		"data": user,
	})
}

func DeleteUser(c *fiber.Ctx) error {
	type DeleteUserRequest struct {
		password string
	}
	db := database.DB
	json := new(DeleteUserRequest)
	user := c.Locals("user").(User)
	if err := c.BodyParser(json); err != nil {
		return c.JSON(fiber.Map{
			"code": 400,
			"message": "Bad request",
		})
	}
	if !comparePasswords(user.Password, []byte(json.password)) {
		return c.Status(fiber.StatusUnauthorized).SendString("Invalid Password")
	}
	db.Model(&user).Association("Sessions").Delete()
	db.Model(&user).Association("Products").Delete()
	db.Delete(&user)
	c.ClearCookie("sessionid")
	return c.SendStatus(fiber.StatusOK)
}

func ChangePassword(c *fiber.Ctx) error {
	type ChangePasswordRequest struct {
		NewPassword string `json:"newPassword"`
	}
	db := database.DB
	user := c.Locals("user").(User)
	json := new(ChangePasswordRequest)
	if err := c.BodyParser(json); err != nil {
		return c.JSON(fiber.Map{
			"code": 400,
			"message": "Bad request",
		})
	}
	if !comparePasswords(user.Password, []byte(json.NewPassword)) {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid Password")
	}
	user.Password = hashAndSalt([]byte(json.NewPassword))
	db.Save(&user)
	return c.SendStatus(fiber.StatusOK)
}

func ChangeDiscord(c *fiber.Ctx) error {
	type ChangeDiscordRequest struct {
		NewDiscord int `json:"newDiscord"`
	}

	db := database.DB
	user := c.Locals("user").(User)
	json := new(ChangeDiscordRequest)
	if err := c.BodyParser(json); err != nil {
		return c.JSON(fiber.Map{
			"code": 400,
			"message": "Bad request",
		})
	}
	user.Discord = json.NewDiscord
	db.Save(&user)
	return c.SendStatus(fiber.StatusOK)
}

func hashAndSalt(pwd []byte) string {
	hash, _ := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	return string(hash)
}

func comparePasswords(hashedPwd string, plainPwd []byte) bool {
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	return err == nil
}

// Universal date the Session Will Expire
func SessionExpires(days time.Duration) time.Time {
	return time.Now().Add(days * 24 * time.Hour)
}
