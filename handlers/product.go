package handlers

import (
	"order-inventory/config"
	"order-inventory/models"
	"order-inventory/pricing" // Make sure this package exists
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func CreateProduct(c *fiber.Ctx) error {
	var product models.Product
	if err := c.BodyParser(&product); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}

	result := config.DB.Create(&product)
	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Could not create product",
		})
	}

	return c.Status(201).JSON(product)
}

func GetProducts(c *fiber.Ctx) error {
	var products []models.Product
	result := config.DB.Find(&products)
	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Could not fetch products",
		})
	}

	return c.JSON(products)
}

func GetProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	var product models.Product
	result := config.DB.First(&product, id)
	if result.Error != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "Product not found",
		})
	}

	return c.JSON(product)
}

func UpdateProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	var product models.Product
	if err := c.BodyParser(&product); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}

	result := config.DB.Model(&models.Product{}).Where("id = ?", id).Updates(product)
	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Could not update product",
		})
	}

	return c.JSON(product)
}

func DeleteProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	result := config.DB.Delete(&models.Product{}, id)
	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Could not delete product",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Product deleted successfully",
	})
}

// UpdateStock updates product stock and triggers price recalculation
func UpdateStock(c *fiber.Ctx) error {
	id := c.Params("id")
	productID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid product ID",
		})
	}

	var stockUpdate struct {
		Stock int `json:"stock"`
	}

	if err := c.BodyParser(&stockUpdate); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}

	// Update stock - the trigger will handle price updates
	result := config.DB.Model(&models.Product{}).Where("id = ?", id).Update("stock", stockUpdate.Stock)
	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Could not update stock",
		})
	}

	// Optionally simulate the new price for the response
	simulatedPrice, err := pricing.SimulatePriceChange(uint(productID), stockUpdate.Stock)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Could not simulate price",
		})
	}

	return c.JSON(fiber.Map{
		"message":         "Stock updated successfully",
		"simulated_price": simulatedPrice,
	})
}

// GetPriceHistory retrieves price history for a product
func GetPriceHistory(c *fiber.Ctx) error {
	id := c.Params("id")
	var priceHistory []models.PriceHistory

	result := config.DB.Where("product_id = ?", id).
		Order("created_at desc").
		Find(&priceHistory)

	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Could not fetch price history",
		})
	}

	return c.JSON(priceHistory)
}
