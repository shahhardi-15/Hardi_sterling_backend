# Patient Dashboard - Quick Start Guide

## 🚀 System Overview

The Patient Dashboard is a full-stack application that allows patients to:
- View available doctors and their specializations
- Book appointment slots with preferred doctors
- View their complete appointment history with filtering
- Cancel scheduled appointments
- View and manage their profile

## 📦 What Was Built

### Backend (Go) ✅
- **5 API endpoints** for patient appointments
- **Database models** for doctors, appointments, and slots
- **Repositories** for clean data access
- **Handlers** with authentication and validation
- **Database schema** with sample data and indexes

### Frontend (Vue.js) ✅
- **3 reusable components** (Book, History, Profile)
- **Pinia store** for state management
- **Responsive dashboard** with modern UI
- **API service** with automatic JWT handling
- **Form validation** and error handling

### Database (PostgreSQL) ✅
- **3 main tables**: doctors, appointments, appointment_slots
- **Sample data**: 5 doctors with pre-populated slots
- **Indexes** for optimal query performance
- **Constraints** for data integrity

---

## 🔧 60-Second Setup

### Step 1: Database Migration (Windows)
```powershell
cd F:\Hardi_sterling_backend

# Run the migration
& "C:\Program Files\PostgreSQL\18\bin\psql.exe" -U postgres -h localhost -d sterling_hms `
  -f "database\appointments_migration.sql"
```

Verify it worked:
```powershell
# Connect to database
& "C:\Program Files\PostgreSQL\18\bin\psql.exe" -U postgres -h localhost -d sterling_hms

# Check tables
\dt

# Check sample doctors
SELECT * FROM doctors;

# Exit
\q
```

### Step 2: Start Backend
```powershell
cd F:\Hardi_sterling_backend
go run cmd/main.go
```

Expected output:
```
Server is running on http://localhost:8080
Environment: development
```

### Step 3: Start Frontend
```powershell
cd F:\Hardi_Sterling_frontend\sterling-hms-frontend
npm run dev
```

Expected output:
```
  VITE v4.x.x  ready in xxx ms
  ➜  Local:   http://localhost:5173/
```

### Step 4: Open Browser
Navigate to: **http://localhost:5173**

---

## 🧪 Testing Workflow

### 1. Create Test Account
1. Click "Sign Up" link
2. Enter:
   - First Name: `John`
   - Last Name: `Doe`
   - Email: `john.doe@test.com`
   - Password: `Test1234!` (must have uppercase, lowercase, number)
3. Click "Sign Up"
4. You'll be redirected to Dashboard

### 2. Book an Appointment
1. In the **"Book Appointment"** section:
   - Select Doctor: `Dr. John Smith - General Practitioner`
   - Select Date: Pick any future date (e.g., 2-3 days from today)
   - Available Time Slots will appear below the date
   - Click a time slot (e.g., `09:00`)
   - Enter Reason: `Annual checkup`
   - Click `Book Appointment`
2. Success message appears: **"Appointment booked successfully!"**

### 3. View Appointment History
1. Scroll to **"Appointment History"** section
2. You should see:
   - Your newly booked appointment in the table
   - Date, Time, Doctor Name, Specialization
   - Status badge showing "Scheduled" in blue
   - A "Cancel" button available
3. Test Filter:
   - Select "Scheduled" from status filter
   - Only scheduled appointments should show
   - Select "Cancelled" - should be empty
   - Select "All Statuses" - shows all

### 4. Cancel an Appointment
1. In the Appointment History table
2. Find your booked appointment with "Scheduled" status
3. Click the "Cancel" button
4. Confirm the cancellation dialog
5. Success message appears
6. Status changes to "Cancelled" (red badge)

### 5. Check Your Profile
1. Scroll to **"Your Profile"** section (right sidebar)
2. See:
   - Avatar with initials (e.g., "JD")
   - Your name: `John Doe`
   - Email: `john.doe@test.com`
   - Statistics cards:
     - Total Appointments: `1` (the one you booked)
     - Upcoming: `0` (or `1` if still scheduled)
     - Completed: `0`
3. Click "Logout" to test logout

---

## 🧬 Testing Different Scenarios

### Scenario 1: Multiple Appointments
```
1. Book appointment with Dr. John Smith on April 1 @ 09:00
2. Book appointment with Dr. Sarah Johnson on April 5 @ 14:00
3. Book appointment with Dr. Michael Brown on April 10 @ 16:00
4. View history - should show all 3
5. Cancel one - check it shows as cancelled
6. Profile should show: Total: 3, Scheduled: 2, Completed: 0
```

### Scenario 2: Slot Management
```
1. Select Dr. Sarah Johnson (Cardiologist)
2. Select April 15
3. Try to book 09:00 slot twice
   - First attempt succeeds
   - Second attempt fails with "slot not available"
4. Try different time slot - should succeed
```

### Scenario 3: Date Validation
```
1. Try to select a past date in date picker
   - Should be disabled (grayed out)
2. Try to book for more than 30 days away
   - Should show "no slots available"
3. Try dates more than 30 days out - no slots appear
```

### Scenario 4: Doctor Filtering
```
1. View all doctors via /api/doctors
2. Each specialty has doctors:
   - General Practitioner: 1 doctor
   - Cardiologist: 1 doctor
   - Dermatologist: 1 doctor
   - Neurologist: 1 doctor
   - Orthopedist: 1 doctor
3. Slots show grouped by specialization in dropdown
```

---

## 🔌 API Testing (Using Postman/curl)

### Get Doctors (No Auth Needed)
```bash
curl -X GET "http://localhost:8080/api/doctors"

# Response:
{
  "message": "Doctors retrieved successfully",
  "doctors": [
    {
      "id": 1,
      "name": "Dr. John Smith",
      "specialization": "General Practitioner",
      "email": "john.smith@hospital.com",
      "phone": "+1234567890",
      "isAvailable": true,
      "createdAt": "2026-03-27T..."
    },
    ...
  ]
}
```

### Get Available Slots (No Auth Needed)
```bash
curl -X GET "http://localhost:8080/api/doctors/available-slots?doctorId=1"

# Response:
{
  "message": "Available slots retrieved successfully",
  "slots": [
    {
      "id": 1,
      "doctorId": 1,
      "slotDate": "2026-03-30",
      "timeSlot": "09:00",
      "isAvailable": true
    },
    ...
  ]
}
```

### Get Patient Profile (Auth Required)
```bash
# First, get token from signin
curl -X POST "http://localhost:8080/api/auth/signin" \
  -H "Content-Type: application/json" \
  -d '{"email":"john.doe@test.com","password":"Test1234!"}'

# Copy the "token" from response, then:
curl -X GET "http://localhost:8080/api/patient/profile" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"

# Response:
{
  "message": "Patient profile retrieved successfully",
  "profile": {
    "id": 1,
    "firstName": "John",
    "lastName": "Doe",
    "email": "john.doe@test.com"
  }
}
```

### Book Appointment (Auth Required)
```bash
curl -X POST "http://localhost:8080/api/appointments/book" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{
    "doctorId": 1,
    "appointmentDate": "2026-04-01",
    "timeSlot": "10:00",
    "reason": "Regular checkup"
  }'

# Response:
{
  "message": "Appointment booked successfully",
  "appointment": {
    "id": 1,
    "patientId": 1,
    "doctorId": 1,
    "appointmentDate": "2026-04-01",
    "timeSlot": "10:00",
    "reason": "Regular checkup",
    "status": "scheduled",
    "createdAt": "2026-03-27T..."
  }
}
```

### Get Appointment History (Auth Required)
```bash
curl -X GET "http://localhost:8080/api/appointments/history?page=1&limit=10" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"

# Response:
{
  "message": "Appointment history retrieved successfully",
  "appointments": [
    {
      "id": 1,
      "patientId": 1,
      "doctorId": 1,
      "appointmentDate": "2026-04-01",
      "timeSlot": "10:00",
      "reason": "Regular checkup",
      "status": "scheduled",
      "doctor": {
        "id": 1,
        "name": "Dr. John Smith",
        "specialization": "General Practitioner",
        ...
      },
      "createdAt": "2026-03-27T..."
    }
  ],
  "total": 1
}
```

### Cancel Appointment (Auth Required)
```bash
curl -X DELETE "http://localhost:8080/api/appointments/1" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"

# Response:
{
  "message": "Appointment cancelled successfully"
}
```

---

## 📊 Database Queries for Testing

### Check Doctors
```sql
SELECT id, name, specialization, is_available FROM doctors;
```

### Check Available Slots
```sql
SELECT doctor_id, slot_date, time_slot, is_available 
FROM appointment_slots 
WHERE slot_date >= CURRENT_DATE 
ORDER BY slot_date, time_slot;
```

### Check Booked Appointments
```sql
SELECT a.id, a.patient_id, d.name as doctor_name, a.appointment_date, 
       a.time_slot, a.status, a.reason
FROM appointments a
JOIN doctors d ON a.doctor_id = d.id
ORDER BY a.appointment_date DESC;
```

### Check Slot Availability After Booking
```sql
-- Should show false for booked slots
SELECT doctor_id, slot_date, time_slot, is_available
FROM appointment_slots
WHERE doctor_id = 1 AND slot_date = '2026-04-01' AND time_slot = '10:00';
```

---

## ⚠️ Troubleshooting

### Error: `UNIQUE violation on doctor_id, appointment_date, time_slot`
**Cause**: Trying to book an already taken slot
**Fix**: Select a different time slot

### Error: `Selected time slot is not available`
**Cause**: Slot may have been booked by another user
**Fix**: Refresh page and select a different slot

### Error: `Invalid date format. Use YYYY-MM-DD`
**Cause**: Frontend should handle this, but verify date format
**Fix**: Use YYYY-MM-DD format in API calls

### Error: `Appointment not found` when cancelling
**Cause**: Appointment ID doesn't exist or belongs to different user
**Fix**: Verify appointment ID from history

### Frontend shows `loading...` forever
**Cause**: Backend not running or API endpoint incorrect
**Fix**: 
1. Check backend is running on port 8080
2. Check browser DevTools Network tab
3. Verify API base URL in `src/api/auth.js`

### No doctors showing in dropdown
**Cause**: Migration not run or doctors table empty
**Fix**: 
1. Run migration: `psql -d sterling_hms -f database/appointments_migration.sql`
2. Verify: `SELECT COUNT(*) FROM doctors;` should be 5

### Can't cancel appointment
**Cause**: Only "scheduled" status appointments can be cancelled
**Fix**: Check appointment status is "scheduled" not "completed" or "cancelled"

---

## 📈 Performance Notes

- **First load**: May take a few seconds while loading doctors and initial appointments
- **Slot loading**: Dynamically loads when date changes (slight delay is normal)
- **Pagination**: History shows 10 items per page for performance
- **Caching**: Doctor list cached after first load

---

## 🔐 Security Features Verified

✅ JWT authentication on patient endpoints
✅ Patients can only view/manage their own appointments
✅ Slot availability check prevents double booking
✅ Past dates disabled in date picker
✅ CORS configured for frontend domain
✅ Password hashing on backend
✅ Input validation on all endpoints
✅ Database constraints prevent integrity violations

---

## 📝 Files Reference

| File | Purpose |
|------|---------|
| `database/appointments_migration.sql` | Database schema & sample data |
| `internal/models/models.go` | Data models (Doctor, Appointment, etc.) |
| `internal/repositories/appointment_repository.go` | Database queries |
| `internal/handlers/appointment_handler.go` | API handlers |
| `cmd/main.go` | Routes & server setup |
| `src/stores/appointment.js` | State management |
| `src/api/appointment.js` | API calls |
| `src/components/BookAppointment.vue` | Booking form |
| `src/components/AppointmentHistory.vue` | History table |
| `src/components/PatientProfile.vue` | Profile display |
| `src/views/Dashboard.vue` | Main dashboard |

---

## 🎯 Success Criteria

After following this guide, you should be able to:
- ✅ Sign up and login
- ✅ View available doctors grouped by specialization
- ✅ See time slots for selected doctor and date
- ✅ Book appointments successfully
- ✅ View appointment history with full details
- ✅ Filter appointments by status
- ✅ Cancel scheduled appointments
- ✅ View patient profile and statistics
- ✅ Logout properly

If you can do all of these, the Patient Dashboard is fully functional! 🎉
