# E-Commerce API - Golang

A complete RESTful e-commerce backend built with Go, Gin, MongoDB, and Razorpay payment integration.

## Features

- 🔐 **Authentication & Authorization** - JWT-based auth with role management
- 🛍️ **Product Management** - CRUD operations with pagination and filtering
- 🛒 **Shopping Cart** - Add, update, remove items with auto price calculation
- 📦 **Order Management** - Create orders with automatic stock management
- 💳 **Payment Integration** - Razorpay payment gateway with signature verification
- ✅ **Input Validation** - Request validation middleware
- 🔒 **Security** - Password hashing, JWT tokens, HMAC signature verification

## Tech Stack

- **Language**: Go 1.25
- **Framework**: Gin
- **Database**: MongoDB
- **Payment**: Razorpay
- **Authentication**: JWT
- **Validation**: go-playground/validator

## Project Structure

```
eccomerce-golang/
├── cmd/server/          # Server entry points
├── internal/
│   ├── config/          # Database configuration
│   ├── handlers/        # HTTP request handlers
│   │   ├── auth.go      # Login, Register
│   │   ├── cart.go      # Cart operations
│   │   ├── order.go     # Order management
│   │   ├── payment.go   # Payment processing
│   │   └── product.go   # Product CRUD
│   ├── middleware/      # Auth & validation middleware
│   ├── models/          # Data models
│   ├── routes/          # Route definitions
│   ├── utils/           # Helper functions
│   └── validation/      # Input validation schemas
├── .env                 # Environment variables
├── main.go             # Application entry point
└── seed_products.go    # Database seeding script
```

## Installation

### Prerequisites
- Go 1.25+
- MongoDB
- Razorpay account (for payments)

### Setup

1. **Clone the repository**
```bash
git clone <repository-url>
cd eccomerce-golang
```

2. **Install dependencies**
```bash
go mod download
```

3. **Configure environment variables**
Create `.env` file:
```env
PORT=8080
MONGO_URI=mongodb://localhost:27017
DB_NAME=ecommerce
JWT_SECRET=your_jwt_secret_key
RAZORPAY_KEY=your_razorpay_key
RAZORPAY_SECRET=your_razorpay_secret
RAZORPAY_WEBHOOK_SECRET=your_webhook_secret
```

4. **Run the application**
```bash
go run main.go
```

5. **Seed database (optional)**
```bash
go run seed_products.go
```

## API Endpoints

### Authentication
| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| POST | `/api/auth/register` | Register new user | No |
| POST | `/api/auth/login` | Login user | No |

### Products
| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| GET | `/api/products` | Get all products (paginated) | No |
| GET | `/api/products/:id` | Get single product | No |
| POST | `/api/products/create` | Create product | Yes |
| PUT | `/api/products/:id` | Update product | No |
| DELETE | `/api/products/:id` | Delete product | Yes |

**Query Parameters for GET /products:**
- `page` - Page number (default: 1)
- `limit` - Items per page (default: 10)
- `category` - Filter by category
- `min_price` - Minimum price
- `max_price` - Maximum price

### Cart
| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| POST | `/api/cart/add` | Add item to cart | Yes |
| GET | `/api/cart` | Get user cart | Yes |
| PUT | `/api/cart/update` | Update cart item | Yes |
| DELETE | `/api/cart/remove/:product_id` | Remove item | Yes |

### Orders
| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| POST | `/api/orders/create` | Create order | Yes |
| GET | `/api/orders` | Get user orders | Yes |
| GET | `/api/orders/:id` | Get single order | Yes |
| PUT | `/api/orders/:id/status` | Update order status | Yes |
| DELETE | `/api/orders/:id` | Cancel order | Yes |

### Payments
| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| POST | `/api/payments/:order_id` | Create payment | Yes |
| POST | `/api/payments/verify` | Verify payment | Yes |
| POST | `/api/payments/webhook` | Razorpay webhook | No |

## API Usage Examples

### 1. Register User
```bash
POST /api/auth/register
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "password123",
  "role": "customer"
}
```

### 2. Login
```bash
POST /api/auth/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
      "id": "...",
      "name": "John Doe",
      "email": "john@example.com",
      "role": "customer"
    }
  }
}
```

### 3. Get Products with Pagination
```bash
GET /api/products?page=1&limit=10&category=Electronics&min_price=100&max_price=500
```

### 4. Add to Cart
```bash
POST /api/cart/add
Authorization: Bearer <token>
Content-Type: application/json

{
  "product_id": "67a6cb28d2f1237cff176510",
  "quantity": 2
}
```

### 5. Create Order
```bash
POST /api/orders/create
Authorization: Bearer <token>
Content-Type: application/json

{
  "items": [
    {
      "product_id": "67a6cb28d2f1237cff176510",
      "quantity": 2
    }
  ]
}
```

### 6. Create Payment
```bash
POST /api/payments/67a6cb28d2f1237cff176510
Authorization: Bearer <token>
```

**Response:**
```json
{
  "gateway_order_id": "order_xxxxx",
  "amount": 199.98,
  "currency": "INR"
}
```

### 7. Verify Payment
```bash
POST /api/payments/verify
Authorization: Bearer <token>
Content-Type: application/json

{
  "order_id": "67a6cb28d2f1237cff176510",
  "payment_id": "pay_xxxxx",
  "signature": "razorpay_signature",
  "gateway_order_id": "order_xxxxx"
}
```

## Key Features Explained

### 1. **Automatic Stock Management**
When an order is created:
- Validates stock availability
- Decreases product stock automatically
- Returns error if insufficient stock

### 2. **Cart Total Calculation**
Cart automatically calculates total price by:
- Fetching current product prices from database
- Multiplying by quantities
- Updating on add/update/remove operations

### 3. **Payment Flow**
1. Create order → Get order ID
2. Create payment → Get Razorpay order ID
3. User pays via Razorpay frontend
4. Verify payment signature
5. Update order status to "confirmed"

### 4. **Pagination & Filtering**
Products endpoint supports:
- Page-based pagination
- Category filtering
- Price range filtering
- Customizable page size

### 5. **Security Features**
- Password hashing with bcrypt
- JWT token authentication
- ObjectID conversion in middleware (once)
- HMAC signature verification for payments
- Input validation on all endpoints

## Database Models

### User
```go
{
  id: ObjectID,
  name: string,
  email: string,
  password: string (hashed),
  role: string,
  created_at: DateTime,
  updated_at: DateTime
}
```

### Product
```go
{
  id: ObjectID,
  name: string,
  description: string,
  price: float64,
  stock: int,
  category: string,
  created_at: DateTime,
  updated_at: DateTime
}
```

### Cart
```go
{
  id: ObjectID,
  user_id: ObjectID,
  items: [{
    product_id: ObjectID,
    quantity: int
  }],
  total_price: float64,
  created_at: DateTime,
  updated_at: DateTime
}
```

### Order
```go
{
  id: ObjectID,
  user_id: ObjectID,
  items: [{
    product_id: ObjectID,
    quantity: int,
    price: float64
  }],
  total_price: float64,
  status: string,
  created_at: DateTime,
  updated_at: DateTime
}
```

### Payment
```go
{
  id: ObjectID,
  order_id: ObjectID,
  user_id: ObjectID,
  amount: float64,
  currency: string,
  status: string,
  gateway_id: string,
  payment_id: string,
  created_at: DateTime,
  updated_at: DateTime
}
```

## Development

### Hot Reload
Using Air for hot reload:
```bash
air
```

### Seed Database
Add 100,000 test products:
```bash
go run seed_products.go
```

## Testing with Razorpay

1. Get test API keys from [Razorpay Dashboard](https://dashboard.razorpay.com/)
2. Use test card: `4111 1111 1111 1111`
3. Any future expiry date
4. Any CVV

## Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| PORT | Server port | Yes |
| MONGO_URI | MongoDB connection string | Yes |
| DB_NAME | Database name | Yes |
| JWT_SECRET | Secret for JWT signing | Yes |
| RAZORPAY_KEY | Razorpay API key | Yes |
| RAZORPAY_SECRET | Razorpay API secret | Yes |
| RAZORPAY_WEBHOOK_SECRET | Webhook secret | No |

## Error Handling

All endpoints return consistent error format:
```json
{
  "success": false,
  "message": "Error message",
  "errors": {
    "field": "Validation error"
  }
}
```

## Contributing

1. Fork the repository
2. Create feature branch
3. Commit changes
4. Push to branch
5. Create Pull Request

## License

MIT License

## Author

Your Name

## Support

For issues and questions, please open an issue on GitHub.
