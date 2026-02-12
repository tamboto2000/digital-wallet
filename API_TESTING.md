# API Testing Guide

This guide provides step-by-step instructions for testing the Digital Wallet API.

## Prerequisites

- Server running on `http://localhost:8080`
- PostgreSQL database with sample data loaded
- curl, Postman, or any HTTP client

## Test Scenarios

### 1. Health Check

Verify the service is running.

```bash
curl http://localhost:8080/health
```

**Expected Response:**
```json
{
  "status": "healthy"
}
```

---

### 2. Balance Inquiry Tests

#### Test 2.1: Get balance for existing user

```bash
curl http://localhost:8080/api/balance/1
```

**Expected Response (200 OK):**
```json
{
  "user_id": 1,
  "balance": 1000.00
}
```

#### Test 2.2: Get balance for non-existent user

```bash
curl http://localhost:8080/api/balance/999
```

**Expected Response (404 Not Found):**
```json
{
  "error": "Wallet not found"
}
```

#### Test 2.3: Get balance with invalid user ID

```bash
curl http://localhost:8080/api/balance/abc
```

**Expected Response (400 Bad Request):**
```json
{
  "error": "Invalid user ID"
}
```

---

### 3. Withdrawal Tests

#### Test 3.1: Successful withdrawal

```bash
curl -X POST http://localhost:8080/api/withdraw/1 \
  -H "Content-Type: application/json" \
  -d '{"amount": 100.50}'
```

**Expected Response (200 OK):**
```json
{
  "user_id": 1,
  "amount": 100.50,
  "balance_before": 1000.00,
  "balance_after": 899.50,
  "transaction_id": 1
}
```

**Verify:**
```bash
curl http://localhost:8080/api/balance/1
```

Should show updated balance: `899.50`

#### Test 3.2: Withdrawal with insufficient funds

```bash
curl -X POST http://localhost:8080/api/withdraw/2 \
  -H "Content-Type: application/json" \
  -d '{"amount": 1000.00}'
```

**Expected Response (400 Bad Request):**
```json
{
  "error": "Insufficient funds"
}
```

**Note:** User 2 starts with $500.00

#### Test 3.3: Withdrawal with zero amount

```bash
curl -X POST http://localhost:8080/api/withdraw/1 \
  -H "Content-Type: application/json" \
  -d '{"amount": 0}'
```

**Expected Response (400 Bad Request):**
```json
{
  "error": "Amount must be greater than 0"
}
```

#### Test 3.4: Withdrawal with negative amount

```bash
curl -X POST http://localhost:8080/api/withdraw/1 \
  -H "Content-Type: application/json" \
  -d '{"amount": -50.00}'
```

**Expected Response (400 Bad Request):**
```json
{
  "error": "Amount must be greater than 0"
}
```

#### Test 3.5: Withdrawal for non-existent user

```bash
curl -X POST http://localhost:8080/api/withdraw/999 \
  -H "Content-Type: application/json" \
  -d '{"amount": 50.00}'
```

**Expected Response (404 Not Found):**
```json
{
  "error": "Wallet not found"
}
```

#### Test 3.6: Withdrawal with invalid JSON body

```bash
curl -X POST http://localhost:8080/api/withdraw/1 \
  -H "Content-Type: application/json" \
  -d 'invalid json'
```

**Expected Response (400 Bad Request):**
```json
{
  "error": "Invalid request body"
}
```

---

### 4. Concurrency Test

Test that concurrent withdrawals don't cause race conditions.

Create a file `concurrent_test.sh`:

```bash
#!/bin/bash

# Make 10 concurrent $50 withdrawals from user 1
for i in {1..10}; do
  curl -X POST http://localhost:8080/api/withdraw/1 \
    -H "Content-Type: application/json" \
    -d '{"amount": 50.00}' &
done

wait

# Check final balance
curl http://localhost:8080/api/balance/1
```

**Expected Behavior:**
- User 1 starts with $1000
- Some requests should succeed, some should fail with "insufficient funds"
- Final balance should be exactly: $1000 - (number of successful withdrawals × $50)
- No negative balance should occur

---

### 5. Edge Cases

#### Test 5.1: Large withdrawal amount

```bash
curl -X POST http://localhost:8080/api/withdraw/3 \
  -H "Content-Type: application/json" \
  -d '{"amount": 2500.50}'
```

**Expected Response (200 OK):**
Should succeed and leave balance at $0.00

#### Test 5.2: Decimal precision

```bash
curl -X POST http://localhost:8080/api/withdraw/1 \
  -H "Content-Type: application/json" \
  -d '{"amount": 33.33}'
```

**Expected Response (200 OK):**
Should handle decimal amounts correctly

---

## Database Verification

After running tests, verify transaction history:

```bash
psql -U postgres -d digital_wallet -c "SELECT * FROM transactions ORDER BY created_at DESC LIMIT 10;"
```

This should show all recorded transactions with:
- Transaction ID
- Wallet ID
- Amount
- Type (withdraw)
- Balance before
- Balance after
- Created timestamp

---

## Postman Collection

Import `postman_collection.json` for a complete set of pre-configured test requests:

1. Open Postman
2. Click "Import"
3. Select `postman_collection.json`
4. Run the entire collection or individual requests

---

## Performance Testing

### Using Apache Bench (ab)

```bash
# Install apache bench
sudo apt-get install apache2-utils  # Ubuntu/Debian
brew install httpd                   # macOS

# Test balance endpoint (100 requests, 10 concurrent)
ab -n 100 -c 10 http://localhost:8080/api/balance/1

# Test withdrawal endpoint
ab -n 50 -c 5 -p withdraw.json -T application/json http://localhost:8080/api/withdraw/1
```

Create `withdraw.json`:
```json
{"amount": 1.00}
```

---

## Test Checklist

- [ ] Server starts without errors
- [ ] Health check returns 200 OK
- [ ] Balance inquiry works for existing user
- [ ] Balance inquiry returns 404 for non-existent user
- [ ] Successful withdrawal updates balance
- [ ] Insufficient funds returns 400 error
- [ ] Zero/negative amounts return 400 error
- [ ] Invalid user ID returns 400 error
- [ ] Non-existent user returns 404 error
- [ ] Concurrent withdrawals don't create negative balance
- [ ] All transactions are recorded in database
- [ ] Balance never goes below zero

---

## Troubleshooting

### Server won't start
- Check if PostgreSQL is running: `systemctl status postgresql`
- Verify database exists: `psql -U postgres -l | grep digital_wallet`
- Check environment variables in `.env` file

### Connection refused errors
- Ensure server is running on port 8080
- Check firewall settings
- Verify no other service is using port 8080

### Database errors
- Verify PostgreSQL is accessible
- Check database credentials in `.env`
- Ensure schema has been applied: `psql -U postgres -d digital_wallet -f schema.sql`

### Unexpected responses
- Check server logs for detailed error messages
- Verify request format matches API documentation
- Ensure Content-Type header is set for POST requests
