# Quick Start Guide

Get the Digital Wallet API up and running in 5 minutes!

## Option 1: Using Docker (Recommended - Easiest)

```bash
# 1. Navigate to project directory
cd digital-wallet

# 2. Start everything with Docker Compose
docker-compose up -d

# 3. Wait for services to be ready (about 10-20 seconds)
docker-compose logs -f app

# 4. Test the API
curl http://localhost:8080/health
curl http://localhost:8080/api/balance/1
```

That's it! The API is now running on `http://localhost:8080`

**To stop:**
```bash
docker-compose down
```

---

## Option 2: Local Setup (Traditional)

### Prerequisites
- Go 1.21+ installed
- PostgreSQL 12+ running locally

### Steps

```bash
# 1. Navigate to project directory
cd digital-wallet

# 2. Install Go dependencies
go mod download

# 3. Create database
createdb -U postgres digital_wallet

# 4. Apply schema
psql -U postgres -d digital_wallet -f schema.sql

# 5. Configure environment
cp .env.example .env
# Edit .env with your database credentials

# 6. Run the application
go run main.go
```

**Test it:**
```bash
curl http://localhost:8080/health
```

---

## Testing the API

### Check Balance
```bash
curl http://localhost:8080/api/balance/1
```

Expected response:
```json
{
  "user_id": 1,
  "balance": 1000.00
}
```

### Make a Withdrawal
```bash
curl -X POST http://localhost:8080/api/withdraw/1 \
  -H "Content-Type: application/json" \
  -d '{"amount": 100.00}'
```

Expected response:
```json
{
  "user_id": 1,
  "amount": 100.00,
  "balance_before": 1000.00,
  "balance_after": 900.00,
  "transaction_id": 1
}
```

---

## Sample Test Users

The database is pre-populated with 3 test users:

| User ID | Initial Balance |
|---------|----------------|
| 1       | $1,000.00      |
| 2       | $500.00        |
| 3       | $2,500.50      |

---

## Next Steps

1. **Read the full documentation:** `README.md`
2. **Try all test scenarios:** `API_TESTING.md`
3. **Import Postman collection:** `postman_collection.json`
4. **Review the code structure**

---

## Common Commands

```bash
# Using Make (if available)
make help           # Show all available commands
make build          # Build the application
make run            # Run the application
make test           # Run tests
make docker-up      # Start with Docker
make docker-down    # Stop Docker services

# Manual commands
go build            # Build binary
go run main.go      # Run directly
go test ./...       # Run tests
```

---

## Need Help?

- Check `README.md` for detailed documentation
- See `API_TESTING.md` for comprehensive testing guide
- Review code comments for implementation details

---

## Troubleshooting

**Docker not working?**
- Ensure Docker and Docker Compose are installed
- Check if ports 8080 and 5432 are available

**Database connection failed?**
- Verify PostgreSQL is running: `systemctl status postgresql` (Linux) or `brew services list` (macOS)
- Check credentials in `.env` file
- Ensure database exists: `psql -U postgres -l | grep digital_wallet`

**Port already in use?**
- Change `SERVER_PORT` in `.env` file
- Update port mapping in `docker-compose.yml` if using Docker
