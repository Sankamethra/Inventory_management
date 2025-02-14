package handlers

import (
	"order-inventory/config"
	"order-inventory/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type CreateOrderRequest struct {
	Items []OrderItemRequest `json:"items" validate:"required,dive"`
}

type OrderItemRequest struct {
	ProductID uint `json:"product_id" validate:"required"`
	Quantity  int  `json:"quantity" validate:"required,gt=0"`
}

// CreateOrder handles the creation of a new order
func CreateOrder(c *fiber.Ctx) error {
	var request CreateOrderRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request format",
			"error":   err.Error(),
		})
	}

	// Validate request
	if len(request.Items) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Order must contain at least one item",
		})
	}

	// Get user ID from JWT token
	user := c.Locals("user").(jwt.MapClaims)
	userID := uint(user["user_id"].(float64))

	// Start transaction
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create order with initial total_price
	order := models.Order{
		UserID:     userID,
		Status:     "pending",
		TotalPrice: 0, // Initialize with 0
	}

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create order",
			"error":   err.Error(),
		})
	}

	var totalAmount float64 = 0

	// Process each order item
	for _, item := range request.Items {
		var product models.Product
		if err := tx.First(&product, item.ProductID).Error; err != nil {
			tx.Rollback()
			return c.Status(404).JSON(fiber.Map{
				"status":     "error",
				"message":    "Product not found",
				"product_id": item.ProductID,
			})
		}

		// Create order item
		orderItem := models.OrderItem{
			OrderID:   order.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     product.Price,
		}

		if err := tx.Create(&orderItem).Error; err != nil {
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to create order item",
				"error":   err.Error(),
			})
		}

		// Update product stock
		if err := tx.Model(&product).Update("stock", product.Stock-item.Quantity).Error; err != nil {
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to update product stock",
				"error":   err.Error(),
			})
		}

		totalAmount += product.Price * float64(item.Quantity)
	}

	// Update order total_price
	if err := tx.Model(&order).Update("total_price", totalAmount).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update order total",
			"error":   err.Error(),
		})
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to commit transaction",
			"error":   err.Error(),
		})
	}

	// Fetch complete order details
	var completeOrder models.Order
	if err := config.DB.Preload("Items").
		Preload("Items.Product").
		First(&completeOrder, order.ID).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch complete order details",
			"error":   err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"status":  "success",
		"message": "Order created successfully",
		"data":    completeOrder,
	})
}

// GetUserOrders retrieves all orders for the authenticated user
func GetUserOrders(c *fiber.Ctx) error {
	user := c.Locals("user").(jwt.MapClaims)
	userID := uint(user["user_id"].(float64))

	var orders []models.Order
	result := config.DB.Debug().
		Preload("Items").
		Preload("Items.Product").
		Preload("User").
		Where("user_id = ?", userID).
		Order("created_at desc").
		Find(&orders)

	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Could not fetch orders",
			"error":   result.Error.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   orders,
	})
}

// GetUserOrderDetails retrieves detailed information about a specific order
func GetUserOrderDetails(c *fiber.Ctx) error {
	orderID := c.Params("id")
	user := c.Locals("user").(jwt.MapClaims)
	userID := uint(user["user_id"].(float64))

	var order models.Order
	result := config.DB.
		Preload("Items").
		Preload("Items.Product").
		Where("id = ? AND user_id = ?", orderID, userID).
		First(&order)

	if result.Error != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Order not found",
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   order,
	})
}

// GetUserDashboardStats retrieves order statistics for the user dashboard
func GetUserDashboardStats(c *fiber.Ctx) error {
	user := c.Locals("user").(jwt.MapClaims)
	userID := uint(user["user_id"].(float64))

	var stats struct {
		TotalOrders     int64          `json:"total_orders"`
		TotalSpent      float64        `json:"total_spent"`
		RecentOrders    []models.Order `json:"recent_orders"`
		PendingOrders   int64          `json:"pending_orders"`
		CompletedOrders int64          `json:"completed_orders"`
	}

	// Get total orders
	if err := config.DB.Model(&models.Order{}).
		Where("user_id = ?", userID).
		Count(&stats.TotalOrders).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch order statistics",
			"error":   err.Error(),
		})
	}

	// Get total spent
	if err := config.DB.Model(&models.Order{}).
		Where("user_id = ? AND status = ?", userID, "completed").
		Select("COALESCE(SUM(total_price), 0)").
		Row().
		Scan(&stats.TotalSpent); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to calculate total spent",
			"error":   err.Error(),
		})
	}

	// Get recent orders
	if err := config.DB.
		Preload("Items").
		Preload("Items.Product").
		Where("user_id = ?", userID).
		Order("created_at desc").
		Limit(5).
		Find(&stats.RecentOrders).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch recent orders",
			"error":   err.Error(),
		})
	}

	// Get order status counts
	if err := config.DB.Model(&models.Order{}).
		Where("user_id = ? AND status = ?", userID, "pending").
		Count(&stats.PendingOrders).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to count pending orders",
			"error":   err.Error(),
		})
	}

	if err := config.DB.Model(&models.Order{}).
		Where("user_id = ? AND status = ?", userID, "completed").
		Count(&stats.CompletedOrders).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to count completed orders",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   stats,
	})
}
