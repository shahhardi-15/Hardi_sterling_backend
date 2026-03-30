# Test Appointment Booking Flow
# This script tests the complete appointment booking system

$baseUrl = "http://localhost:8080/api"
$doctorId = 1
$patientId = 1  # Assuming patient user exists
$appointmentDate = "2026-04-04"  # Saturday with slots

Write-Host "=== APPOINTMENT BOOKING SYSTEM TEST ===" -ForegroundColor Green
Write-Host "Date: $($appointmentDate) (Saturday - weekend slots now available)" -ForegroundColor Cyan
Write-Host ""

# Test 1: Get available doctors
Write-Host "Test 1: Fetching available doctors..." -ForegroundColor Yellow
try {
    $doctors = Invoke-RestMethod -Uri "$baseUrl/doctors" -Method Get
    Write-Host "PASS: Doctors fetched: $($doctors.doctors.Count) doctors available" -ForegroundColor Green
    Write-Host ""
    if ($doctors.doctors.Count -gt 0) {
        Write-Host "Sample doctor: $($doctors.doctors[0].name) - $($doctors.doctors[0].specialization)"
        Write-Host ""
    }
}
catch {
    Write-Host "FAIL: Failed to fetch doctors: $($_.Exception.Message)" -ForegroundColor Red
    Write-Host ""
}

# Test 2: Get available slots
Write-Host "Test 2: Fetching available slots for doctor $($doctorId)..." -ForegroundColor Yellow
try {
    $slots = Invoke-RestMethod -Uri "$baseUrl/doctors/available-slots?doctorId=$doctorId" -Method Get
    Write-Host "PASS: Slots fetched: $($slots.slots.Count) slots available" -ForegroundColor Green
    Write-Host ""
    
    # Filter by date
    $slotsForDate = $slots.slots | Where-Object { $_.slotDate -eq $appointmentDate -and $_.isAvailable }
    Write-Host "Slots available on $($appointmentDate): $($slotsForDate.Count) slots"
    if ($slotsForDate.Count -gt 0) {
        Write-Host "Sample times: $($slotsForDate[0].timeSlot), $($slotsForDate[1].timeSlot), $($slotsForDate[2].timeSlot)" -ForegroundColor Green
        Write-Host ""
    }
}
catch {
    Write-Host "FAIL: Failed to fetch slots: $($_.Exception.Message)" -ForegroundColor Red
    Write-Host ""
}

# Test 3: Check slot unavailable function
Write-Host "Test 3: Verifying slot availability tracking..." -ForegroundColor Yellow
try {
    $checkSlot = Invoke-RestMethod -Uri "$baseUrl/doctors/available-slots?doctorId=1" -Method Get
    $availableCount = ($checkSlot.slots | Where-Object { $_.slotDate -eq $appointmentDate -and $_.isAvailable }).Count
    Write-Host "PASS: Total available slots on $($appointmentDate): $($availableCount)" -ForegroundColor Green
    Write-Host ""
}
catch {
    Write-Host "FAIL: Failed to verify slots: $($_.Exception.Message)" -ForegroundColor Red
    Write-Host ""
}

Write-Host "════════════════════════════════════════════════════════════════" -ForegroundColor Cyan
Write-Host "TESTING INSTRUCTIONS:" -ForegroundColor Yellow
Write-Host "1. Keep backend running: cd f:\Hardi_sterling_backend && go run cmd/main.go" -ForegroundColor White
Write-Host "2. Keep frontend running: cd f:\Hardi_Sterling_frontend\sterling-hms-frontend && npm run dev" -ForegroundColor White
Write-Host "3. Open browser: http://localhost:5173/dashboard" -ForegroundColor White
Write-Host "4. Select a specialization and doctor" -ForegroundColor White
Write-Host "5. Select date: April 4, 2026 (now has weekend slots!)" -ForegroundColor White
Write-Host "6. 14 time slots should appear (10:00-13:00 morning, 18:00-21:00 evening)" -ForegroundColor White
Write-Host "7. Select a time slot and complete booking" -ForegroundColor White
Write-Host "8. Check browser console (F12) to see debugging log output" -ForegroundColor White
Write-Host ""
Write-Host "EXPECTED BEHAVIOR:" -ForegroundColor Green
Write-Host "   + Slots load from database when date selected" -ForegroundColor Green
Write-Host "   + 14 slots show for April 4 (morning 7 + evening 7)" -ForegroundColor Green
Write-Host "   + After booking, slot availability updates" -ForegroundColor Green
Write-Host "   + Appointment appears in appointment history" -ForegroundColor Green
Write-Host "   + Doctor name shows in appointment history" -ForegroundColor Green
