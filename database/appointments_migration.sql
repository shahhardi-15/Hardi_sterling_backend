-- Create doctors table
CREATE TABLE IF NOT EXISTS doctors (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  specialization VARCHAR(255) NOT NULL,
  email VARCHAR(255),
  phone VARCHAR(20),
  is_available BOOLEAN DEFAULT true,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create appointments table
CREATE TABLE IF NOT EXISTS appointments (
  id SERIAL PRIMARY KEY,
  patient_id INTEGER NOT NULL,
  doctor_id INTEGER NOT NULL,
  appointment_date DATE NOT NULL,
  time_slot VARCHAR(10) NOT NULL,
  reason VARCHAR(500),
  status VARCHAR(50) DEFAULT 'scheduled', -- scheduled, completed, cancelled, no-show
  notes TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (patient_id) REFERENCES users(id) ON DELETE CASCADE,
  FOREIGN KEY (doctor_id) REFERENCES doctors(id) ON DELETE RESTRICT,
  UNIQUE(doctor_id, appointment_date, time_slot)
);

-- Create appointment slots table for availability
CREATE TABLE IF NOT EXISTS appointment_slots (
  id SERIAL PRIMARY KEY,
  doctor_id INTEGER NOT NULL,
  slot_date DATE NOT NULL,
  time_slot VARCHAR(10) NOT NULL,
  is_available BOOLEAN DEFAULT true,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (doctor_id) REFERENCES doctors(id) ON DELETE CASCADE,
  UNIQUE(doctor_id, slot_date, time_slot)
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_appointments_patient_id ON appointments(patient_id);
CREATE INDEX IF NOT EXISTS idx_appointments_doctor_id ON appointments(doctor_id);
CREATE INDEX IF NOT EXISTS idx_appointments_date ON appointments(appointment_date);
CREATE INDEX IF NOT EXISTS idx_appointments_status ON appointments(status);
CREATE INDEX IF NOT EXISTS idx_slots_doctor_id ON appointment_slots(doctor_id);
CREATE INDEX IF NOT EXISTS idx_slots_date ON appointment_slots(slot_date);
CREATE INDEX IF NOT EXISTS idx_doctors_specialization ON doctors(specialization);

-- Insert sample doctors
INSERT INTO doctors (name, specialization, email, phone, is_available) VALUES
  ('Dr. John Smith', 'General Practitioner', 'john.smith@hospital.com', '+1234567890', true),
  ('Dr. Sarah Johnson', 'Cardiologist', 'sarah.johnson@hospital.com', '+1234567891', true),
  ('Dr. Michael Brown', 'Dermatologist', 'michael.brown@hospital.com', '+1234567892', true),
  ('Dr. Emily Davis', 'Neurologist', 'emily.davis@hospital.com', '+1234567893', true),
  ('Dr. Robert Wilson', 'Orthopedist', 'robert.wilson@hospital.com', '+1234567894', true)
ON CONFLICT DO NOTHING;

-- Insert sample available slots for next 30 days
INSERT INTO appointment_slots (doctor_id, slot_date, time_slot, is_available)
SELECT 
  d.id,
  CURRENT_DATE + i AS slot_date,
  slots.time_slot,
  true
FROM doctors d
CROSS JOIN LATERAL (
  SELECT generate_series(1, 30) AS i
) dates
CROSS JOIN LATERAL (
  SELECT unnest(ARRAY['09:00', '09:30', '10:00', '10:30', '11:00', '14:00', '14:30', '15:00', '15:30', '16:00']) AS time_slot
) slots
WHERE (CURRENT_DATE + i)::dow NOT IN (0, 6) -- Exclude Sundays and Saturdays
ON CONFLICT DO NOTHING;
