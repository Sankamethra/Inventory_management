package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"order-inventory/config"
	"order-inventory/handlers"
	"order-inventory/middleware"
)

func setupRoutes(app *fiber.App) {
	// Public routes
	app.Post("/api/login", handlers.Login)
	app.Post("/api/signup", handlers.Signup)

	// Protected routes
	api := app.Group("/api", middleware.AuthMiddleware())
	
	// Product routes
	api.Post("/products", handlers.CreateProduct)
	api.Get("/products", handlers.GetProducts)
	api.Get("/products/:id", handlers.GetProduct)
	api.Put("/products/:id", handlers.UpdateProduct)
	api.Delete("/products/:id", handlers.DeleteProduct)
	api.Put("/products/:id/stock", handlers.UpdateStock)
	api.Get("/products/:id/price-history", handlers.GetPriceHistory)

	// Order routes for users
	api.Post("/orders", handlers.CreateOrder)
	api.Get("/orders", handlers.GetUserOrders)
	api.Get("/orders/:id", handlers.GetUserOrderDetails)
	api.Get("/dashboard", handlers.GetUserDashboardStats)

	// Admin routes
	admin := api.Group("/admin", middleware.AdminMiddleware())
	admin.Get("/orders", handlers.GetAllOrders)
	admin.Get("/stats", handlers.GetSystemStats)
}

func main() {
	// Initialize Database
	config.ConnectDatabase()

	// Initialize Fiber App
	app := fiber.New()

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New())

	// Setup Routes
	setupRoutes(app)

	// Start Server
	log.Fatal(app.Listen(":3000"))
}