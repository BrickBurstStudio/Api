package router

import (
	"github.com/NikSchaefer/go-fiber/handlers"
	"github.com/NikSchaefer/go-fiber/middleware"
	"github.com/gofiber/fiber/v2"
)

func Initalize(router *fiber.App) {
	router.Use(middleware.Security)

	router.Get("/", func(c *fiber.Ctx) error {
		return c.Status(200).SendString("Hello, World!")
	})

	router.Use(middleware.Json)

	users := router.Group("/users")
	users.Delete("/", middleware.Authenticated, handlers.DeleteUser)
	users.Patch("/password", middleware.Authenticated, handlers.ChangePassword)
	users.Patch("/link", middleware.Authenticated, handlers.ChangeDiscord)
	users.Post("/", middleware.Authenticated, handlers.GetUserInfo)
	users.Put("/", handlers.CreateUser)
	users.Post("/login", handlers.Login)
	users.Delete("/logout", handlers.Logout)

	products := router.Group("/scripts", middleware.Authenticated)
	products.Put("/", handlers.CreateProduct)
	products.Post("/all", handlers.GetProducts)
	products.Delete("/", handlers.DeleteProduct)
	products.Post("/", handlers.GetProductById)
	products.Patch("/", handlers.UpdateProduct)

	keys := router.Group("/keys")
	keys.Put("/", handlers.CreateKey)
	keys.Post("/all", handlers.GetKeys)
	keys.Delete("/", handlers.DeleteKey)
	keys.Post("/", handlers.GetKeyById)

	// files := router.Group("/files")
	// files.Put("/", handlers.CreateFile)
	// files.Post("/", handlers.GetFile)
	// files.Delete("/", handlers.DeleteFile)
	// files.Patch("/", handlers.UpdateFile)


	router.Use(func(c *fiber.Ctx) error {
		return c.Status(404).SendString("404 Not Found")
	})

}
