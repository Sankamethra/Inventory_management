package pricing

import (
	"order-inventory/config"
	"order-inventory/models"
)

// SimulatePriceChange calculates what the price would be without actually updating it
func SimulatePriceChange(productID uint, simulatedStock int) (float64, error) {
	var product models.Product
	if err := config.DB.First(&product, productID).Error; err != nil {
		return 0, err
	}

	// Get recent order count
	var orderCount int64
	config.DB.Model(&models.OrderItem{}).
		Joins("JOIN orders ON orders.id = order_items.order_id").
		Where("order_items.product_id = ? AND orders.created_at >= NOW() - INTERVAL '7 days'", 
			productID).
		Count(&orderCount)

	// Calculate factors
	stockFactor := calculateStockFactor(simulatedStock)
	demandFactor := calculateDemandFactor(orderCount)

	// Calculate simulated price
	return product.BasePrice * stockFactor * demandFactor, nil
}

// Helper functions can remain the same as they're now used for simulation
func calculateStockFactor(stock int) float64 {
	switch {
	case stock <= 5:
		return 1.3
	case stock <= 20:
		return 1.2
	case stock <= 50:
		return 1.1
	default:
		return 1.0
	}
}

func calculateDemandFactor(orderCount int64) float64 {
	switch {
	case orderCount >= 100:
		return 1.25
	case orderCount >= 50:
		return 1.15
	case orderCount >= 20:
		return 1.1
	default:
		return 1.0
	}
}