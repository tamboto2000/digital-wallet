# Digital Wallet API

A RESTful API for managing digital wallet transactions built with Go and PostgreSQL.

## Features

- ✅ Balance inquiry API
- ✅ Withdrawal API with transaction history
- ✅ Concurrent transaction handling with database locks
- ✅ Input validation and error handling
- ✅ Transaction atomicity using database transactions
- ✅ RESTful API design

## Tech Stack

- **Language**: Go 1.21+
- **Database**: PostgreSQL
- **Router**: Gorilla Mux
- **Database Driver**: lib/pq

## Project Structure

```
digital-wallet/
├── main.go                 # Application entry point
├── database/
│   └── db.go              # Database connection setup
├── models/
│   └── wallet.go          # Wallet business logic
├── handlers/
│   └── wallet_handlers.go # API handlers
├── schema.sql             # Database schema
├── go.mod                 # Go dependencies
├── .env.example           # Environment variables template
└── README.md              # This file
```

## Prerequisites

- Go 1.21 or higher
- PostgreSQL 12 or higher
- Git

## Installation & Setup

### 1. Clone the repository

```bash
git clone <your-repo-url>
cd digital-wallet
```

### 2. Install dependencies

```bash
go mod download
```

### 3. Setup PostgreSQL database

```bash
# Login to PostgreSQL
psql -U postgres

# Create database
CREATE DATABASE digital_wallet;

# Exit psql
\q

# Run schema
psql -U postgres -d digital_wallet -f schema.sql
```

### 4. Configure environment variables

```bash
# Copy example env file
cp .env.example .env

# Edit .env with your database credentials
# Example:
# DB_HOST=localhost
# DB_PORT=5432
# DB_USER=postgres
# DB_PASSWORD=your_password
# DB_NAME=digital_wallet
# SERVER_PORT=8080
```

### 5. Run the application

```bash
go run main.go
```

The server will start on `http://localhost:8080`

## API Documentation

### Base URL
```
http://localhost:8080
```

### Endpoints

#### 1. Health Check
Check if the service is running.

**Endpoint:** `GET /health`

**Response:**
```json
{
  "status": "healthy"
}
```

---

#### 2. Get Balance
Retrieve the current balance of a user's wallet.

**Endpoint:** `GET /api/balance/{user_id}`

**Parameters:**
- `user_id` (path parameter): User ID (integer)

**Success Response (200 OK):**
```json
{
  "user_id": 1,
  "balance": 1000.00
}
```

**Error Responses:**

- `400 Bad Request`: Invalid user ID
```json
{
  "error": "Invalid user ID"
}
```

- `404 Not Found`: Wallet not found
```json
{
  "error": "Wallet not found"
}
```

**Example:**
```bash
curl http://localhost:8080/api/balance/1
```

---

#### 3. Withdraw
Withdraw funds from a user's wallet.

**Endpoint:** `POST /api/withdraw/{user_id}`

**Parameters:**
- `user_id` (path parameter): User ID (integer)

**Request Body:**
```json
{
  "amount": 100.50
}
```

**Success Response (200 OK):**
```json
{
  "user_id": 1,
  "amount": 100.50,
  "balance_before": 1000.00,
  "balance_after": 899.50,
  "transaction_id": 1
}
```

**Error Responses:**

- `400 Bad Request`: Invalid user ID
```json
{
  "error": "Invalid user ID"
}
```

- `400 Bad Request`: Invalid request body
```json
{
  "error": "Invalid request body"
}
```

- `400 Bad Request`: Invalid amount (zero or negative)
```json
{
  "error": "Amount must be greater than 0"
}
```

- `400 Bad Request`: Insufficient funds
```json
{
  "error": "Insufficient funds"
}
```

- `404 Not Found`: Wallet not found
```json
{
  "error": "Wallet not found"
}
```

**Example:**
```bash
curl -X POST http://localhost:8080/api/withdraw/1 \
  -H "Content-Type: application/json" \
  -d '{"amount": 100.50}'
```

## Testing the API

### Using curl

**Check health:**
```bash
curl http://localhost:8080/health
```

**Get balance:**
```bash
curl http://localhost:8080/api/balance/1
```

**Withdraw funds:**
```bash
curl -X POST http://localhost:8080/api/withdraw/1 \
  -H "Content-Type: application/json" \
  -d '{"amount": 50.00}'
```

### Using Postman

1. Import the following collection or create requests manually:
   - GET `http://localhost:8080/health`
   - GET `http://localhost:8080/api/balance/1`
   - POST `http://localhost:8080/api/withdraw/1` with JSON body: `{"amount": 50.00}`

## Database Schema

### Wallets Table
```sql
CREATE TABLE wallets (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL UNIQUE,
    balance DECIMAL(15, 2) NOT NULL DEFAULT 0.00,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT positive_balance CHECK (balance >= 0)
);
```

### Transactions Table
```sql
CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    wallet_id BIGINT NOT NULL,
    amount DECIMAL(15, 2) NOT NULL,
    type VARCHAR(10) NOT NULL CHECK (type IN ('withdraw', 'deposit')),
    balance_before DECIMAL(15, 2) NOT NULL,
    balance_after DECIMAL(15, 2) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (wallet_id) REFERENCES wallets(id) ON DELETE CASCADE
);
```

## Key Features Implemented

### 1. Concurrency Safety
- Uses database row-level locking (`SELECT ... FOR UPDATE`) to prevent race conditions
- Ensures atomic updates with database transactions

### 2. Transaction Management
- All balance changes are recorded in the transactions table
- Full audit trail of all withdrawals
- Atomic operations (either both wallet update and transaction record succeed, or both fail)

### 3. Error Handling
- Comprehensive validation for user inputs
- Proper HTTP status codes
- Clear error messages
- Database error handling

### 4. Security Considerations
- Input validation to prevent SQL injection (using parameterized queries)
- Balance constraints at database level
- Transaction integrity with ACID properties

## Sample Data

The schema includes 3 test users:
- User ID 1: Balance $1,000.00
- User ID 2: Balance $500.00
- User ID 3: Balance $2,500.50

## Future Enhancements

- [ ] User authentication and authorization
- [ ] Deposit API
- [ ] Transfer between wallets
- [ ] Transaction history API
- [ ] Pagination for transaction lists
- [ ] Rate limiting
- [ ] API documentation with Swagger
- [ ] Unit and integration tests
- [ ] Docker containerization
- [ ] CI/CD pipeline

## Development

### Running tests
```bash
go test ./...
```

### Building the application
```bash
go build -o digital-wallet
./digital-wallet
```

## License

This project is for demonstration purposes.

## Author

Created as part of a technical assessment.
