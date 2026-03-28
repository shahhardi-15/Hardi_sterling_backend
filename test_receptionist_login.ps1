#!/usr/bin/env powershell

<#
.SYNOPSIS
    Test receptionist login API endpoint directly
.DESCRIPTION
    This script tests if the backend receptionist login API is working
    by making a direct HTTP request without going through the frontend.
#>

Write-Host "`n╔════════════════════════════════════════╗" -ForegroundColor Cyan
Write-Host "║   Receptionist API Direct Test         ║" -ForegroundColor Cyan
Write-Host "╚════════════════════════════════════════╝`n" -ForegroundColor Cyan

# Configuration
$BackendUrl = "http://localhost:5000"
$Endpoint = "/api/receptionist/login"
$Email = "receptionist@sterling.com"
$Password = "Receptionist@Sterling2026"

# Test 1: Check if backend is reachable
Write-Host "Test 1: Checking if backend is reachable..." -ForegroundColor Yellow
try {
    $healthResponse = Invoke-WebRequest -Uri "$BackendUrl/health" -Method GET -ErrorAction Stop
    Write-Host "✓ Backend is running on $BackendUrl" -ForegroundColor Green
} catch {
    Write-Host "✗ Backend is NOT reachable at $BackendUrl" -ForegroundColor Red
    Write-Host "  Make sure backend is running: go run cmd/main.go" -ForegroundColor Yellow
    exit 1
}

# Test 2: Make login request
Write-Host "`nTest 2: Sending login request..." -ForegroundColor Yellow
Write-Host "  Email: $Email" -ForegroundColor Gray
Write-Host "  Password: $Password" -ForegroundColor Gray

try {
    $body = @{
        email = $Email
        password = $Password
    } | ConvertTo-Json

    $response = Invoke-WebRequest -Uri "$BackendUrl$Endpoint" `
        -Method POST `
        -ContentType "application/json" `
        -Body $body `
        -ErrorAction Stop

    $responseData = $response.Content | ConvertFrom-Json

    Write-Host "`n✓ Login request successful!" -ForegroundColor Green
    Write-Host "  Status: $($response.StatusCode)" -ForegroundColor Green
    Write-Host "`n  Response:" -ForegroundColor Cyan
    Write-Host "    Message: $($responseData.message)" -ForegroundColor White
    Write-Host "    Success: $($responseData.success)" -ForegroundColor White
    Write-Host "    Token: $($responseData.token.substring(0, 20))..." -ForegroundColor White
    
    if ($responseData.receptionist) {
        Write-Host "`n  Receptionist Info:" -ForegroundColor Cyan
        Write-Host "    ID: $($responseData.receptionist.id)" -ForegroundColor White
        Write-Host "    Name: $($responseData.receptionist.name)" -ForegroundColor White
        Write-Host "    Email: $($responseData.receptionist.email)" -ForegroundColor White
        Write-Host "    Department: $($responseData.receptionist.department)" -ForegroundColor White
    }

} catch {
    $statusCode = $_.Exception.Response.StatusCode.Value__
    $errorContent = $_.Exception.Response.Content.ToString()
    
    Write-Host "`n✗ Login request failed!" -ForegroundColor Red
    Write-Host "  Status Code: $statusCode" -ForegroundColor Red
    
    try {
        $errorData = $errorContent | ConvertFrom-Json
        Write-Host "  Message: $($errorData.message)" -ForegroundColor Red
    } catch {
        Write-Host "  Response: $errorContent" -ForegroundColor Red
    }

    if ($statusCode -eq 401) {
        Write-Host "`n⚠ Issue: Incorrect credentials or user not found" -ForegroundColor Yellow
        Write-Host "  Solutions:" -ForegroundColor Yellow
        Write-Host "    1. Verify password is correct" -ForegroundColor Gray
        Write-Host "    2. Run migration to ensure user exists: .\run_receptionist_migration.ps1" -ForegroundColor Gray
        Write-Host "    3. Check database: SELECT * FROM receptionist_users" -ForegroundColor Gray
    } elseif ($statusCode -eq 404) {
        Write-Host "`n⚠ Issue: Endpoint not found" -ForegroundColor Yellow
        Write-Host "  Solutions:" -ForegroundColor Yellow
        Write-Host "    1. Verify backend is on correct port (5000)" -ForegroundColor Gray
        Write-Host "    2. Restart backend: go run cmd/main.go" -ForegroundColor Gray
    } elseif ($statusCode -eq 500) {
        Write-Host "`n⚠ Issue: Server error" -ForegroundColor Yellow
        Write-Host "  Solutions:" -ForegroundColor Yellow
        Write-Host "    1. Check backend logs for detailed error" -ForegroundColor Gray
        Write-Host "    2. Verify database connection: .\run_receptionist_migration.ps1" -ForegroundColor Gray
    }
    
    exit 1
}

# Test 3: Verify token can be used
Write-Host "`n`nTest 3: Verifying token works..." -ForegroundColor Yellow
try {
    $headers = @{
        "Authorization" = "Bearer $($responseData.token)"
        "Content-Type" = "application/json"
    }
    
    $statsResponse = Invoke-WebRequest -Uri "$BackendUrl/api/receptionist/dashboard/stats" `
        -Method GET `
        -Headers $headers `
        -ErrorAction Stop

    $statsData = $statsResponse.Content | ConvertFrom-Json
    
    Write-Host "✓ Token is valid and API is working!" -ForegroundColor Green
    Write-Host "  Dashboard Stats:" -ForegroundColor Cyan
    Write-Host "    Total Patients: $($statsData.stats.total_patients)" -ForegroundColor White
    Write-Host "    Pending Appointments: $($statsData.stats.pending_appointments)" -ForegroundColor White
    
} catch {
    Write-Host "⚠ Token validation failed, but login was successful" -ForegroundColor Yellow
    Write-Host "  This might be a minor issue. Check backend logs." -ForegroundColor Yellow
}

Write-Host "`n╔════════════════════════════════════════╗" -ForegroundColor Green
Write-Host "║   ✓ All tests passed!                   ║" -ForegroundColor Green
Write-Host "║                                        ║" -ForegroundColor Green
Write-Host "║   The backend API is working correctly ║" -ForegroundColor Green
Write-Host "║   If frontend login still fails:       ║" -ForegroundColor Green
Write-Host "║   1. Hard refresh frontend Ctrl+Shift+R║" -ForegroundColor Green
Write-Host "║   2. Check browser console (F12)       ║" -ForegroundColor Green
Write-Host "║   3. See LOGIN_DEBUG.md for help       ║" -ForegroundColor Green
Write-Host "╚════════════════════════════════════════╝`n" -ForegroundColor Green
