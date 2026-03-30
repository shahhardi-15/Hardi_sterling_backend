-- Verify appointment slots exist and are available
-- Check total slots
SELECT COUNT(*) as total_slots_in_db FROM appointment_slots;

-- Check slots for today and tomorrow
SELECT 
    doctor_id, 
    slot_date, 
    COUNT(*) as slots_count,
    COUNT(CASE WHEN is_available = true THEN 1 END) as available_slots
FROM appointment_slots
WHERE slot_date IN (CURRENT_DATE, CURRENT_DATE + INTERVAL '1 day')
GROUP BY doctor_id, slot_date
ORDER BY slot_date, doctor_id
LIMIT 20;

-- Show sample available slots for first doctor today
SELECT 
    id,
    doctor_id,
    slot_date,
    time_slot,
    is_available
FROM appointment_slots
WHERE doctor_id = 1 
    AND slot_date = CURRENT_DATE
    AND is_available = true
ORDER BY time_slot
LIMIT 15;

-- Check if any appointments exist (should be empty after truncate unless new bookings)
SELECT COUNT(*) as total_appointments FROM appointments;
