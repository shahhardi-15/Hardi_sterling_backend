-- Check slots for April 4, 2026 and other dates
SELECT slot_date, COUNT(*) as total_slots_for_date
FROM appointment_slots
WHERE slot_date IN ('2026-04-04', '2026-03-30', '2026-03-31')
GROUP BY slot_date
ORDER BY slot_date;

-- Show sample slots for April 4
SELECT doctor_id, time_slot, is_available
FROM appointment_slots
WHERE slot_date = '2026-04-04' AND doctor_id = 1
ORDER BY time_slot;
