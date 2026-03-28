-- Create doctor authentication table
CREATE TABLE IF NOT EXISTS doctor_users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    specialization VARCHAR(255),
    phone VARCHAR(20),
    role VARCHAR(50) DEFAULT 'doctor',
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create index on email for faster lookups
CREATE INDEX IF NOT EXISTS idx_doctor_users_email ON doctor_users(email);
CREATE INDEX IF NOT EXISTS idx_doctor_users_role ON doctor_users(role);

-- Create doctor_patient_assignment table
CREATE TABLE IF NOT EXISTS doctor_patient_assignment (
    id SERIAL PRIMARY KEY,
    doctor_id INTEGER NOT NULL REFERENCES doctor_users(id) ON DELETE CASCADE,
    patient_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    assigned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_active BOOLEAN DEFAULT true,
    UNIQUE(doctor_id, patient_id)
);

CREATE INDEX IF NOT EXISTS idx_doctor_patient_assignment_doctor_id ON doctor_patient_assignment(doctor_id);
CREATE INDEX IF NOT EXISTS idx_doctor_patient_assignment_patient_id ON doctor_patient_assignment(patient_id);

-- Sample doctor data
INSERT INTO doctor_users (email, name, password_hash, specialization, phone) VALUES
  ('dr.smith@sterling.com', 'Dr. John Smith', '$2a$10$YourHashedPasswordHere1', 'General Medicine', '+1234567890'),
  ('dr.johnson@sterling.com', 'Dr. Sarah Johnson', '$2a$10$YourHashedPasswordHere2', 'Cardiology', '+1234567891'),
  ('dr.brown@sterling.com', 'Dr. Michael Brown', '$2a$10$YourHashedPasswordHere3', 'Dermatology', '+1234567892')
ON CONFLICT DO NOTHING;
