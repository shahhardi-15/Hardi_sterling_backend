# Sterling HMS - Quick Start Guide

## Project Overview

Sterling HMS (Hospital Management System) consists of:
- **Frontend**: Vue 3 + Vite with Tailwind CSS
- **Backend**: Go with Gin framework
- **Database**: PostgreSQL

## Installation Steps

### Step 1: Install Go (If not already installed)

Download and install Go 1.22+ from https://golang.org/dl/

Verify installation:
```bash
go version
```

### Step 2: PostgreSQL Setup

#### Option A: Create Database with psql

```bash
# Open PostgreSQL
psql -U postgres

# Create database
CREATE DATABASE sterling_hms;

# Exit psql
\q
```

#### Option B: Using createdb command

```bash
createdb sterling_hms -U postgres
```

#### Initialize Database Schema

```bash
# Navigate to backend directory
cd f:\Hardi_sterling_backend

# Run schema script
psql -U postgres -d sterling_hms -f database/schema.sql
```

### Step 3: Backend Setup

```bash
# Navigate to backend
cd f:\Hardi_sterling_backend

# Download dependencies
go mod download

# Verify environment variables
cat .env
```

**Update `.env` if needed:**
```env
DB_HOST=localhost
DB_PORT=5432
DB_NAME=sterling_hms
DB_USER=postgres
DB_PASSWORD=your_actual_password  # Update this!
```

### Step 4: Start Backend

#### Option A: Direct execution
```bash
go run cmd/main.go
```

#### Option B: Build and run
```bash
go build -o sterling-hms-backend cmd/main.go
./sterling-hms-backend
```

#### Option C: With hot reload (development)
```bash
# Install air
go install github.com/cosmtrek/air@latest

# Run with hot reload
air
```

Expected output:
```
Server is running on http://localhost:5000
Environment: development
```

### Step 5: Frontend Setup

In a new terminal:

```bash
cd f:\Hardi_Sterling_frontend\sterling-hms-frontend

# Install dependencies
npm install

# Start development server
npm run dev
```

Frontend will be available at: `http://localhost:5173`

## Verify Setup

### Backend Health Check
```bash
curl http://localhost:5000/health
```

Should return:
```json
{"message": "Server is running"}
```

### Test Sign Up
```bash
curl -X POST http://localhost:5000/api/auth/signup \
  -H "Content-Type: application/json" \
  -d '{
    "firstName": "John",
    "lastName": "Doe",
    "email": "john@example.com",
    "password": "securepassword123"
  }'
```

### Frontend Access
Open browser and go to: `http://localhost:5173`

You should see the login page.

## Troubleshooting

### Port Already in Use (Windows PowerShell)
```bash
# Find process using port 5000
Get-NetTCPConnection -LocalPort 5000

# Kill process
Stop-Process -Id <PID> -Force
```

### Database Connection Error
1. Verify PostgreSQL is running:
   ```bash
   psql -U postgres -c "SELECT version();"
   ```

2. Check database exists:
   ```bash
   psql -U postgres -l | grep sterling_hms
   ```

3. Verify credentials in `.env`

### Go Module Issues
```bash
# Clean and download again
go clean -modcache
go mod download
go mod tidy
```

### Frontend Won't Start
```bash
cd f:\Hardi_Sterling_frontend\sterling-hms-frontend
rm -r node_modules package-lock.json
npm install
npm run dev
```

## Development Workflow

### 1. Make Changes

**Backend:**
- Edit files in `internal/`
- If using `air`, changes are auto-reloaded
- Otherwise, restart `go run cmd/main.go`

**Frontend:**
- Edit files in `src/`
- Changes auto-hot-reload in browser

### 2. Test API Endpoints

Using curl, Postman, or the frontend UI:

```bash
# Sign Up
POST http://localhost:5000/api/auth/signup

# Sign In
POST http://localhost:5000/api/auth/signin

# Get Current User (requires token)
GET http://localhost:5000/api/auth/me
Authorization: Bearer <token>
```

### 3. Database Queries

```bash
# Connect to database
psql -U postgres -d sterling_hms

# View users table
SELECT * FROM users;

# Check indexes
\d users
```

## File Structure Reference

```
Hardi_Sterling_frontend/
  sterling-hms-frontend/
    src/
      ├── views/          # Pages (Login, SignUp, Dashboard)
      ├── stores/         # Pinia stores (auth store)
      ├── api/            # API client
      ├── router/         # Vue Router
      └── App.vue

Hardi_sterling_backend/
  ├── cmd/
  │   └── main.go         # Entry point
  ├── internal/
  │   ├── config/         # Configuration
  │   ├── handlers/       # HTTP handlers
  │   ├── middleware/     # Auth middleware
  │   ├── models/         # Data models
  │   └── repositories/   # Database layer
  ├── database/
  │   └── schema.sql      # Database schema
  └── go.mod
```

## Next Steps

After getting the basics running:

1. Create more pages (Patients, Appointments, Staff management)
2. Add more database tables for different entities
3. Implement business logic in handlers
4. Add pagination, filtering, sorting
5. Implement error logging
6. Add unit and integration tests
7. Deploy to production

## Useful Commands

```bash
# Backend
cd f:\Hardi_sterling_backend
go mod tidy              # Tidy dependencies
go run cmd/main.go       # Run directly
go build -o app .        # Build binary
go test ./...            # Run tests

# Frontend
cd f:\Hardi_Sterling_frontend\sterling-hms-frontend
npm install              # Install dependencies
npm run dev              # Start dev server
npm run build            # Build for production
npm run preview          # Preview production build
```

## Support

For issues, check:
1. Backend README: `f:\Hardi_sterling_backend\README.md`
2. Make sure environment variables are set correctly
3. PostgreSQL is running and accessible
4. No port conflicts (5000 for backend, 5173 for frontend)
