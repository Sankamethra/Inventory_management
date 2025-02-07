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

## Testing API Endpoints

### Authentication

```bash
# Signup
curl -X POST http://localhost:3000/api/signup \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123",
    "first_name": "John",
    "last_name": "Doe"
  }'

# Login
curl -X POST http://localhost:3000/api/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'

```

### Products

```bash
# Create Product
curl -X POST http://localhost:3000/api/products \
  -H "Authorization: Bearer <your-token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Product",
    "description": "Test Description",
    "price": 99.99,
    "base_price": 89.99,
    "stock": 100
  }'
```

```bash
# Get Products
curl -X GET http://localhost:3000/api/products \
  -H "Authorization: Bearer <your-token>"
```

### Orders

```bash
# Create Order
curl -X POST http://localhost:3000/api/orders \
  -H "Authorization: Bearer <your-token>" \
  -H "Content-Type: application/json" \
  -d '{
    "items": [
      {
        "product_id": 1,
        "quantity": 2
      }
    ]
  }'

```
```bash
# Get User Orders
curl -X GET http://localhost:3000/api/orders \
  -H "Authorization: Bearer <your-token>"
```

### Admin Routes

```bash
# Get System Stats (Admin only)
curl -X GET http://localhost:3000/api/admin/stats \
  -H "Authorization: Bearer <admin-token>"
```
```bash
# Get All Orders (Admin only)
curl -X GET http://localhost:3000/api/admin/orders \
  -H "Authorization: Bearer <admin-token>"
```