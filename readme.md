# OTP Flow Implementation

A simple One-Time Password (OTP) system built with Go, Gin, and PostgreSQL.

## Features

- Generate 6-digit OTP codes with 2-minute expiration
- Validate OTP codes with expiration checking
- Clean architecture with repository, service, and handler layers
- Comprehensive unit tests with mocking
- PostgreSQL database persistence

## Prerequisites

- Go 1.16 or higher
- PostgreSQL 12 or higher
- Docker (optional, for running PostgreSQL)

## Project Structure

```
prenup/
├── main.go                          # Application entry point
├── go.mod / go.sum                  # Dependencies
├── .env                             # Environment variables (git ignored)
├── .gitignore                       # Git ignore rules
├── config/
│   └── database.go                  # Database connection & setup
├── models/
│   └── otp.go                       # Data models
├── handlers/
│   ├── otp_handler.go               # HTTP handlers
│   └── otp_handler_test.go          # Handler tests
├── repository/
│   └── otp_repository.go            # Database operations
├── services/
│   ├── otp_service.go               # Business logic
│   └── otp_service_test.go          # Service tests
```

## Setup PostgreSQL

```bash
createuser -P postgres
createdb -U postgres otp_db
```

## Environment Variables

Create a `.env` file in the project root:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=otp_db
```

## Running the Project

### 1. Install Dependencies
```bash
go mod download
```

### 2. Run the Server
```bash
go run main.go
```

Expected output:
```
2025/11/03 21:23:58 Server starting on port 8080...
```

### 3. Test the Endpoints

**Request OTP:**
```bash
curl -X POST http://localhost:8080/otp/request \
  -H "Content-Type: application/json" \
  -d '{"user_id": "Robert"}'
```

Response:
```json
{
  "user_id": "Robert",
  "otp": "123456"
}
```

**Validate OTP:**
```bash
curl -X POST http://localhost:8080/otp/validate \
  -H "Content-Type: application/json" \
  -d '{"user_id": "Robert", "otp": "123456"}'
```

Response (success):
```json
{
  "user_id": "Robert",
  "message": "OTP validated successfully."
}
```

Response (failure):
```json
{
  "error": "otp_not_found",
  "error_description": "OTP Not Found"
}
```

## Running Tests

```bash
go test ./... -v
```

All tests should pass:
- 6 service tests
- 5 handler tests
- Total: 11 passing tests

## API Endpoints

### POST /otp/request
Generate a new OTP for a user.

**Request:**
```json
{
  "user_id": "Robert"
}
```

**Response:** `200 OK`
```json
{
  "user_id": "Robert",
  "otp": "123456"
}
```

### POST /otp/validate
Validate an OTP code.

**Request:**
```json
{
  "user_id": "Robert",
  "otp": "123456"
}
```

**Response:** `200 OK`
```json
{
  "user_id": "Robert",
  "message": "OTP validated successfully."
}
```

**Response:** `401 Unauthorized` (invalid/expired OTP)
```json
{
  "error": "otp_not_found",
  "error_description": "OTP Not Found"
}
```
