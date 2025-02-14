# Order and Inventory Management System

A backend system in Go (Fiber) that manages orders and inventory with dynamic pricing based on demand and stock levels.

## Features
- User and Admin authentication
- Product management with dynamic pricing
- Order processing and tracking
- Inventory management
- Admin dashboard with statistics
- Dynamic pricing based on stock and demand

## Tech Stack
- Go Fiber framework
- PostgreSQL with GORM
- JWT Authentication

## Project Setup

### 1. Prerequisites
- Go (version 1.19 or later)
- PostgreSQL (version 13 or later)
- Git

### 2. Installation
```bash
# Clone the repository
git clone https://github.com/yourusername/order-inventory-management.git
cd order-inventory-management

# Install dependencies
go mod download
go mod tidy
```

### 3. Database Setup
```bash
# Create PostgreSQL database
psql -U postgres
CREATE DATABASE order_management;
\q

# Database migrations will run automatically when you start the server
```

### 4. Environment Configuration
Create a `.env` file in the root directory:
```env
# Database Configuration
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=postgress
DB_NAME=order_management
DB_PORT=5432

# Authentication
JWT_SECRET=orderinventorymanagement
SETUP_SECRET=orderinventorymanagement  # For admin creation
```

### 5. Run the Server
```bash
# Start the server
go run main.go

# Server will start on http://localhost:3000
```

## API Usage

### 1. Initial Admin Setup
```bash
# Create the first admin account
curl -X POST http://localhost:3000/api/setup/admin \
  -H "Content-Type: application/json" \
  -H "Setup-Secret: orderinventorymanagement" \
  -d '{
    "email": "admin@example.com",
    "password": "adminpass123",
    "first_name": "Admin",
    "last_name": "User"
  }'
```

### 2. Authentication

#### User Registration
```bash
curl -X POST http://localhost:3000/api/signup \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123",
    "first_name": "John",
    "last_name": "Doe"
  }'
```

#### Login (Users & Admin)
```bash
curl -X POST http://localhost:3000/api/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'

# Save the token from the response for future requests
```

### 3. Product Management (Admin Only)

```bash
# Create product
curl -X POST http://localhost:3000/api/admin/products \
  -H "Authorization: Bearer ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Product Name",
    "description": "Description",
    "price": 100.00,
    "base_price": 90.00,
    "stock": 50
  }'

# Update stock
curl -X PUT http://localhost:3000/api/admin/products/1/stock \
  -H "Authorization: Bearer ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "stock": 40
  }'
```

### 4. Order Management

#### For Users
```bash
# View products
curl -X GET http://localhost:3000/api/products \
  -H "Authorization: Bearer USER_TOKEN"

# Place order
curl -X POST http://localhost:3000/api/orders \
  -H "Authorization: Bearer USER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "items": [
      {
        "product_id": 1,
        "quantity": 2
      }
    ]
  }'

# View orders
curl -X GET http://localhost:3000/api/orders \
  -H "Authorization: Bearer USER_TOKEN"
```

#### For Admin
```bash
# View all orders with filters
curl -X GET "http://localhost:3000/api/admin/orders?status=pending&sort_by=created_at" \
  -H "Authorization: Bearer ADMIN_TOKEN"

# View system statistics
curl -X GET http://localhost:3000/api/admin/stats \
  -H "Authorization: Bearer ADMIN_TOKEN"
```

### 5. Dynamic Pricing System

The system automatically adjusts product prices based on:

1. Stock Levels
   - ≤5 items: +30% price increase
   - ≤20 items: +20% price increase
   - ≤50 items: +10% price increase

2. Demand (Orders per week)
   - ≥100 orders: +25% price increase
   - ≥50 orders: +15% price increase
   - ≥20 orders: +10% price increase

### Available Query Parameters

#### Product Listing
```
GET /api/products
- search: Search by product name
- min_price: Minimum price filter
- max_price: Maximum price filter
- sort_by: name, price, created_at
- sort_order: asc, desc
```

#### Order Listing (Admin)
```
GET /api/admin/orders
- status: pending, completed, cancelled
- sort_by: created_at, total_price
- sort_order: asc, desc
```

## Troubleshooting

1. Database Connection Issues:
   - Verify PostgreSQL is running
   - Check database credentials in .env
   - Ensure database exists

2. Authentication Issues:
   - Verify JWT_SECRET in .env
   - Check token expiration
   - Ensure proper token format in Authorization header