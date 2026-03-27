# Setup Instructions - Forgot Password Feature

## Prerequisites

- Go 1.22+
- PostgreSQL 12+
- Existing Sterling HMS backend running

## Step 1: Database Migration

Run the password reset migration to create the required tables:

```bash
# From the backend directory
cd f:\Hardi_sterling_backend

# Run the migration
psql -U postgres -d sterling_hms -f database/password_reset_migration.sql
```

Verify tables were created:

```bash
psql -U postgres -d sterling_hms -c "\dt password_reset*"
psql -U postgres -d sterling_hms -c "\dt otp_*"
```

## Step 2: Environment Configuration

Add the following to your `.env` file in the backend directory:

```env
# Email Configuration (Gmail example)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SENDER_EMAIL=your-email@gmail.com
SENDER_PASSWORD=your_app_password  # Google App Password, not regular password
FROM_NAME=Sterling HMS

# Environment (development logs emails to console instead of sending)
ENV=development
```

### Email Setup Instructions

#### For Gmail:
1. Enable 2-Step Verification
2. Generate App Password: https://myaccount.google.com/apppasswords
3. Use the generated password in `SENDER_PASSWORD`

#### For Custom SMTP:
Replace `SMTP_HOST` and `SMTP_PORT` with your provider's settings.

#### For Development:
Set `ENV=development` to log emails to console instead of sending them.

## Step 3: Test the Installation

Start the backend server:

```bash
cd f:\Hardi_sterling_backend
go run cmd/main.go
```

Test the forgot password endpoint:

```bash
curl -X POST http://localhost:5000/api/auth/forgot-password \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com"}'
```

Expected response:
```json
{
  "message": "If an account exists with this email, you will receive an OTP.",
  "otpSent": true
}
```

## Step 4: Enable in Frontend

The backend is ready! Now update the Vue.js frontend to use these endpoints. See [FRONTEND_SETUP.md](./FRONTEND_SETUP.md) for frontend implementation.

## File Structure Added

```
internal/
├── handlers/
│   └── auth_handler.go               # Updated with 4 new methods
│       ├── ForgotPassword()
│       ├── VerifyOTP()
│       ├── ResetPassword()
│       └── ResendOTP()
├── repositories/
│   └── password_reset_repository.go  # NEW - Password reset DB operations
├── models/
│   └── models.go                     # Updated with new DTOs
└── utils/
    ├── email_service.go              # NEW - Email sending
    ├── password_validator.go         # NEW - Password validation
    ├── security.go                   # NEW - OTP/Token/Lockout management
    └── audit_log.go                  # NEW - Audit logging

database/
└── password_reset_migration.sql      # NEW - Database schema

cmd/
└── main.go                           # Updated with new routes
```

## API Endpoints Available

After completing setup, these endpoints are available:

```
POST   /api/auth/forgot-password      Request password reset OTP
POST   /api/auth/verify-otp           Verify OTP and get reset token
POST   /api/auth/reset-password       Update password with reset token
POST   /api/auth/resend-otp           Resend OTP to email
```

## Security Features Enabled

✅ Rate limiting (5 requests/hour per email)
✅ Account lockout (3 failed OTP attempts = 30 min)
✅ OTP hashing with SHA256
✅ Token hashing with SHA256
✅ Password strength validation (8+ chars, uppercase, lowercase, number, special)
✅ Generic error messages (prevent user enumeration)
✅ Audit logging with IP address
✅ Email confirmation sent
✅ Token expiry (OTP: 10 min, Reset: 60 min)

## Testing Checklist

- [ ] Database migration successful (`\dt` shows new tables)
- [ ] Environment variables configured
- [ ] Server starts without errors
- [ ] Forgot password endpoint returns 200
- [ ] OTP appears in logs (development mode)
- [ ] Rate limiting works (6th request returns 429)
- [ ] Wrong OTP locked out after 3 attempts
- [ ] Database logs recorded (check password_reset_logs table)

## Troubleshooting

### Migration Failed
```bash
# Check migration syntax
psql -U postgres -d sterling_hms
# Then paste the SQL from password_reset_migration.sql
```

### Email Not Sending
- Verify `ENV=development` is set (check console for logged OTP)
- Check SMTP credentials in `.env`
- For Gmail: Use App Password, not regular password

### Build Errors
```bash
# Update dependencies
go mod tidy
go mod download
go run cmd/main.go
```

### Database Connection Error
```bash
# Verify PostgreSQL is running
# Check connection string in .env
psql -U postgres -d sterling_hms -c "SELECT 1;"
```

## Next Steps

1. ✅ Backend setup complete
2. → [Frontend setup](./FRONTEND_SETUP.md) - Implement Vue.js components
3. → Test end-to-end flow
4. → Deploy to production with HTTPS

## Support

For issues:
1. Check FORGOT_PASSWORD_FEATURE.md for detailed documentation
2. Review audit logs in `password_reset_logs` table
3. Check server console for error messages
4. Verify database tables exist: `psql -d sterling_hms -c "\dt"`

