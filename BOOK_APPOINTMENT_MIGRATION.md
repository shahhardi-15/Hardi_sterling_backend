# Book Appointment Form Enhancement - Migration Guide

## Overview
This document provides step-by-step instructions to apply the enhanced Book Appointment form changes to your system.

## Prerequisites
- Backend server running
- Frontend server running  
- Access to PostgreSQL database

## Backend Setup

### Step 1: Apply Database Migration
Run the migration SQL script to add new columns to the doctors table:

```bash
# From the backend directory
psql -U your_user -d your_database -f database/doctors_enhancement_migration.sql
```

Or manually run the migration in your PostgreSQL client:

```sql
-- Add new columns to doctors table
ALTER TABLE doctors
ADD COLUMN IF NOT EXISTS experience_years INTEGER DEFAULT 0,
ADD COLUMN IF NOT EXISTS qualification VARCHAR(255),
ADD COLUMN IF NOT EXISTS address TEXT;

-- Update sample doctors with realistic data
UPDATE doctors SET 
  experience_years = 15,
  qualification = 'MD - General Medicine',
  address = '123 Medical Center, City Hospital, Street 1, Main City'
WHERE name = 'Dr. John Smith';

UPDATE doctors SET 
  experience_years = 12,
  qualification = 'MD - Cardiology, Board Certified',
  address = '456 Cardiology Wing, City Hospital, Heart Street, Main City'
WHERE name = 'Dr. Sarah Johnson';

UPDATE doctors SET 
  experience_years = 10,
  qualification = 'MD - Dermatology',
  address = '789 Dermatology Clinic, City Hospital, Skin Lane, Main City'
WHERE name = 'Dr. Michael Brown';

UPDATE doctors SET 
  experience_years = 8,
  qualification = 'MD - Neurology',
  address = '321 Neurology Department, City Hospital, Brain Avenue, Main City'
WHERE name = 'Dr. Emily Davis';

UPDATE doctors SET 
  experience_years = 18,
  qualification = 'MD - Orthopedic Surgery, Board Certified',
  address = '654 Orthopedic Center, City Hospital, Bone Street, Main City'
WHERE name = 'Dr. Robert Wilson';
```

### Step 2: Verify Database Changes
Check that the columns were added:

```sql
SELECT id, name, specialization, experience_years, qualification, address
FROM doctors;
```

### Step 3: Restart Backend Server
Rebuild and restart the backend server to pick up the new endpoints:

```bash
go run cmd/main.go
```

The server should now serve these new endpoints:
- `GET /api/specializations`
- `GET /api/doctors/by-specialization?specialization={name}`

## Frontend Setup

### Step 1: No Installation Required
All frontend changes are already in place in the updated files:
- `src/stores/appointment.js` - Updated store with new actions
- `src/api/appointment.js` - Updated API calls
- `src/components/BookAppointment.vue` - Completely redesigned component

### Step 2: Verify Frontend is Running
Make sure the frontend development server is running:

```bash
cd sterling-hms-frontend
npm run dev
```

### Step 3: Test the New Form
1. Navigate to http://localhost:5174
2. Login to a patient account
3. Go to Dashboard
4. Click "Book Appointment"
5. Test the new form flow:
   - Select a specialization
   - See doctors for that specialty displayed as cards
   - Click a doctor to view full details
   - Select date and time slots
   - Add reason and optional notes
   - Submit booking

## API Response Examples

### GET /api/specializations
```json
{
  "message": "Specializations retrieved successfully",
  "specializations": [
    {
      "name": "Cardiologist"
    },
    {
      "name": "Dermatologist"
    },
    {
      "name": "General Practitioner"
    },
    {
      "name": "Neurologist"
    },
    {
      "name": "Orthopedist"
    }
  ]
}
```

### GET /api/doctors/by-specialization?specialization=Cardiologist
```json
{
  "message": "Doctors retrieved successfully",
  "doctors": [
    {
      "id": 2,
      "name": "Dr. Sarah Johnson",
      "specialization": "Cardiologist",
      "email": "sarah.johnson@hospital.com",
      "phone": "+1234567891",
      "experienceYears": 12,
      "qualification": "MD - Cardiology, Board Certified",
      "address": "456 Cardiology Wing, City Hospital, Heart Street, Main City",
      "isAvailable": true,
      "createdAt": "2024-01-15T10:00:00Z",
      "updatedAt": "2024-01-15T10:00:00Z"
    }
  ]
}
```

### POST /api/appointments/book (with notes)
Request:
```json
{
  "doctorId": 2,
  "appointmentDate": "2024-04-15",
  "timeSlot": "14:00",
  "reason": "Regular checkup for heart condition",
  "notes": "Please bring recent test results"
}
```

Response:
```json
{
  "message": "Appointment booked successfully",
  "appointment": {
    "id": 123,
    "patientId": 1,
    "doctorId": 2,
    "appointmentDate": "2024-04-15",
    "timeSlot": "14:00",
    "reason": "Regular checkup for heart condition",
    "status": "scheduled",
    "notes": "Please bring recent test results",
    "createdAt": "2024-03-27T10:30:00Z",
    "updatedAt": "2024-03-27T10:30:00Z"
  }
}
```

## Features Summary

### New UI Features
✅ Specialization dropdown filtering  
✅ Doctor cards with visual design  
✅ Doctor details panel after selection  
✅ Experience years display  
✅ Qualification information  
✅ Phone number display  
✅ Address display  
✅ Availability status indicators  
✅ Additional notes field (optional)  
✅ Step-by-step form guidance  
✅ Form validation before submission  

### New Backend Features
✅ Get specializations endpoint  
✅ Get doctors by specialization endpoint  
✅ Support for doctor experience_years  
✅ Support for doctor qualification  
✅ Support for doctor address  
✅ Support for appointment notes  

## Troubleshooting

### Database Migration Error
If you see "column already exists" error, it means the migration was already applied. No action needed.

### API Endpoint Not Found
Make sure backend is restarted after code changes. The new endpoints are:
- `/api/specializations`
- `/api/doctors/by-specialization`

### Frontend Form Not Loading
Check browser console for errors. Ensure backend is running and CORS is configured correctly.

### Doctor Data Not Showing
Verify doctors table has the new columns and sample data is populated. Run verification query above.

## Rollback (if needed)

To revert to the previous version:

1. Remove the new columns from doctors table:
```sql
ALTER TABLE doctors
DROP COLUMN IF EXISTS experience_years,
DROP COLUMN IF EXISTS qualification,
DROP COLUMN IF EXISTS address;
```

2. Reset frontend files to previous version
3. Restart both servers

## Support

If you encounter issues:
1. Check backend logs for errors
2. Verify database migration was successful
3. Ensure both frontend and backend servers are running
4. Check browser console for frontend errors
5. Verify CORS configuration in backend
