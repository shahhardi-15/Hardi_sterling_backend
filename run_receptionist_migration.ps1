#!/usr/bin/env powershell

<#
.SYNOPSIS
    Runs the receptionist system migration on existing Sterling HMS database
.DESCRIPTION
    This script applies the receptionist_migration.sql to add receptionist user table
    and appointment approval workflow to an already-initialized Sterling HMS database.
    Use this if you've already run setup.ps1 but need to add receptionist functionality.
#>

$psqlPath = "C:\Program Files\PostgreSQL\18\bin\psql.exe"
$dbHost = "127.0.0.1"
$dbUser = "postgres"
$dbPort = "5432"
$dbName = "sterling"

# Check if psql exists
if (-not (Test-Path $psqlPath)) {
    Write-Host "ERROR: psql not found at $psqlPath" -ForegroundColor Red
    Write-Host "Please install PostgreSQL or update the path in this script" -ForegroundColor Yellow
    exit 1
}

# Test database connection
Write-Host "Testing database connection to $dbName..." -ForegroundColor Cyan
try {
    $output = & $psqlPath -U $dbUser -h $dbHost -p $dbPort -d $dbName -c "SELECT 1;" 2>&1
    if ($LASTEXITCODE -ne 0) {
        throw "Connection failed"
    }
    Write-Host "✓ Database connection successful" -ForegroundColor Green
} catch {
    Write-Host "ERROR: Could not connect to database" -ForegroundColor Red
    Write-Host "Please ensure PostgreSQL is running and the connection details are correct" -ForegroundColor Yellow
    exit 1
}

# Run receptionist migration
Write-Host "`nApplying receptionist system migration..." -ForegroundColor Cyan
& $psqlPath -U $dbUser -h $dbHost -p $dbPort -d $dbName -f "database/receptionist_migration.sql"

if ($LASTEXITCODE -eq 0) {
    Write-Host "`n✓ Migration completed successfully!" -ForegroundColor Green
    Write-Host "`nDefault receptionist credentials:" -ForegroundColor Green
    Write-Host "  Email: receptionist@sterling.com" -ForegroundColor White
    Write-Host "  Password: receptionist@123" -ForegroundColor White
    Write-Host "`nYou can now log in with these credentials at http://localhost:5173/login"
} else {
    Write-Host "`nERROR: Migration failed" -ForegroundColor Red
    Write-Host "Please check the error message above" -ForegroundColor Yellow
    exit 1
}
