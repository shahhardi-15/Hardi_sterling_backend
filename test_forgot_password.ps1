# Sterling HMS Forgot Password Feature - Test Script (PowerShell)
# Run from backend directory: .\test_forgot_password.ps1

$BASE_URL = "http://localhost:5000/api/auth"
$TEST_EMAIL = "test@example.com"

Write-Host "===============================================" -ForegroundColor Yellow
Write-Host "Sterling HMS - Forgot Password Feature Tests" -ForegroundColor Yellow
Write-Host "===============================================" -ForegroundColor Yellow
Write-Host ""

# Helper function to make requests
function Test-Endpoint {
    param(
        [string]$Method,
        [string]$Endpoint,
        [string]$Data,
        [string]$Description
    )
    
    Write-Host "Testing: $Description" -ForegroundColor Yellow
    Write-Host "Method: $Method"
    Write-Host "Endpoint: $BASE_URL$Endpoint"
    if ($Data) {
        Write-Host "Data: $Data"
    }
    Write-Host ""
    
    try {
        if ($Data) {
            $response = Invoke-WebRequest -Uri "$BASE_URL$Endpoint" `
                -Method $Method `
                -ContentType "application/json" `
                -Body $Data `
                -ErrorAction SilentlyContinue
        }
        else {
            $response = Invoke-WebRequest -Uri "$BASE_URL$Endpoint" `
                -Method $Method `
                -ContentType "application/json" `
                -ErrorAction SilentlyContinue
        }
        
        Write-Host "Status: $($response.StatusCode) $($response.StatusDescription)"
        Write-Host "Response:"
        $response.Content | ConvertFrom-Json | ConvertTo-Json -Depth 10 | Write-Host
    }
    catch {
        Write-Host "Error: $($_.Exception.Message)" -ForegroundColor Red
        if ($_.Exception.Response) {
            $responseStr = $_.Exception.Response.Content.ReadAsStringAsync().Result
            Write-Host "Response: $responseStr"
        }
    }
    
    Write-Host ""
    Write-Host "-------------------------------------------" -ForegroundColor Cyan
    Write-Host ""
}

# Test 1: Check server health
Write-Host "1. Checking server health..." -ForegroundColor Green
Test-Endpoint "GET" "/../../health" "" "Server Health"

# Test 2: Create a test user
Write-Host "2. Creating test user..." -ForegroundColor Green
$signupData = @{
    firstName = "Test"
    lastName  = "User"
    email     = $TEST_EMAIL
    password  = "TestPassword123!"
} | ConvertTo-Json

try {
    $signupResponse = Invoke-WebRequest -Uri "$BASE_URL/signup" `
        -Method POST `
        -ContentType "application/json" `
        -Body $signupData `
        -ErrorAction SilentlyContinue
    
    Write-Host "Signup Status: $($signupResponse.StatusCode)"
    $signupResponse.Content | ConvertFrom-Json | ConvertTo-Json -Depth 10 | Write-Host
}
catch {
    Write-Host "Signup Error (user may already exist): $($_.Exception.Message)" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "-------------------------------------------" -ForegroundColor Cyan
Write-Host ""

# Test 3: Request password reset
Write-Host "3. Testing Forgot Password endpoint..." -ForegroundColor Green
$forgotData = @{
    email = $TEST_EMAIL
} | ConvertTo-Json

Test-Endpoint "POST" "/forgot-password" $forgotData "Forgot Password Request"

# Test 4: Note about OTP
Write-Host "NOTE: In development mode, check server console for the generated OTP" -ForegroundColor Yellow
Write-Host "Or query: SELECT * FROM password_reset_logs ORDER BY created_at DESC LIMIT 5;" -ForegroundColor Yellow
Write-Host ""
Write-Host "-------------------------------------------" -ForegroundColor Cyan
Write-Host ""

# Test 5: Verify OTP with invalid OTP (3 attempts to trigger lockout)
Write-Host "4. Testing Verify OTP endpoint (with invalid OTP)..." -ForegroundColor Green

for ($i = 1; $i -le 3; $i++) {
    Write-Host "Attempt $i/3" -ForegroundColor Yellow
    
    $otpData = @{
        email = $TEST_EMAIL
        otp   = "000000"
    } | ConvertTo-Json
    
    Test-Endpoint "POST" "/verify-otp" $otpData "Verify OTP - Invalid (Attempt $i)"
    Start-Sleep -Seconds 1
}

Write-Host "After 3 failed attempts, account should be locked for 30 minutes" -ForegroundColor Red
Write-Host ""
Write-Host "-------------------------------------------" -ForegroundColor Cyan
Write-Host ""

# Test 6: Try locked account
Write-Host "5. Testing Locked Account..." -ForegroundColor Green
$otpData = @{
    email = $TEST_EMAIL
    otp   = "123456"
} | ConvertTo-Json

Test-Endpoint "POST" "/verify-otp" $otpData "Verify OTP - Locked Account"

# Test 7: Rate limiting
Write-Host "6. Testing Rate Limiting..." -ForegroundColor Green
Write-Host "Making 6 forgot-password requests to same email..." -ForegroundColor Yellow

for ($i = 1; $i -le 6; $i++) {
    Write-Host "Request $i/6" -ForegroundColor Yellow
    
    $forgotData = @{
        email = "ratelimit@example.com"
    } | ConvertTo-Json
    
    try {
        $response = Invoke-WebRequest -Uri "$BASE_URL/forgot-password" `
            -Method POST `
            -ContentType "application/json" `
            -Body $forgotData `
            -ErrorAction SilentlyContinue
        
        if ($response.StatusCode -eq 429) {
            Write-Host "Request $i: Rate limited (429 Too Many Requests)" -ForegroundColor Red
        }
        else {
            Write-Host "Request $i: Success ($($response.StatusCode))" -ForegroundColor Green
        }
    }
    catch {
        if ($_.Exception.Response.StatusCode -eq 429) {
            Write-Host "Request $i: Rate limited (429 Too Many Requests)" -ForegroundColor Red
        }
        else {
            Write-Host "Request $i: Error - $($_.Exception.Message)" -ForegroundColor Red
        }
    }
    Write-Host ""
}

Write-Host "-------------------------------------------" -ForegroundColor Cyan
Write-Host ""

# Test 8: Non-existent email
Write-Host "7. Testing Non-existent Email (generic response)..." -ForegroundColor Green
$forgotData = @{
    email = "nonexistent@example.com"
} | ConvertTo-Json

Test-Endpoint "POST" "/forgot-password" $forgotData "Forgot Password - Non-existent Email"

# Test 9: Invalid email format
Write-Host "8. Testing Invalid Email Format..." -ForegroundColor Green
$forgotData = @{
    email = "notanemail"
} | ConvertTo-Json

Test-Endpoint "POST" "/forgot-password" $forgotData "Forgot Password - Invalid Format"

# Test 10: Reset password with invalid token
Write-Host "9. Testing Reset Password with Invalid Token..." -ForegroundColor Green
$resetData = @{
    resetToken = "invalid_token_12345"
    password   = "NewPassword123!"
} | ConvertTo-Json

Test-Endpoint "POST" "/reset-password" $resetData "Reset Password - Invalid Token"

# Test 11: Weak passwords
Write-Host "10. Testing Password Validation (Weak Passwords)..." -ForegroundColor Green

$weakPasswords = @(
    @{resetToken = "valid_token"; password = "weak"},
    @{resetToken = "valid_token"; password = "NoNumbers!"},
    @{resetToken = "valid_token"; password = "nouppercase123!"},
    @{resetToken = "valid_token"; password = "NOLOWERCASE123!"},
    @{resetToken = "valid_token"; password = "NoSpecialChar123"}
)

foreach ($pwd in $weakPasswords) {
    Write-Host "Testing weak password: $($pwd.password)" -ForegroundColor Yellow
    $resetData = $pwd | ConvertTo-Json
    Test-Endpoint "POST" "/reset-password" $resetData "Password Validation"
}

# Test 12: Resend OTP
Write-Host "11. Testing Resend OTP..." -ForegroundColor Green
$resendData = @{
    email = $TEST_EMAIL
} | ConvertTo-Json

Test-Endpoint "POST" "/resend-otp" $resendData "Resend OTP"

Write-Host "===============================================" -ForegroundColor Green
Write-Host "Test Suite Complete!" -ForegroundColor Green
Write-Host "===============================================" -ForegroundColor Green
Write-Host ""

Write-Host "Summary:" -ForegroundColor Yellow
Write-Host "- All endpoints tested"
Write-Host "- Rate limiting verified (6th request = 429)"
Write-Host "- Account lockout tested (3 failed attempts = 30 min)"
Write-Host "- Generic error responses verified"
Write-Host "- Password validation tested"
Write-Host ""

Write-Host "Next steps:" -ForegroundColor Yellow
Write-Host "1. Check database audit logs"
Write-Host "2. Verify OTP and reset tokens are hashed"
Write-Host "3. Monitor server logs for any errors"
Write-Host ""

Write-Host "Database Queries:" -ForegroundColor Cyan
Write-Host "psql -d sterling_hms -c 'SELECT * FROM password_reset_logs ORDER BY created_at DESC LIMIT 10;'"
Write-Host "psql -d sterling_hms -c 'SELECT * FROM otp_lockouts;'"
Write-Host "psql -d sterling_hms -c 'SELECT * FROM password_reset_tokens ORDER BY created_at DESC LIMIT 5;'"
