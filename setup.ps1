#!/usr/bin/env powershell

$psqlPath = "C:\Program Files\PostgreSQL\18\bin\psql.exe"
$dbHost = "127.0.0.1"
$dbUser = "postgres"
$dbPort = "5432"

Write-Host "Setting up PostgreSQL database..."

# Create the sterling database
$createDbSql = @"
CREATE DATABASE sterling;
"@

$createDbSql | & $psqlPath -U $dbUser -h $dbHost -p $dbPort -d template1

# Initialize database schema
Write-Host "Initializing database schema..."
& $psqlPath -U $dbUser -h $dbHost -p $dbPort -d sterling -f "database/schema.sql"

# Run receptionist system migration
Write-Host "Setting up receptionist system..."
& $psqlPath -U $dbUser -h $dbHost -p $dbPort -d sterling -f "database/receptionist_migration.sql"

Write-Host "Database setup complete!"
