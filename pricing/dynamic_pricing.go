package pricing

import (
	"math"
	"order-inventory/config"
	"order-inventory/models"
	"time"
)

// PriceAdjustment handles dynamic price calculations
func CalculateNewPrice(productID uint) (float64, error) {
	var product models.Product
	if err := config.DB.First(&product, productID).Error; err != nil {
		return 0, err
	}

	// Calculate demand factor based on recent orders (last 7 days)
	var orderCount int64
	config.DB.Model(&models.OrderItem{}).
		Joins("JOIN orders ON orders.id = order_items.order_id").
		Where("order_items.product_id = ? AND orders.created_at >= ?", 
			productID, time.Now().AddDate(0, 0, -7)).
		Count(&orderCount)

	// Calculate price factors
	stockFactor := calculateStockFactor(product.Stock)
	demandFactor := calculateDemandFactor(orderCount)

	// Calculate new price
	newPrice := product.BasePrice * stockFactor * demandFactor
	
	// Round to 2 decimal places
	newPrice = math.Round(newPrice*100) / 100

	// Create price history record
	priceHistory := models.PriceHistory{
		ProductID: productID,
		OldPrice:  product.Price,
		NewPrice:  newPrice,
		Reason:    generatePriceChangeReason(stockFactor, demandFactor),
	}
	
	// Update product price and save price history
	tx := config.DB.Begin()
	if err := tx.Create(&priceHistory).Error; err != nil {
		tx.Rollback()
		return 0, err
	}
	
	if err := tx.Model(&product).Update("price", newPrice).Error; err != nil {
		tx.Rollback()
		return 0, err
	}
	
	tx.Commit()
	return newPrice, nil
}

func calculateStockFactor(stock int) float64 {
	switch {
	case stock <= 5:
		return 1.3 // 30% markup for very low stock
	case stock <= 20:
		return 1.2 // 20% markup for low stock
	case stock <= 50:
		return 1.1 // 10% markup for medium stock
	default:
		return 1.0 // No markup for normal stock levels
	}
}

func calculateDemandFactor(orderCount int64) float64 {
	switch {
	case orderCount >= 100:
		return 1.25 // 25% markup for very high demand
	case orderCount >= 50:
		return 1.15 // 15% markup for high demand
	case orderCount >= 20:
		return 1.1 // 10% markup for medium demand
	default:
		return 1.0 // No markup for low demand
	}
}

func generatePriceChangeReason(stockFactor, demandFactor float64) string {
	var reason string
	
	if stockFactor > 1.0 {
		reason += "Stock level adjustment "
	}
	if demandFactor > 1.0 {
		if len(reason) > 0 {
			reason += "and "
		}
		reason += "Demand adjustment"
	}
	
	if len(reason) == 0 {
		reason = "Regular price update"
	}
	
	return reason
}