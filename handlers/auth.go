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
		return c.SendStatus(fiber.StatusBadRequest)
	}

	found := User{}
	query := User{Username: json.Username}
	err := db.First(&found, &query).Error
	if err == gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusNotFound).SendString("User not Found")
	}
	if !comparePasswords(found.Password, []byte(json.Password)) {
		return c.Status(fiber.StatusBadRequest).SendString("Incorrect Password")
	}
	session := Session{UserRefer: found.ID, Expires: SessionExpires(7), SessionID: guuid.New()}
	db.Create(&session)
	c.Cookie(&fiber.Cookie{
		Name:     "sessionID",
		Expires:  SessionExpires(7),
		Value:    session.SessionID.String(),
		HTTPOnly: false,
	})
	return c.Status(fiber.StatusOK).JSON(session)
}

func Logout(c *fiber.Ctx) error {
	db := database.DB
	json := new(Session)
	if err := c.BodyParser(json); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	session := Session{}
	query := Session{SessionID: json.SessionID}
	err := db.First(&session, &query).Error
	if err == gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusBadRequest).SendString("Session Not Found")
	}
	db.Delete(&session)
	c.ClearCookie("sessionID")
	return c.SendStatus(fiber.StatusOK)
}

func CreateUser(c *fiber.Ctx) error {
	db := database.DB
	
	json := new(User)
	if err := c.BodyParser(json); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	password := hashAndSalt([]byte(json.Password))

	new := User{
		Username: strings.ToLower(json.Username),
		Discord: json.Discord,
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
	c.Cookie(&fiber.Cookie{
		Name:     "sessionID",
		Expires:  SessionExpires(7),
		Value:    session.SessionID.String(),
		HTTPOnly: false,
	})
	return c.Status(fiber.StatusOK).JSON(session)
}

func GetUserInfo(c *fiber.Ctx) error {
	user := c.Locals("user").(User)
	return c.Status(200).JSON(user)
}

func DeleteUser(c *fiber.Ctx) error {
	type DeleteUserRequest struct {
		password string
	}
	db := database.DB
	json := new(DeleteUserRequest)
	user := c.Locals("user").(User)
	if err := c.BodyParser(json); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	if !comparePasswords(user.Password, []byte(json.password)) {
		return c.Status(fiber.StatusUnauthorized).SendString("Invalid Password")
	}
	db.Model(&user).Association("Sessions").Delete()
	db.Model(&user).Association("Products").Delete()
	db.Delete(&user)
	c.ClearCookie("sessionID")
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
		return c.SendStatus(fiber.StatusBadRequest)
	}
	if !comparePasswords(user.Password, []byte(json.NewPassword)) {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid Password")
	}
	user.Password = hashAndSalt([]byte(json.NewPassword))
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
