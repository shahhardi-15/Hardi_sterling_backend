-- Create appointments table
CREATE TABLE IF NOT EXISTS appointments (
    id SERIAL PRIMARY KEY,
    patient_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    doctor_id INTEGER NOT NULL REFERENCES doctors(id),
    appointment_date DATE NOT NULL,
    time_slot VARCHAR(10) NOT NULL,
    reason TEXT,
    notes TEXT,
    status VARCHAR(20) DEFAULT 'scheduled' CHECK (status IN ('scheduled', 'completed', 'cancelled', 'no-show')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(doctor_id, appointment_date, time_slot)
);

-- Create appointment_slots table
CREATE TABLE IF NOT EXISTS appointment_slots (
    id SERIAL PRIMARY KEY,
    doctor_id INTEGER NOT NULL REFERENCES doctors(id),
    slot_date DATE NOT NULL,
    time_slot VARCHAR(10) NOT NULL,
    is_available BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(doctor_id, slot_date, time_slot)
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_appointments_patient_id ON appointments(patient_id);
CREATE INDEX IF NOT EXISTS idx_appointments_doctor_id ON appointments(doctor_id);
CREATE INDEX IF NOT EXISTS idx_appointments_appointment_date ON appointments(appointment_date);
CREATE INDEX IF NOT EXISTS idx_appointment_slots_doctor_id ON appointment_slots(doctor_id);
CREATE INDEX IF NOT EXISTS idx_appointment_slots_slot_date ON appointment_slots(slot_date);
CREATE INDEX IF NOT EXISTS idx_appointment_slots_is_available ON appointment_slots(is_available);

-- Clear existing slots
TRUNCATE TABLE appointment_slots CASCADE;

-- Seed appointment slots for next 30 days with specific times
-- Morning: 10:00, 10:30, 11:00, 11:30, 12:00, 12:30, 13:00 (10 AM to 1 PM)
-- Evening: 18:00, 18:30, 19:00, 19:30, 20:00, 20:30, 21:00 (6 PM to 9 PM)
-- Total: 14 slots per day per doctor (ALL DAYS including weekends)
INSERT INTO appointment_slots (doctor_id, slot_date, time_slot, is_available)
SELECT 
    d.id,
    (CURRENT_DATE + INTERVAL '1 day' * (n.day_offset))::date AS slot_date,
    t.time_slot,
    true
FROM 
    doctors d,
    (SELECT generate_series(0, 29) AS day_offset) AS n,
    (SELECT UNNEST(ARRAY[
        '10:00', '10:30', '11:00', '11:30', '12:00', '12:30', '13:00',
        '18:00', '18:30', '19:00', '19:30', '20:00', '20:30', '21:00'
    ]) AS time_slot) AS t
WHERE 
    d.is_available = true
ON CONFLICT (doctor_id, slot_date, time_slot) 
DO NOTHING;

-- Verify slots were created
SELECT COUNT(*) as total_slots FROM appointment_slots;
SELECT doctor_id, COUNT(*) as slots_count FROM appointment_slots GROUP BY doctor_id ORDER BY doctor_id LIMIT 5;
