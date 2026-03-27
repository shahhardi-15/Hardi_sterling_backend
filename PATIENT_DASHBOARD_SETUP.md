# Patient Dashboard Implementation Guide

## Overview
This guide covers the complete implementation of a patient appointment booking dashboard for Sterling HMS, featuring Vue.js frontend, Go backend, and PostgreSQL database.

## Project Structure

### Backend Files Created/Modified
```
internal/
├── models/
│   └── models.go                     (Added: Doctor, Appointment, AppointmentSlot models)
├── handlers/
│   └── appointment_handler.go        (NEW: Appointment and patient handlers)
├── repositories/
│   └── appointment_repository.go     (NEW: Database operations)
└── middleware/
    └── auth.go                       (Existing: Used for auth protection)

database/
└── appointments_migration.sql        (NEW: Schema for appointments, doctors, slots)

cmd/
└── main.go                           (Modified: Added appointment routes)
```

### Frontend Files Created
```
src/
├── stores/
│   └── appointment.js                (NEW: Pinia store for appointment state)
├── api/
│   └── appointment.js                (NEW: API service for appointment endpoints)
├── components/
│   ├── BookAppointment.vue           (NEW: Booking form component)
│   ├── AppointmentHistory.vue        (NEW: History table component)
│   └── PatientProfile.vue            (NEW: Profile display component)
└── views/
    └── Dashboard.vue                 (Modified: Updated with 3-section layout)
```

## Database Setup

### 1. Run the Migration
```bash
cd F:\Hardi_sterling_backend

# Use psql to execute the migration
psql -U postgres -h localhost -d sterling_hms -f database/appointments_migration.sql
```

### 2. What Gets Created
- **doctors**: Stores doctor information with specializations
- **appointments**: Patient-doctor appointment records
- **appointment_slots**: Available time slots for doctors
- Sample doctors and slots are automatically inserted

### 3. Sample Data
The migration automatically inserts:
- 5 sample doctors with different specializations
- Appointment slots for next 30 days (excluding weekends)
- 10 time slots per day per doctor (09:00-16:00)

## Backend Setup

### 1. Dependencies
Ensure you have these in your `go.mod`:
```go
require (
    github.com/gin-gonic/gin v1.x.x
    github.com/golang-jwt/jwt/v5 v5.x.x
    github.com/lib/pq v1.x.x
)
```

### 2. Database Connection
The code uses existing `config.DB` from your config package. Ensure your database is running:
```bash
# Start PostgreSQL (Windows)
psql -U postgres -h localhost -d sterling_hms
```

### 3. API Endpoints

#### Public Endpoints
```
GET    /api/doctors                    # Get all available doctors
GET    /api/doctors/available-slots    # Get slots for a doctor
```

#### Protected Endpoints (require JWT auth)
```
GET    /api/patient/profile            # Get patient's profile
GET    /api/appointments/history       # Get patient's appointments (paginated)
POST   /api/appointments/book          # Create new appointment
DELETE /api/appointments/{id}          # Cancel appointment
```

### 4. Request/Response Examples

#### Book Appointment
```json
POST /api/appointments/book
{
  "doctorId": 1,
  "appointmentDate": "2026-04-15",
  "timeSlot": "10:00",
  "reason": "Regular checkup"
}

Response:
{
  "message": "Appointment booked successfully",
  "appointment": {
    "id": 1,
    "patientId": 2,
    "doctorId": 1,
    "appointmentDate": "2026-04-15",
    "timeSlot": "10:00",
    "reason": "Regular checkup",
    "status": "scheduled",
    "createdAt": "2026-03-27T..."
  }
}
```

#### Get Available Slots
```json
GET /api/doctors/available-slots?doctorId=1

Response:
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

### 5. Error Handling
The handlers implement comprehensive error handling:
- **400**: Invalid request data
- **401**: Unauthorized (missing/invalid token)
- **403**: Forbidden (accessing others' data)
- **404**: Resource not found
- **409**: Conflict (slot already booked)
- **500**: Server error

## Frontend Setup

### 1. Install Dependencies
```bash
cd F:\Hardi_Sterling_frontend\sterling-hms-frontend

npm install pinia  # If not already installed
```

### 2. Store Configuration
The Pinia store (`appointment.js`) manages:
- Patient profile data
- Appointment history
- Available doctors
- Available time slots
- Form state and loading

### 3. API Configuration
Update your API base URL in `src/api/auth.js`:
```javascript
const api = axios.create({
  baseURL: 'http://localhost:8080/api',  // Adjust to match backend
  headers: {
    'Content-Type': 'application/json'
  }
})
```

### 4. Component Features

#### BookAppointment Component
- Doctor selection (grouped by specialization)
- Date picker (prevents past dates)
- Dynamic time slot loading
- Validation for all fields
- Success/error feedback

#### AppointmentHistory Component
- Table view of all appointments
- Status filtering (scheduled, completed, cancelled, no-show)
- Pagination support
- Cancel appointment action (only for scheduled)
- Appointment details with doctor info

#### PatientProfile Component
- Patient information display
- Avatar with initials
- Statistics (total, upcoming, completed appointments)
- Logout functionality

### 5. Styling
All components use scoped CSS with:
- Responsive design (mobile, tablet, desktop)
- Color scheme: Green (#4CAF50) for primary actions
- Consistent spacing and typography
- Smooth animations and transitions

## Running the Application

### Backend
```bash
cd F:\Hardi_sterling_backend

# Run with auto-reload
go run cmd/main.go

# Or build and run
go build -o sterling-hms cmd/main.go
./sterling-hms
```

### Frontend
```bash
cd F:\Hardi_Sterling_frontend\sterling-hms-frontend

# Development server
npm run dev

# Production build
npm run build
```

Access the dashboard at: `http://localhost:5173/dashboard`

## Security Features

### Backend
1. **JWT Authentication**
   - All patient endpoints require valid token
   - Token contains user ID and email
   - Configurable expiration

2. **Data Ownership Validation**
   - Patients can only access their own appointments
   - Cancel only allowed for apt status = 'scheduled'
   - Doctor and slot validation before booking

3. **Slot Management**
   - Slots marked unavailable after booking
   - Slots freed up on cancellation
   - UNIQUE constraint prevents double booking

### Frontend
1. **Auth Middleware**
   - Route protection via `requiresAuth` meta field
   - Token stored in auth store
   - Auto-logout on token expiration

2. **Input Validation**
   - Client-side validation before submission
   - Date picker prevents past dates
   - Required fields enforced

## Testing

### Manual Testing Steps

1. **Create Test User**
   - Sign up with email: `patient@example.com`, password: `Test1234!`

2. **Test Booking**
   - Navigate to Dashboard
   - Select a doctor (e.g., Dr. John Smith - General Practitioner)
   - Pick a future date
   - Select available time slot (e.g., 09:00)
   - Enter reason for visit
   - Click "Book Appointment"

3. **View History**
   - Check "Appointment History" section
   - Should show the newly booked appointment
   - Verify status is "scheduled"
   - Test filter by status

4. **Cancel Appointment**
   - Click "Cancel" button on a scheduled appointment
   - Confirm cancellation
   - Status should change to "cancelled"

5. **Test Profile**
   - View patient profile card
   - Verify appointment statistics match history
   - Test logout functionality

### Automated Testing

Backend unit tests for repositories:
```bash
go test ./internal/repositories/... -v
```

Frontend unit tests:
```bash
npm run test
```

## Troubleshooting

### Common Issues

**1. "No available slots" error**
- Ensure migration was run successfully
- Check appointment_slots table has data
- Verify dates are in YYYY-MM-DD format

**2. "Cannot cancel appointment"**
- Only scheduled appointments can be cancelled
- Check appointment status in database
- Ensure correct appointment ID

**3. API 500 error**
- Check backend logs for database errors
- Verify database connection in config
- Check JWT secret is configured

**4. Frontend shows "loading" indefinitely**
- Check network tab in DevTools
- Verify API endpoint is correct
- Check token expiration

**5. CORS errors**
- Verify CORS middleware in main.go includes frontend URL
- Check request headers match allowed origins

## Extending the System

### Adding Features
1. **Email Notifications**
   - Send confirmation emails on booking
   - Reminder emails before appointment

2. **Appointment Rescheduling**
   - Allow changing date/time for scheduled appointments
   - Implement PUT /api/appointments/{id} handler

3. **Doctor Ratings**
   - Add rating field to appointments
   - Display average rating per doctor

4. **Administrative Features**
   - Doctor availability schedule management
   - Appointment history analytics

## Performance Optimization

### Database
- Indexes on frequently queried fields (patient_id, doctor_id, date)
- Connection pooling configured in config.DB
- Query pagination to limit results

### Frontend
- Lazy load components with Vue Router
- Debounce API calls on search
- Cache doctor list after first load
- Pagination for appointment history

## Deployment Checklist

- [ ] Database migration applied
- [ ] Environment variables configured (.env)
- [ ] CORS origins updated for production
- [ ] JWT secret configured securely
- [ ] Frontend built for production (`npm run build`)
- [ ] Backend compiled with optimizations
- [ ] Database backups configured
- [ ] Error logging implemented
- [ ] SSL/TLS certificates configured
- [ ] Load testing completed

## Support

For issues or questions:
1. Check the troubleshooting section above
2. Review backend logs: `go run cmd/main.go`
3. Check frontend DevTools console
4. Verify database connection: `psql -U postgres -h localhost -d sterling_hms`
