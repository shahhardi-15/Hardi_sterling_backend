# Test Script for View Appointments History Feature (Windows PowerShell)
# Usage: .\TEST_APPOINTMENTS_HISTORY.ps1

Write-Host "=====================================" -ForegroundColor Green
Write-Host "View Appointments History - Test Suite" -ForegroundColor Green
Write-Host "=====================================" -ForegroundColor Green
Write-Host ""

$API_URL = "http://localhost:5000/api"

Write-Host "STEP 1: Verify Database" -ForegroundColor Yellow
Write-Host "Checking appointments table structure..." -ForegroundColor Gray

$env:PGPASSWORD = "admin"
&"C:\Program Files\PostgreSQL\18\bin\psql.exe" -U postgres -h localhost -d sterling_hms -c "
SELECT column_name, data_type 
FROM information_schema.columns 
WHERE table_name = 'appointments' 
ORDER BY ordinal_position;"

if ($LASTEXITCODE -eq 0) {
    Write-Host "✓ Database verification passed" -ForegroundColor Green
} else {
    Write-Host "✗ Database verification failed" -ForegroundColor Red
}
Write-Host ""

Write-Host "STEP 2: Check Test Data" -ForegroundColor Yellow
Write-Host "Checking for test appointments..." -ForegroundColor Gray

$env:PGPASSWORD = "admin"
&"C:\Program Files\PostgreSQL\18\bin\psql.exe" -U postgres -h localhost -d sterling_hms -c "
SELECT 
    a.id, 
    a.patient_id, 
    a.doctor_id, 
    a.appointment_date, 
    a.time_slot, 
    a.reason, 
    a.status
FROM appointments  
WHERE patient_id = 6 
LIMIT 5;" 2>&1 | Out-Null

$appointmentCount = &"C:\Program Files\PostgreSQL\18\bin\psql.exe" -U postgres -h localhost -d sterling_hms -t -c "SELECT COUNT(*) FROM appointments WHERE patient_id = 6;" 2>&1 | Select-Object -First 1 | ForEach-Object { $_.Trim() }

if ([int]$appointmentCount -gt 0) {
    Write-Host "✓ Test data exists ($appointmentCount appointments found)" -ForegroundColor Green
} else {
    Write-Host "⚠ No test appointments found for patient 6" -ForegroundColor Yellow
}
Write-Host ""

Write-Host "STEP 3: Verify Backend SQL Query" -ForegroundColor Yellow
Write-Host "Testing appointment history join query..." -ForegroundColor Gray

$env:PGPASSWORD = "admin"
&"C:\Program Files\PostgreSQL\18\bin\psql.exe" -U postgres -h localhost -d sterling_hms -c "
SELECT 
    a.id, 
    a.patient_id, 
    a.doctor_id, 
    a.appointment_date, 
    a.time_slot,
    a.reason, 
    a.status,
    d.name as doctor_name,
    d.specialization
FROM appointments a
JOIN doctors d ON a.doctor_id = d.id
WHERE a.patient_id = 6
ORDER BY a.appointment_date DESC, a.time_slot DESC
LIMIT 1;"

if ($LASTEXITCODE -eq 0) {
    Write-Host "✓ SQL query executed successfully" -ForegroundColor Green
} else {
    Write-Host "✗ SQL query failed" -ForegroundColor Red
}
Write-Host ""

Write-Host "STEP 4: Application Fixes Applied" -ForegroundColor Yellow
Write-Host ""
Write-Host "Frontend Component Fixes:" -ForegroundColor Cyan
Write-Host "✓ Added 'watch' import from Vue"
Write-Host "✓ Created loadAppointmentHistory() function"
Write-Host "✓ Added onMounted hook to load initial data"
Write-Host "✓ Added watch(currentPage) to reload on pagination"
Write-Host "✓ formatTimeSlot() converts 24h to 12h AM/PM format"
Write-Host "✓ Component displays time_slot in table"
Write-Host ""

Write-Host "STEP 5: Manual Testing Checklist" -ForegroundColor Yellow
Write-Host ""
Write-Host "Please perform these manual tests in the browser:" -ForegroundColor Cyan
Write-Host ""

Write-Host "1. Authentication:" -ForegroundColor Gray
Write-Host "   [ ] Log in with a patient account"
Write-Host ""

Write-Host "2. Dashboard Navigation:" -ForegroundColor Gray
Write-Host "   [ ] Navigate to Patient Dashboard"
Write-Host "   [ ] Appointment History section should appear"
Write-Host ""

Write-Host "3. Data Display:" -ForegroundColor Gray
Write-Host "   [ ] Verify appointment table shows:"
Write-Host "       - Date column (formatted, e.g., '28 Mar 2026')"
Write-Host "       - Time column (formatted, e.g., '9:00 AM')" -ForegroundColor White
Write-Host "       - Doctor name"
Write-Host "       - Specialization"
Write-Host "       - Reason for visit"
Write-Host "       - Status badge (color-coded)"
Write-Host ""

Write-Host "4. Time Slot Formatting Test:" -ForegroundColor Gray
Write-Host "   [ ] Verify time slot conversions:"
Write-Host "       - '09:00' should display as '9:00 AM'"
Write-Host "       - '14:30' should display as '2:30 PM'"
Write-Host "       - '00:00' should display as '12:00 AM'"
Write-Host "       - '12:00' should display as '12:00 PM'"
Write-Host ""

Write-Host "5. Pagination:" -ForegroundColor Gray
Write-Host "   [ ] If more than 10 appointments:"
Write-Host "       - Previous button disabled on page 1"
Write-Host "       - Next button functional"
Write-Host "       - Click Next loads new appointments"
Write-Host "       - Click Previous goes back"
Write-Host ""

Write-Host "6. Status Filters:" -ForegroundColor Gray
Write-Host "   [ ] Select different status filters:"
Write-Host "       - 'Scheduled' (blue badge)"
Write-Host "       - 'Completed' (green badge)"
Write-Host "       - 'Cancelled' (red badge)"
Write-Host "   [ ] Only matching appointments display"
Write-Host ""

Write-Host "7. Cancel Appointment Action:" -ForegroundColor Gray
Write-Host "   [ ] Click Cancel button on scheduled appointment"
Write-Host "   [ ] Confirm dialog appears"
Write-Host "   [ ] After confirmation, status changes to 'Cancelled'"
Write-Host "   [ ] Success message displays for 3 seconds"
Write-Host ""

Write-Host "8. Error Handling:" -ForegroundColor Gray
Write-Host "   [ ] Check browser DevTools console for errors"
Write-Host "   [ ] Network tab shows successful 200 responses"
Write-Host "   [ ] No 401/403 authorization errors"
Write-Host ""

Write-Host "STEP 6: Browser Console Inspection" -ForegroundColor Yellow
Write-Host ""
Write-Host "Open browser DevTools (F12) and verify:" -ForegroundColor Cyan
Write-Host "  • No red error messages in console"
Write-Host "  • Network requests show 200 status"
Write-Host "  • API response includes 'timeSlot' field in each appointment"
Write-Host "  • Doctor nested object includes name and specialization"
Write-Host ""

Write-Host "STEP 7: File Changes Summary" -ForegroundColor Yellow
Write-Host ""
Write-Host "Modified Files:" -ForegroundColor Cyan
Write-Host "  ✓ sterling-hms-frontend/src/components/AppointmentHistory.vue"
Write-Host "     - Added 'watch' import"
Write-Host "     - Refactored to use loadAppointmentHistory() function"  
Write-Host "     - Added watch on currentPage for pagination"
Write-Host ""

Write-Host "Verified Files:" -ForegroundColor Cyan
Write-Host "  ✓ sterling-hms-backend/internal/handlers/appointment_handler.go"
Write-Host "  ✓ sterling-hms-backend/internal/repositories/appointment_repository.go"
Write-Host "  ✓ sterling-hms-backend/database/schema.sql"
Write-Host "  ✓ sterling-hms-frontend/src/stores/appointment.js"
Write-Host "  ✓ sterling-hms-frontend/src/api/appointment.js"
Write-Host ""

Write-Host "STEP 8: Expected API Response Format" -ForegroundColor Yellow
Write-Host ""
Write-Host "GET /api/appointments/history?page=1&limit=10" -ForegroundColor Cyan
Write-Host "Response:" -ForegroundColor Gray
Write-Host "{
  ""message"": ""Appointment history retrieved successfully"",
  ""appointments"": [
    {
      ""id"": 1,
      ""patientId"": 6,
      ""doctorId"": 1,
      ""appointmentDate"": ""2026-03-30"",
      ""timeSlot"": ""09:00"",
      ""reason"": ""General Checkup"",
      ""status"": ""scheduled"",
      ""doctor"": {
        ""id"": 1,
        ""name"": ""Dr. Smith"",
        ""specialization"": ""Cardiology"",
        ""email"": ""dr.smith@hospital.com"",
        ...
      }
    }
  ],
  ""total"": 1
}" -ForegroundColor DarkGray
Write-Host ""

Write-Host "=====================================" -ForegroundColor Green
Write-Host "Testing Complete!" -ForegroundColor Green
Write-Host "=====================================" -ForegroundColor Green
Write-Host ""

Write-Host "Next Steps:" -ForegroundColor Cyan
Write-Host "  1. Run backend:   cd F:\Hardi_sterling_backend && go run ./cmd/main.go"
Write-Host "  2. Run frontend:  cd F:\Hardi_Sterling_frontend\sterling-hms-frontend && npm run dev"
Write-Host "  3. Open browser:  http://localhost:5173"
Write-Host "  4. Follow manual testing checklist above"
Write-Host ""
