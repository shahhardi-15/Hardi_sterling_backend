# Sterling HMS Admin System - Quick Start Guide

## Getting Started

### Prerequisites
- Node.js and npm (for frontend)
- Go 1.21+ (for backend)
- PostgreSQL 18+
- Backend server running on http://localhost:5000
- Frontend dev server running on http://localhost:5173

### Step 1: Ensure Database Migrations are Applied

The migrations should have been applied automatically. To manually apply them:

```bash
cd f:\Hardi_sterling_backend

# Set PostgreSQL password
$env:PGPASSWORD='admin'

# Apply admin migrations
& "C:\Program Files\PostgreSQL\18\bin\psql.exe" -U postgres -h localhost -d sterling -f "database\admin_users_migration.sql"

# Apply appointments migrations (if not done)
& "C:\Program Files\PostgreSQL\18\bin\psql.exe" -U postgres -h localhost -d sterling -f "database\appointments_migration.sql"
```

### Step 2: Start the Backend Server

```bash
cd f:\Hardi_sterling_backend
go run .\cmd\main.go
```

You should see:
```
2026/03/27 15:42:21 Server is running on http://localhost:5000
```

### Step 3: Start the Frontend Dev Server

```bash
cd f:\Hardi_Sterling_frontend\sterling-hms-frontend
npm run dev
```

Frontend will be available at http://localhost:5173

### Step 4: Test Admin Login

1. **Navigate to Admin Login:**
   - Open browser to: http://localhost:5173/admin/login

2. **Enter Credentials:**
   - Email: `adminsterling@gmail.com`
   - Password: `admin@123`

3. **Submit:**
   - Click "Sign In as Admin" or press Enter

4. **Expected Result:**
   - Redirected to `/admin/dashboard`
   - Dashboard displays statistics cards
   - Welcome message shows "Admin Sterling"

### Step 5: Test Dashboard

Once logged in:

1. **View Statistics:**
   - Total Patients: Shows number of registered patients
   - Total Appointments: Shows number of appointments
   - Total Doctors: Shows number of doctors (should be 5 from sample data)
   - Total Staff: Shows number of admin users

2. **Quick Links:**
   - Manage Patients
   - View Appointments
   - Manage Doctors
   - View Reports

3. **Logout:**
   - Click "Logout" button in top-right
   - Redirected back to `/admin/login`
   - Token removed from localStorage

## Testing Admin API Directly

### Test 1: Admin Login

```powershell
$uri = "http://localhost:5000/api/admin/login"
$headers = @{"Content-Type" = "application/json"}
$body = @{
    email = "adminsterling@gmail.com"
    password = "admin@123"
} | ConvertTo-Json

$response = Invoke-WebRequest -Uri $uri -Method Post -Headers $headers -Body $body
$response.Content | ConvertFrom-Json | ConvertTo-Json
```

**Expected Response:**
```json
{
  "message": "Login successful",
  "success": true,
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "admin": {
    "id": 1,
    "email": "adminsterling@gmail.com",
    "name": "Admin Sterling",
    "role": "admin",
    "createdAt": "2026-03-27T...",
    "updatedAt": "2026-03-27T...",
    "isActive": true
  }
}
```

### Test 2: Get Dashboard Statistics

```powershell
$token = "<paste token from login response>"
$uri = "http://localhost:5000/api/admin/dashboard/stats"
$headers = @{
    "Authorization" = "Bearer $token"
    "Content-Type" = "application/json"
}

$response = Invoke-WebRequest -Uri $uri -Method Get -Headers $headers
$response.Content | ConvertFrom-Json | ConvertTo-Json
```

**Expected Response:**
```json
{
  "message": "Statistics retrieved successfully",
  "success": true,
  "stats": {
    "totalPatients": 0,
    "totalAppointments": 0,
    "totalDoctors": 5,
    "totalStaff": 1
  }
}
```

### Test 3: Invalid Login

```powershell
$uri = "http://localhost:5000/api/admin/login"
$headers = @{"Content-Type" = "application/json"}
$body = @{
    email = "wrong@email.com"
    password = "wrongpassword"
} | ConvertTo-Json

try {
    $response = Invoke-WebRequest -Uri $uri -Method Post -Headers $headers -Body $body
} catch {
    $_.Exception.Response.Content | ConvertFrom-Json | ConvertTo-Json
}
```

**Expected Response:**
```json
{
  "message": "Invalid credentials",
  "success": false
}
```

## Common Issues & Troubleshooting

### Admin Won't Login
**Problem:** Getting "Invalid credentials" error

**Solutions:**
1. Verify database migrations were applied to correct database:
   ```bash
   $env:PGPASSWORD='admin'; & "C:\Program Files\PostgreSQL\18\bin\psql.exe" -U postgres -h localhost -d sterling -c "SELECT * FROM admin_users;"
   ```

2. Verify admin user exists:
   ```
   id  |               email               |      name      | role  | is_active
   ----+-----------------------------------+----------------+-------+-----------
   1   | adminsterling@gmail.com           | Admin Sterling | admin | t
   ```

3. Check password hash is correct:
   ```bash
   & "C:\Program Files\PostgreSQL\18\bin\psql.exe" -U postgres -h localhost -d sterling -f "database\fix_admin_password.sql"
   ```

### Dashboard Shows No Data
**Problem:** Statistics show 0 for all fields

**Solution:** This is normal! Add sample data first:
```bash
# Create a test patient
# Add appointments, etc.
```

### Backend Server Won't Start
**Problem:** Address already in use or module errors

**Solutions:**
1. Kill existing processes:
   ```bash
   Get-Process | Where-Object {$_.ProcessName -like "*go*"} | Stop-Process -Force
   ```

2. Rebuild:
   ```bash
   cd f:\Hardi_sterling_backend
   go clean
   go build .\cmd\main.go
   ```

### CORS Errors
**Problem:** Frontend can't reach backend API

**Solutions:**
1. Verify backend is running on http://localhost:5000
2. Check CORS allowlist in `cmd/main.go` includes http://localhost:5173
3. Check frontend environment variable `VITE_API_URL=http://localhost:5000`

## Features Implemented

✅ **Authentication**
- Secure admin login with bcrypt password hashing
- JWT token generation and validation
- 24-hour token expiry
- Role-based access control (admin only)

✅ **Authorization Middleware**
- Protects admin dashboard routes
- Validates JWT tokens
- Checks admin role

✅ **Admin Dashboard**
- Display statistics (patients, appointments, doctors, staff)
- System status monitoring
- Quick links to admin modules
- Mobile-responsive design

✅ **Security Features**
- Generic error messages (prevents email enumeration)
- Audit logging for all admin actions
- HTTPS/CORS support
- Token stored securely in localStorage
- Authorization header validation

✅ **API Endpoints**
- POST /api/admin/login - Admin authentication
- GET /api/admin/dashboard/stats - Retrieve statistics
- POST /api/admin/logout - Logout admin user

## Future Enhancements

- [ ] Admin management (add/remove admins)
- [ ] Patient management interface
- [ ] Appointment management interface
- [ ] Doctor management interface
- [ ] Report generation and export
- [ ] Advanced analytics
- [ ] User activity logs viewer
- [ ] Password change functionality
- [ ] Multi-factor authentication
- [ ] Session management

## Support

For issues or questions, check:
1. [ADMIN_IMPLEMENTATION_SUMMARY.md](./ADMIN_IMPLEMENTATION_SUMMARY.md) - Detailed technical documentation
2. Backend logs at http://localhost:5000/
3. Browser console for frontend errors
4. Database logs from PostgreSQL

---

**Last Updated:** March 27, 2026
