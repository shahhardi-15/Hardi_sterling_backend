# Forgot Password Feature - Backend Implementation

## Overview

Complete forgot password feature for Sterling HMS with OTP-based verification, rate limiting, account lockout, and comprehensive audit logging.

## Architecture

### Database Tables

#### `password_reset_tokens`
Stores password reset tokens and OTP hashes.

```sql
- id (SERIAL PRIMARY KEY)
- user_id (FOREIGN KEY → users)
- token_hash (VARCHAR, UNIQUE) - Hash of the reset token
- otp_hash (VARCHAR) - Hash of the OTP
- expires_at (TIMESTAMP) - Token expiration time
- reset_at (TIMESTAMP) - When the password was reset
- is_used (BOOLEAN) - Whether the token was already used
- created_at (TIMESTAMP) - Creation timestamp
```

#### `password_reset_logs`
Audit trail for all password reset actions.

```sql
- id (SERIAL PRIMARY KEY)
- user_id (FOREIGN KEY → users)
- email (VARCHAR) - Email used in the request
- action (VARCHAR) - Action type (see constants below)
- ip_address (VARCHAR) - Client IP address
- user_agent (TEXT) - User agent string
- success (BOOLEAN) - Whether the action succeeded
- error_message (TEXT) - Error details if failed
- created_at (TIMESTAMP) - Log timestamp
```

#### `otp_lockouts`
Tracks failed OTP attempts and account lockouts.

```sql
- id (SERIAL PRIMARY KEY)
- user_id (FOREIGN KEY → users, UNIQUE)
- failed_attempts (INTEGER) - Number of failed attempts
- locked_until (TIMESTAMP) - When the lockout expires
- created_at (TIMESTAMP)
- updated_at (TIMESTAMP)
```

### API Endpoints

#### 1. POST `/api/auth/forgot-password`

**Request:**
```json
{
  "email": "user@example.com"
}
```

**Response (Success):**
```json
{
  "message": "If an account exists with this email, you will receive an OTP.",
  "otpSent": true
}
```

**Features:**
- Rate limited: 5 requests per hour per email
- Generic response to prevent user enumeration
- OTP generated and emailed (10-minute expiry)
- Reset token generated for later use
- All attempts logged with IP address

---

#### 2. POST `/api/auth/verify-otp`

**Request:**
```json
{
  "email": "user@example.com",
  "otp": "123456"
}
```

**Response (Success):**
```json
{
  "message": "OTP verified successfully",
  "resetToken": "a1b2c3d4..." // Use this for password reset
}
```

**Status Codes:**
- `200` - OTP verified successfully
- `401` - Invalid OTP
- `429` - Too many failed attempts (account locked for 30 min)

**Security:**
- 3 failed attempts → 30-minute lockout
- Failed attempts tracked per user
- OTP expiry: 10 minutes

---

#### 3. POST `/api/auth/reset-password`

**Request:**
```json
{
  "resetToken": "a1b2c3d4...",
  "password": "NewSecurePassword123!"
}
```

**Response (Success):**
```json
{
  "message": "Password reset successfully. Please log in with your new password.",
  "success": true
}
```

**Password Validation:**
- Minimum 8 characters
- At least one uppercase letter (A-Z)
- At least one lowercase letter (a-z)
- At least one number (0-9)
- At least one special character (!@#$%^&*)
- Not longer than 128 characters
- Not a commonly used weak password

**Security:**
- Token must be valid and not expired (60-minute window)
- Token marked as used after reset
- Lockout cleared after successful reset
- Confirmation email sent

---

#### 4. POST `/api/auth/resend-otp`

**Request:**
```json
{
  "email": "user@example.com"
}
```

**Response:**
```json
{
  "message": "OTP resent successfully",
  "otpSent": true
}
```

**Features:**
- Generates new OTP if existing one expired
- Resends to same email
- Rate limited to prevent abuse

---

## Security Features

### 1. **OTP Hashing**
- OTP hashed using SHA256 before storage
- Never stored in plaintext
- Not transmitted after verification

### 2. **Reset Token Hashing**
- Reset token hashed using SHA256 before storage
- Token hash used as identifier in verify response
- Cannot be reverse-engineered

### 3. **Rate Limiting**
- **Forgot Password**: 5 requests per hour per email
- **OTP Resend**: Limited by forgot password rate limit
- Prevents brute force attacks

### 4. **Account Lockout**
- 3 failed OTP attempts → 30-minute lockout
- Automatic lockout after threshold exceeded
- Lockout cleared after successful password reset

### 5. **Generic Error Messages**
Prevents user enumeration attacks:
- Non-existent accounts treated same as existing ones
- "If an account exists..." response pattern
- No email confirmation whether account exists

### 6. **Expiry Times**
- OTP: 10 minutes
- Reset Token: 60 minutes
- Automatic cleanup of expired tokens

### 7. **Comprehensive Audit Logging**
All attempts logged with:
- IP address
- User agent
- Email
- User ID (if applicable)
- Action type
- Success/failure status
- Error details

### 8. **Email Security**
- OTP sent via email
- Password reset confirmation email
- Success email for password changes
- Security warnings in email templates

---

## Implementation Details

### Package Structure

```
internal/
├── handlers/
│   └── auth_handler.go          # Password reset handlers
├── repositories/
│   └── password_reset_repository.go  # Database queries
├── models/
│   └── models.go                # Data models
└── utils/
    ├── email_service.go         # Email sending
    ├── password_validator.go     # Password validation
    ├── security.go              # OTP/Token generation, lockout
    └── audit_log.go             # Audit logging
```

### Configuration

Add to your `.env` file:

```env
# Email Configuration
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SENDER_EMAIL=noreply@sterlinghms.com
SENDER_PASSWORD=your_app_specific_password
FROM_NAME=Sterling HMS

# Environment
ENV=development  # Set to production for real emails
```

---

## Database Migration

Run the migration to create tables:

```bash
psql -U postgres -d sterling_hms -f database/password_reset_migration.sql
```

Tables created:
- `password_reset_tokens`
- `password_reset_logs`
- `otp_lockouts`

Indexes created for performance optimization.

---

## Key Constants

```go
// Default expiry times
OTPExpiry: 10 minutes
ResetTokenExpiry: 60 minutes
OTPLockoutDuration: 30 minutes

// Rate limiting
MaxResetRequests: 5 per hour
MaxOTPAttempts: 3 before lockout

// Security
OTPLength: 6 digits
TokenLength: 32 bytes (hex-encoded)
PasswordMinLength: 8
PasswordMaxLength: 128
```

---

## Flow Diagrams

### Complete Password Reset Flow

```
User initiates forgot password
         ↓
Validate email & rate limit
         ↓
Generate OTP (6 digits) + Reset Token (32 bytes)
         ↓
Hash both values using SHA256
         ↓
Store in password_reset_tokens table
         ↓
Send OTP via email
         ↓
User receives email with OTP
         ↓
User submits OTP for verification
         ↓
Verify OTP hash matches stored hash
         ↓
Return reset token to frontend
         ↓
User submits new password + reset token
         ↓
Validate password strength
         ↓
Hash new password with bcrypt
         ↓
Update user password
         ↓
Mark token as used
         ↓
Clear any lockouts
         ↓
Send confirmation email
         ↓
Success response
```

### Failed Attempt Handling

```
User submits wrong OTP
         ↓
Increment failed_attempts counter
         ↓
If attempts >= 3:
   → Lock account for 30 minutes
   → Return 429 Too Many Requests
         ↓
Log failed attempt with IP address
         ↓
User can try again (if not locked out)
```

---

## Error Codes

| Code | Scenario |
|------|----------|
| 200 | Success |
| 400 | Invalid request data |
| 401 | Invalid/expired OTP or token |
| 429 | Rate limited or account locked |
| 500 | Server error |

---

## Monitoring

### Audit Logs Query Examples

**Track password resets by user:**
```sql
SELECT * FROM password_reset_logs 
WHERE user_id = 1 
ORDER BY created_at DESC;
```

**Find suspicious activity (multiple failed attempts):**
```sql
SELECT email, COUNT(*) as failed_attempts, 
       MAX(created_at) as last_attempt
FROM password_reset_logs
WHERE action = 'otp_failed' 
      AND created_at > NOW() - INTERVAL '24 hours'
GROUP BY email
HAVING COUNT(*) > 5
ORDER BY failed_attempts DESC;
```

**IP-based analysis:**
```sql
SELECT ip_address, COUNT(*) as attempts, action
FROM password_reset_logs
WHERE created_at > NOW() - INTERVAL '1 hour'
GROUP BY ip_address, action
ORDER BY attempts DESC;
```

---

## Testing

### Test the Flow

1. **Forgot Password:**
```bash
curl -X POST http://localhost:5000/api/auth/forgot-password \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com"}'
```

2. **Verify OTP:**
```bash
curl -X POST http://localhost:5000/api/auth/verify-otp \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","otp":"123456"}'
```

3. **Reset Password:**
```bash
curl -X POST http://localhost:5000/api/auth/reset-password \
  -H "Content-Type: application/json" \
  -d '{"resetToken":"token_from_verify","password":"NewPass123!"}'
```

---

## Future Enhancements

1. **SMS OTP** - Alternative to email
2. **Passwordless Login** - Use OTP for direct login
3. **Security Questions** - Additional verification layer
4. **Device Fingerprinting** - Detect suspicious locations
5. **Email Verification** - Verify email during signup
6. **Two-Factor Authentication** - MFA support
7. **Recovery Codes** - Backup access methods
8. **Anomaly Detection** - ML-based suspicious activity detection

---

## Troubleshooting

### OTP Not Received
- Check `.env` email configuration
- Set `ENV=development` to log OTPs to console
- Verify SMTP credentials are correct

### Token Expired Error
- OTP expires in 10 minutes
- Reset token expires in 60 minutes
- User must restart process if expired

### Account Locked
- 30-minute automatic lockout after 3 failed OTP attempts
- Lockout cleared after successful password reset
- Or wait 30 minutes and try again

---

## Security Checklist

- [x] OTP hashed before storage
- [x] Tokens hashed before storage
- [x] Rate limiting implemented
- [x] Account lockout on failed attempts
- [x] Generic error messages
- [x] IP address logging
- [x] User agent logging
- [x] Password strength validation
- [x] Email confirmation sent
- [x] Automatic token expiry
- [x] HTTPS recommended for production
- [x] CORS properly configured

