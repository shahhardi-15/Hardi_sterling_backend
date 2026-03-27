# Sterling HMS Backend - Go

Go/Gin backend for Sterling Hospital Management System with JWT authentication.

## Prerequisites

- Go 1.22 or higher
- PostgreSQL 12 or higher
- Git

## Project Structure

```
├── cmd/
│   └── main.go              # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go        # Configuration and database setup
│   ├── handlers/
│   │   └── auth_handler.go  # Authentication handlers
│   ├── middleware/
│   │   └── auth.go          # JWT middleware
│   ├── models/
│   │   └── models.go        # Data models and request/response structures
│   └── repositories/
│       └── user_repository.go  # Database operations
├── database/
│   └── schema.sql           # Database schema
├── go.mod                   # Go module definition
├── .env                     # Environment variables
└── README.md                # This file
```

## Setup Instructions

### 1. Install Dependencies

```bash
cd f:\Hardi_sterling_backend
go mod download
```

### 2. PostgreSQL Setup

Create the database and initialize the schema:

```bash
# Create database
createdb sterling_hms

# Initialize schema
psql -U postgres -d sterling_hms -f database/schema.sql
```

Or using psql interactive mode:

```bash
psql -U postgres

CREATE DATABASE sterling_hms;
\c sterling_hms
\i database/schema.sql
```

### 3. Environment Variables

Create a `.env` file in the root directory with:

```env
PORT=5000
DB_HOST=localhost
DB_PORT=5432
DB_NAME=sterling_hms
DB_USER=postgres
DB_PASSWORD=your_database_password

JWT_SECRET=your_super_secret_jwt_key_change_in_production
JWT_EXPIRE=168h

ENV=development
```

**Important:** Change `JWT_SECRET` and `DB_PASSWORD` in production.

### 4. Build and Run

```bash
# Run directly with go run
go run cmd/main.go

# Or build and run
go build -o sterling-hms-backend cmd/main.go
./sterling-hms-backend
```

Server will start on `http://localhost:5000`

## API Endpoints

### Health Check
- **GET** `/health`
- **Response:** `{"message": "Server is running"}`

### Sign Up
- **POST** `/api/auth/signup`
- **Body:**
```json
{
  "firstName": "John",
  "lastName": "Doe",
  "email": "john@example.com",
  "password": "securepassword123"
}
```
- **Response:**
```json
{
  "message": "User registered successfully",
  "user": {
    "id": 1,
    "firstName": "John",
    "lastName": "Doe",
    "email": "john@example.com",
    "createdAt": "2024-03-26T10:30:00Z"
  },
  "token": "eyJhbGciOiJIUzI1NiIs..."
}
```

### Sign In
- **POST** `/api/auth/signin`
- **Body:**
```json
{
  "email": "john@example.com",
  "password": "securepassword123"
}
```
- **Response:** Same as Sign Up

### Get Current User (Protected)
- **GET** `/api/auth/me`
- **Headers:** `Authorization: Bearer <token>`
- **Response:**
```json
{
  "user": {
    "id": 1,
    "firstName": "John",
    "lastName": "Doe",
    "email": "john@example.com",
    "createdAt": "2024-03-26T10:30:00Z"
  }
}
```

## Features

✓ User registration with comprehensive validation
✓ User login with JWT authentication
✓ Password hashing with bcrypt (strength: 10)
✓ Protected routes with JWT middleware
✓ PostgreSQL database integration
✓ CORS enabled for frontend communication
✓ Error handling and validation
✓ Gin web framework
✓ Environment-based configuration
✓ Database connection pooling
✓ Structured logging

## Dependencies

- **gin-gonic/gin**: HTTP web framework
- **golang-jwt/jwt**: JWT token generation and validation
- **lib/pq**: PostgreSQL driver
- **golang.org/x/crypto**: Password hashing
- **gin-contrib/cors**: CORS middleware
- **joho/godotenv**: Environment variable loading

Install all dependencies:
```bash
go mod tidy
```

## Security Notes

- Passwords are hashed using bcrypt with strength 10
- JWT tokens expire according to `JWT_EXPIRE` (default: 168h = 7 days)
- Always change `JWT_SECRET` in production
- HTTPS should be used in production
- All sensitive data should be stored in environment variables
- Never commit `.env` file to version control

## Error Handling

The API returns standard HTTP status codes:
- `201`: Resource created successfully
- `400`: Bad request / validation error
- `401`: Unauthorized
- `404`: Not found
- `409`: Conflict (e.g., email already exists)
- `500`: Internal server error

## Development

### Running in Development Mode

```bash
# Uses .env file for configuration
go run cmd/main.go
```

### Building for Production

```bash
# Build optimized binary
go build -ldflags="-s -w" -o sterling-hms-backend cmd/main.go

# Run with custom environment
ENV=production ./sterling-hms-backend
```

## Troubleshooting

### Database Connection Issues
- Ensure PostgreSQL is running
- Verify DB credentials in `.env`
- Check that database exists: `psql -l`

### Port Already in Use
- Change `PORT` in `.env`
- Or kill existing process: `lsof -ti:5000 | xargs kill -9` (Linux/Mac)

### Module Not Found Errors
- Run `go mod tidy`
- Run `go mod download`

## Frontend Integration

The frontend expects the backend at `http://localhost:5000/api`. Update your frontend `.env`:

```env
VITE_API_URL=http://localhost:5000/api
```

The frontend should automatically include the JWT token in the `Authorization` header for protected routes.
