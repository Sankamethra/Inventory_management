package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"order-inventory/config"
	"order-inventory/models"
	"order-inventory/pricing"
	"log"
)

type CreateOrderRequest struct {
	Items []OrderItemRequest `json:"items"`
}

type OrderItemRequest struct {
	ProductID uint `json:"product_id"`
	Quantity  int  `json:"quantity"`
}

// CreateOrder handles the creation of a new order
func CreateOrder(c *fiber.Ctx) error {
	var request CreateOrderRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}

	// Get user ID from JWT token
	user := c.Locals("user").(jwt.MapClaims)
	userID := uint(user["user_id"].(float64))

	// Start transaction
	tx := config.DB.Begin()

	// Create order
	order := models.Order{
		UserID: userID,
		Status: "pending",
	}

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{
			"message": "Could not create order",
		})
	}

	var totalPrice float64 = 0

	// Create order items
	for _, item := range request.Items {
		var product models.Product
		if err := tx.First(&product, item.ProductID).Error; err != nil {
			tx.Rollback()
			return c.Status(404).JSON(fiber.Map{
				"message": "Product not found",
			})
		}

		// Check stock
		if product.Stock < item.Quantity {
			tx.Rollback()
			return c.Status(400).JSON(fiber.Map{
				"message": "Insufficient stock for product: " + product.Name,
			})
		}

		// Create order item
		orderItem := models.OrderItem{
			OrderID:   order.ID,
			ProductID: product.ID,
			Quantity:  item.Quantity,
			Price:     product.Price,
		}

		if err := tx.Create(&orderItem).Error; err != nil {
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{
				"message": "Could not create order item",
			})
		}

		// Update product stock
		product.Stock -= item.Quantity
		if err := tx.Save(&product).Error; err != nil {
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{
				"message": "Could not update product stock",
			})
		}

		totalPrice += product.Price * float64(item.Quantity)
	}

	// Update order total price
	order.TotalPrice = totalPrice
	if err := tx.Save(&order).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{
			"message": "Could not update order total",
		})
	}

	// Commit transaction
	tx.Commit()

	// Return created order with items
	var completeOrder models.Order
	config.DB.Preload("OrderItems").Preload("OrderItems.Product").First(&completeOrder, order.ID)

	return c.Status(201).JSON(completeOrder)
}

// GetUserOrders retrieves all orders for the authenticated user
func GetUserOrders(c *fiber.Ctx) error {
	// Get user ID from JWT token
	user := c.Locals("user").(jwt.MapClaims)
	userID := uint(user["user_id"].(float64))

	var orders []models.Order
	result := config.DB.
		Preload("OrderItems").
		Preload("OrderItems.Product").
		Preload("User").
		Where("user_id = ?", userID).
		Order("created_at desc").
		Find(&orders)

	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Could not fetch orders",
		})
	}

	return c.JSON(orders)
}

// GetUserOrderDetails retrieves detailed information about a specific order
func GetUserOrderDetails(c *fiber.Ctx) error {
	orderID := c.Params("id")
	user := c.Locals("user").(jwt.MapClaims)
	userID := uint(user["user_id"].(float64))

	var order models.Order
	result := config.DB.
		Preload("OrderItems").
		Preload("OrderItems.Product").
		Where("id = ? AND user_id = ?", orderID, userID).
		First(&order)

	if result.Error != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "Order not found",
		})
	}

	return c.JSON(order)
}

// GetUserDashboardStats retrieves order statistics for the user dashboard
func GetUserDashboardStats(c *fiber.Ctx) error {
	user := c.Locals("user").(jwt.MapClaims)
	userID := uint(user["user_id"].(float64))

	// Get total number of orders
	var totalOrders int64
	config.DB.Model(&models.Order{}).Where("user_id = ?", userID).Count(&totalOrders)

	// Get total spent
	var totalSpent float64
	config.DB.Model(&models.Order{}).
		Where("user_id = ? AND status = ?", userID, "completed").
		Select("COALESCE(SUM(total_price), 0)").
		Row().
		Scan(&totalSpent)

	// Get recent orders
	var recentOrders []models.Order
	config.DB.
		Preload("OrderItems").
		Preload("OrderItems.Product").
		Where("user_id = ?", userID).
		Order("created_at desc").
		Limit(5).
		Find(&recentOrders)

	// Get orders by status
	var pendingOrders, completedOrders int64
	config.DB.Model(&models.Order{}).
		Where("user_id = ? AND status = ?", userID, "pending").
		Count(&pendingOrders)
	config.DB.Model(&models.Order{}).
		Where("user_id = ? AND status = ?", userID, "completed").
		Count(&completedOrders)

	return c.JSON(fiber.Map{
		"total_orders":     totalOrders,
		"total_spent":      totalSpent,
		"recent_orders":    recentOrders,
		"pending_orders":   pendingOrders,
		"completed_orders": completedOrders,
	})
}

func createOrderHandler(c *fiber.Ctx) error {
	var request CreateOrderRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Invalid request"})
	}

	// Get user ID from JWT token
	user := c.Locals("user").(jwt.MapClaims)
	userID := uint(user["user_id"].(float64))

	// Start transaction
	tx := config.DB.Begin()

	// Create order
	order := models.Order{
		UserID: userID,
		Status: "pending",
	}

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"message": "Could not create order"})
	}

	var totalPrice float64 = 0

	// Create order items
	for _, item := range request.Items {
		var product models.Product
		if err := tx.First(&product, item.ProductID).Error; err != nil {
			tx.Rollback()
			return c.Status(404).JSON(fiber.Map{"message": "Product not found"})
		}

		// Check stock
		if product.Stock < item.Quantity {
			tx.Rollback()
			return c.Status(400).JSON(fiber.Map{"message": "Insufficient stock for product: " + product.Name})
		}

		// Create order item
		orderItem := models.OrderItem{
			OrderID:   order.ID,
			ProductID: product.ID,
			Quantity:  item.Quantity,
			Price:     product.Price,
		}

		if err := tx.Create(&orderItem).Error; err != nil {
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{"message": "Could not create order item"})
		}

		// Update product stock
		product.Stock -= item.Quantity
		if err := tx.Save(&product).Error; err != nil {
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{"message": "Could not update product stock"})
		}

		totalPrice += product.Price * float64(item.Quantity)
	}

	// After successful order creation, update prices for all ordered products
	for _, item := range request.Items {
		_, err := pricing.CalculateNewPrice(item.ProductID)
		if err != nil {
			// Log the error but don't fail the order
			log.Printf("Error updating price for product %d: %v", item.ProductID, err)
		}
	}

	// Update order total price
	order.TotalPrice = totalPrice
	if err := tx.Save(&order).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"message": "Could not update order total"})
	}

	// Commit transaction
	tx.Commit()

	// Return created order with items
	var completeOrder models.Order
	config.DB.Preload("OrderItems").Preload("OrderItems.Product").First(&completeOrder, order.ID)

	return c.Status(201).JSON(completeOrder)
} 