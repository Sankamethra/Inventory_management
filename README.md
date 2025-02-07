# Order & Inventory Management System

A robust backend system built with Go (Fiber) for managing orders and inventory with dynamic pricing capabilities.

## Features

- User Authentication (JWT-based)
- Product Management
- Order Processing
- Inventory Tracking
- User Dashboard with Statistics
- Dynamic Pricing System

## Tech Stack

- **Language:** Go
- **Framework:** Fiber
- **Database:** PostgreSQL
- **ORM:** GORM
- **Authentication:** JWT

## Prerequisites

- Go 1.16 or higher
- PostgreSQL 12 or higher
- Make sure to have the following environment variables set in your `.env` file:

```env
DB_HOST=your_db_host
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=your_db_name
DB_PORT=your_db_port
JWT_SECRET=your_jwt_secret
```

## Project Structure

```
├── config/
│   └── database.go         # Database configuration
├── handlers/
│   ├── admin.go           # Admin handlers
│   ├── auth.go            # Authentication handlers
│   ├── order.go           # Order management handlers
│   └── product.go         # Product management handlers
├── middleware/
│   ├── admin.go           # Admin middleware
│   └── auth.go            # Authentication middleware
├── models/
│   ├── order.go           # Order and OrderItem models
│   ├── product.go         # Product and PriceHistory models
│   └── user.go            # User model
├── pricing/
│   └── dynamic_pricing.go # Dynamic pricing logic
├── main.go                # Application entry point
└── README.md
```

## API Endpoints

### Public Endpoints

```
POST /api/login    - User login
POST /api/signup   - User registration
```

### Protected Endpoints (Requires Authentication)

#### Products
```
POST   /api/products      - Create a new product
GET    /api/products      - Get all products
GET    /api/products/:id  - Get product by ID
PUT    /api/products/:id  - Update product
PUT    /api/products/:id/stock        - Update product stock
GET    /api/products/:id/price-history - Get product price history
GET    /api/dashboard                 - Get user dashboard stats
DELETE /api/products/:id  - Delete product
```

#### Orders
```
POST   /api/orders        - Create a new order
GET    /api/orders        - Get user's orders
GET    /api/orders/:id    - Get order details
GET    /api/dashboard     - Get user dashboard statistics
```

## Installation & Setup

1. Clone the repository:
```bash
git clone <repository-url>
```

2. Install dependencies:
```bash
go mod download
```

3. Set up your environment variables in `.env` file

# Create PostgreSQL database
```bash
createdb order_inventory
```

# Copy environment file
```bash
cp .env.example .env
```

# Update .env with your credentials

4. Run the application:
```bash
go run main.go
```

The server will start on `http://localhost:3000`

## Authentication

The system uses JWT (JSON Web Tokens) for authentication. Include the token in the Authorization header:
```
Authorization: Bearer <your-token>
```

## Database Models

### User
- Email (unique)
- Password (hashed)
- FirstName
- LastName
- Role (admin/user)

### Product
- Name
- Description
- Price
- BasePrice
- Stock

### Order
- UserID
- Status
- TotalPrice
- OrderItems

### OrderItem
- OrderID
- ProductID
- Quantity
- Price

### PriceHistory
- ProductID
- OldPrice
- NewPrice
- Reason

## Database Triggers

The system implements two main database triggers for automatic price adjustments and inventory management:

1. **Dynamic Pricing Trigger (`product_price_update`):**
   - Automatically adjusts product prices based on:
     - Current stock levels
     - Recent demand (orders in the last 7 days)
   - Creates price history records for tracking changes
   - Triggered when product stock is updated

2. **Inventory Management Trigger (`inventory_update`):**
   - Automatically updates product stock when new orders are placed
   - Ensures real-time inventory tracking
   - Triggered when new order items are created

### Price Adjustment Factors

Stock-based adjustments:
- Very low stock (≤5): +30% markup
- Low stock (≤20): +20% markup
- Medium stock (≤50): +10% markup
- Normal stock (>50): No markup

Demand-based adjustments:
- Very high demand (≥100 orders/week): +25% markup
- High demand (≥50 orders/week): +15% markup
- Medium demand (≥20 orders/week): +10% markup
- Low demand (<20 orders/week): No markup

## Testing API Endpoints

### Authentication

```bash
# 1. User Signup
curl -X POST http://localhost:3000/api/signup \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123",
    "first_name": "John",
    "last_name": "Doe"
  }'

# 2. User Login (save the returned token)
curl -X POST http://localhost:3000/api/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

### Product Management

```bash
# 1. Create Product (Admin only)
curl -X POST http://localhost:3000/api/products \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Product",
    "description": "Product Description",
    "price": 100.00,
    "base_price": 90.00,
    "stock": 50
  }'

# 2. Get All Products
curl -X GET http://localhost:3000/api/products \
  -H "Authorization: Bearer YOUR_TOKEN"

# 3. Get Single Product
curl -X GET http://localhost:3000/api/products/1 \
  -H "Authorization: Bearer YOUR_TOKEN"

# 4. Update Product Stock (triggers price adjustment)
curl -X PUT http://localhost:3000/api/products/1/stock \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "stock": 5
  }'

# 5. Get Product Price History
curl -X GET http://localhost:3000/api/products/1/price-history \
  -H "Authorization: Bearer YOUR_TOKEN"

# 6. Update Product Details
curl -X PUT http://localhost:3000/api/products/1 \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Updated Product",
    "description": "Updated Description",
    "base_price": 95.00
  }'

# 7. Delete Product (Admin only)
curl -X DELETE http://localhost:3000/api/products/1 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Order Management

```bash
# 1. Create Order
curl -X POST http://localhost:3000/api/orders \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "items": [
      {
        "product_id": 1,
        "quantity": 2
      }
    ]
  }'

# 2. Get User Orders
curl -X GET http://localhost:3000/api/orders \
  -H "Authorization: Bearer YOUR_TOKEN"

# 3. Get Single Order Details
curl -X GET http://localhost:3000/api/orders/1 \
  -H "Authorization: Bearer YOUR_TOKEN"

# 4. Get User Dashboard Stats
curl -X GET http://localhost:3000/api/dashboard \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Admin Operations

```bash
# 1. Get All Orders (with filters)
curl -X GET "http://localhost:3000/api/admin/orders?status=pending&sort_by=created_at&sort_order=desc" \
  -H "Authorization: Bearer ADMIN_TOKEN"

# 2. Get System Statistics
curl -X GET http://localhost:3000/api/admin/stats \
  -H "Authorization: Bearer ADMIN_TOKEN"

# 3. Get Low Stock Products
curl -X GET "http://localhost:3000/api/admin/products?stock_below=20" \
  -H "Authorization: Bearer ADMIN_TOKEN"
```

### Testing Dynamic Pricing

```bash
# 1. Create test product
curl -X POST http://localhost:3000/api/products \
  -H "Authorization: Bearer ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Dynamic Price Test",
    "description": "Test product for dynamic pricing",
    "price": 100.00,
    "base_price": 90.00,
    "stock": 100
  }'

# 2. Create multiple orders to test demand-based pricing
for i in {1..5}; do
  curl -X POST http://localhost:3000/api/orders \
    -H "Authorization: Bearer YOUR_TOKEN" \
    -H "Content-Type: application/json" \
    -d '{
      "items": [
        {
          "product_id": 1,
          "quantity": 2
        }
      ]
    }'
done

# 3. Update stock to test availability-based pricing
curl -X PUT http://localhost:3000/api/products/1/stock \
  -H "Authorization: Bearer ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "stock": 5
  }'

# 4. Check price history to verify dynamic pricing
curl -X GET http://localhost:3000/api/products/1/price-history \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Expected Behaviors

1. Dynamic Pricing:
   - Price increases with low stock:
     * ≤5 items: +30%
     * ≤20 items: +20%
     * ≤50 items: +10%
   - Price increases with high demand:
     * ≥100 orders/week: +25%
     * ≥50 orders/week: +15%
     * ≥20 orders/week: +10%

2. Inventory Management:
   - Stock updates automatically when orders are placed
   - Low stock triggers price adjustments
   - Price history is recorded for all changes

3. Order Processing:
   - Orders are validated against available stock
   - Multiple items can be ordered together
   - Order history is maintained