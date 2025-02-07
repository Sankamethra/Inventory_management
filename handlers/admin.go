package handlers

import (
	"github.com/gofiber/fiber/v2"
	"order-inventory/config"
	"order-inventory/models"
)

// GetAllOrders retrieves all orders with filters
func GetAllOrders(c *fiber.Ctx) error {
	var orders []models.Order
	query := config.DB.Preload("User").Preload("OrderItems").Preload("OrderItems.Product")

	// Apply filters
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if userID := c.Query("user_id"); userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	// Apply sorting
	sortBy := c.Query("sort_by", "created_at")
	sortOrder := c.Query("sort_order", "desc")
	query = query.Order(sortBy + " " + sortOrder)

	result := query.Find(&orders)
	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Could not fetch orders",
		})
	}

	return c.JSON(orders)
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