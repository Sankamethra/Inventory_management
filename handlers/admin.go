package handlers

import (
	"order-inventory/config"
	"order-inventory/models"

	"github.com/gofiber/fiber/v2"
)

// GetAllOrders retrieves all orders with filters
func GetAllOrders(c *fiber.Ctx) error {
	// Get query parameters
	status := c.Query("status")
	sortBy := c.Query("sort_by", "created_at")
	sortOrder := c.Query("sort_order", "desc")

	// Build query
	query := config.DB.Debug().
		Preload("Items").         // Changed from OrderItems
		Preload("Items.Product"). // Changed from OrderItems.Product
		Preload("User")

	// Apply filters
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Apply sorting
	if sortOrder == "asc" {
		query = query.Order(sortBy + " asc")
	} else {
		query = query.Order(sortBy + " desc")
	}

	var orders []models.Order
	if err := query.Find(&orders).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Could not fetch orders",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   orders,
	})
}

// GetSystemStats retrieves system-wide statistics
func GetSystemStats(c *fiber.Ctx) error {
	var stats struct {
		TotalOrders     int64   `json:"total_orders"`
		TotalRevenue    float64 `json:"total_revenue"`
		TotalUsers      int64   `json:"total_users"`
		TotalProducts   int64   `json:"total_products"`
		LowStockItems   int64   `json:"low_stock_items"`
		PendingOrders   int64   `json:"pending_orders"`
		CompletedOrders int64   `json:"completed_orders"`
	}

	// Get total orders and revenue
	config.DB.Model(&models.Order{}).Count(&stats.TotalOrders)
	config.DB.Model(&models.Order{}).Select("COALESCE(SUM(total_price), 0)").Row().Scan(&stats.TotalRevenue)

	// Get user count
	config.DB.Model(&models.User{}).Count(&stats.TotalUsers)

	// Get product stats
	config.DB.Model(&models.Product{}).Count(&stats.TotalProducts)
	config.DB.Model(&models.Product{}).Where("stock < ?", 10).Count(&stats.LowStockItems)

	// Get order status counts
	config.DB.Model(&models.Order{}).Where("status = ?", "pending").Count(&stats.PendingOrders)
	config.DB.Model(&models.Order{}).Where("status = ?", "completed").Count(&stats.CompletedOrders)

	return c.JSON(stats)
}
