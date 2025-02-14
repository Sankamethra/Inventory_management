package main

import (
	"log"

	"order-inventory/config"
	"order-inventory/handlers"
	"order-inventory/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func setupRoutes(app *fiber.App) {
	// Initial admin setup (one-time use)
	app.Post("/api/setup/admin", handlers.SetupAdmin)

	// Public routes
	app.Post("/api/signup", handlers.Signup)
	app.Post("/api/login", handlers.Login)

	// Protected routes
	api := app.Group("/api", middleware.AuthMiddleware())
	{
		// User routes (both users and admins can access)
		api.Get("/products", handlers.GetProducts)
		api.Get("/products/:id", handlers.GetProduct)
		api.Get("/products/:id/price-history", handlers.GetPriceHistory)
		api.Post("/orders", handlers.CreateOrder)
		api.Get("/orders", handlers.GetUserOrders)
		api.Get("/orders/:id", handlers.GetUserOrderDetails)
		api.Get("/dashboard", handlers.GetUserDashboardStats)

		// Admin only routes
		admin := api.Group("/admin", middleware.AdminMiddleware())
		{
			admin.Post("/products", handlers.CreateProduct)
			admin.Put("/products/:id", handlers.UpdateProduct)
			admin.Put("/products/:id/stock", handlers.UpdateStock)
			admin.Delete("/products/:id", handlers.DeleteProduct)
			admin.Get("/orders", handlers.GetAllOrders)
			admin.Get("/stats", handlers.GetSystemStats)
		}
	}
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
