# Forgot Password Implementation - Migration from OTP to Token-Based

## Overview
Successfully migrated the forgot password feature from a complex 3-step OTP-based flow to a simpler 2-step token-based approach with direct email links.

## Changes Summary

### Backend Changes

#### 1. Database Schema (New: `database/password_reset_v2.sql`)
- Removed OTP-related infrastructure
- Dropped `otp_lockouts` table
- Removed `otp_hash` column from `password_reset_tokens` table
- Simplified schema with:
  - `password_reset_tokens` table: id, user_id, token_hash (UNIQUE), expires_at, used_at, is_used, created_at
  - `password_reset_logs` table: id, user_id, email, action, ip_address, user_agent, success, error_message, created_at

#### 2. Data Models (`internal/models/models.go`)
**Removed:**
- `OTPLockout` struct
- `VerifyOTPRequest` / `VerifyOTPResponse` DTOs
- `ResendOTPRequest` / `ResendOTPResponse` DTOs
- `OTPHash` field from `PasswordResetToken`

**Kept/Updated:**
- `PasswordResetToken`: Simplified with TokenHash, ExpiresAt, UsedAt, IsUsed
- `PasswordResetLog`: Unchanged
- `ForgotPasswordRequest`: Email input
- `ForgotPasswordResponse`: Success/Message response
- `ResetPasswordRequest`: ResetToken + Password
- `ResetPasswordResponse`: Success message

#### 3. Repository Layer (`internal/repositories/password_reset_repository.go`)
**Removed Methods:**
- `GetPasswordResetTokenByUserID()` - OTP-specific
- `CountOTPAttempts()` - OTP-specific
- `GetOTPLockout()` / `CreateOrUpdateOTPLockout()` - Lockout management
- `ResetOTPLockout()` - OTP-specific
- `IncrementOTPFailedAttempts()` - OTP-specific
- `IsUserLockedOut()` - OTP lockout check

**Updated Methods:**
- `CreatePasswordResetToken()`: Removed otpHash parameter, now takes only (userID, tokenHash, expiresAt)
- `MarkResetTokenAsUsed()`: Simplified to just mark token as used
- `GetPasswordResetTokenByHash()`: Removed OTPHash scanning

**New Methods:**
- `InvalidateAllUserTokens()`: Cancels all token-based resets for a user
- `CleanupExpiredTokens()`: Removes expired tokens from database

#### 4. Handlers (`internal/handlers/auth_handler.go`)
**Replaced:**
- `ForgotPassword()`: 
  - Now generates single 32-byte SHA256 token only
  - Sends email with reset link containing token in query parameter
  - 1-hour expiry (not 10 minutes)
  - No OTP in email
  
- `ResetPassword()`: 
  - Simplified by removing OTP/lockout validation
  - Only validates token hash and password requirements
  - No lockout management needed
  
**Deleted:**
- `VerifyOTP()`: Entire method removed (no longer needed)
- `ResendOTP()`: Entire method removed (no longer needed)

**Added Import:**
- `"os"` for environment variable access

#### 5. Email Service (`internal/utils/email_service.go`)
**Removed:**
- `SendOTPEmail()`: Was sending 6-digit OTP codes

**Added:**
- `SendPasswordResetEmail()`: Sends clickable reset link (1-hour expiry)
  - Link format: `{FRONTEND_URL}/reset-password?token={token}`
  - HTML template with button and fallback link text

#### 6. Security Utilities (`internal/utils/security.go`)
**Removed:**
- `OTPGenerator` struct and `GenerateOTP()` method
- `LockoutManager` struct and related methods
- OTP-related fields from `PasswordResetFlow`

**Kept:**
- `TokenGenerator`: Still generates secure 32-byte hex tokens
- `RateLimiter`: Rate limiting (5 requests/hour per email)
- Simplified `PasswordResetFlow`: Only tracks token expiry and rate limit window

#### 7. Audit Logging (`internal/utils/audit_log.go`)
**Removed Action Constants:**
- `ActionOTPSent`
- `ActionOTPVerified`
- `ActionOTPFailed`
- `ActionOTPResend`

**Kept:**
- `ActionForgotPasswordRequest`
- `ActionPasswordResetSuccess`
- `ActionPasswordResetFailed` (new)

**Updated Error Messages:**
- Removed OTP-specific error messages
- Updated `InvalidEmail` to reference reset link instead of OTP

#### 8. Routes (`cmd/main.go`)
**Removed Endpoints:**
- `POST /api/auth/verify-otp`
- `POST /api/auth/resend-otp`

**Kept Endpoints:**
- `POST /api/auth/forgot-password` (email input)
- `POST /api/auth/reset-password` (password reset with token)

### Frontend Changes

#### 1. API Client (`src/api/auth.js`)
**Removed Methods:**
- `verifyOtp()`: No longer needed
- `resendOtp()`: No longer needed

**Kept Methods:**
- `forgotPassword(email)`: Request reset link
- `resetPassword(resetToken, password)`: Actually reset password

#### 2. Component (`src/views/ForgotPassword.vue`)
**Changes:**
- Reduced from 3-step to 2-step flow
- Removed OTP input step entirely
- Removed OTP timer logic
- Removed resend OTP functionality
- Component now handles two scenarios:
  1. **Step 1 (Email)**: User enters email, receives reset link via email
  2. **Step 2 (Password)**: User clicks email link or manually enters token, resets password

**Kept Features:**
- Password strength indicator (5-point scale)
- Password requirements checklist
- Error/success messaging
- Responsive design
- Back navigation between steps

**New Logic:**
- Checks URL query parameters for `token` parameter
- If token present in URL, skips Step 1 and shows Step 2 directly

## Security Implications

### Improvements
- ✅ Token sent via email is more secure than in-app OTP entry
- ✅ Email link prevents token interception on same device
- ✅ Single token reduces attack surface
- ✅ Cleaner audit trail without OTP attempts

### Trade-offs
- Email link more browser-dependent (users must click link)
- No intermediate verification step - token is used directly for reset
- Rate limiting maintained at 5 requests/hour to prevent abuse

## Database Migration Instructions

**IMPORTANT:** Before deploying, run the new schema:

```bash
# Option 1: Direct SQL execution
psql -U your_user -d your_db -f database/password_reset_v2.sql

# Option 2: Using Go migrate tool (recommend adding migration tooling)
migrate -path database -database "postgresql://user:pass@localhost/db" up
```

The new schema:
1. Creates fresh `password_reset_tokens` table without OTP fields
2. Creates `password_reset_logs` for audit trail
3. Drops the old `otp_lockouts` table (if exists)

## Testing Checklist

- [ ] Backend compiles without errors
- [ ] Database migrations applied successfully
- [ ] `/api/auth/forgot-password` sends email with reset link
- [ ] Reset link opens to Step 2 with token pre-filled
- [ ] Password reset with valid token succeeds
- [ ] Invalid/expired token shows error
- [ ] Rate limiting blocks after 5 requests
- [ ] Email content displays correctly
- [ ] Frontend handles both email-entry and token-from-URL flows
- [ ] Audit logs capture all actions correctly

## Environment Variables

Ensure these are set in your `.env`:

```
FRONTEND_URL=http://localhost:5173          # or production URL
SENDER_EMAIL=noreply@sterlinghms.com
SENDER_PASSWORD=your_app_password
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
ENV=development                              # or production
```

## Rollback Plan

If needed to revert to OTP version:
1. Restore previous database schema (contains `otp_lockouts` and `otp_hash` fields)
2. Restore previous code from git commit
3. Re-run the original database migration

## Token Flow Diagram

```
User Email Entry
       ↓
Request Forgot Password (/forgot-password)
       ↓
Generate 32-byte Token + Hash
       ↓
Store TokenHash in Database (1-hour expiry)
       ↓
Send Email with Link: /reset-password?token=ABC123XYZ
       ↓
User Clicks Email Link OR Manually Enters Token
       ↓
Post New Password + Token (/reset-password)
       ↓
Verify Token Hash | Hash New Password
       ↓
Update Database | Mark Token as Used
       ↓
Success Response → Redirect to Login
```

## Files Modified Summary

**Backend:**
- ✅ `cmd/main.go` - Routes updated
- ✅ `internal/handlers/auth_handler.go` - 4 methods replaced/removed
- ✅ `internal/models/models.go` - OTP structs removed
- ✅ `internal/repositories/password_reset_repository.go` - Methods refactored
- ✅ `internal/utils/email_service.go` - SendOTPEmail → SendPasswordResetEmail
- ✅ `internal/utils/security.go` - OTP/Lockout utilities removed
- ✅ `internal/utils/audit_log.go` - OTP actions removed
- ✅ `database/password_reset_v2.sql` - New schema (future use)

**Frontend:**
- ✅ `src/api/auth.js` - OTP methods removed
- ✅ `src/views/ForgotPassword.vue` - 3-step → 2-step flow

**Configuration:**
- ✅ Deployment docs (this file)
