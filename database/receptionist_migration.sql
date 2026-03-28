-- Create receptionist_users table
CREATE TABLE IF NOT EXISTS receptionist_users (
  id SERIAL PRIMARY KEY,
  email VARCHAR(255) UNIQUE NOT NULL,
  password_hash VARCHAR(255) NOT NULL,
  name VARCHAR(255) NOT NULL,
  phone VARCHAR(20),
  department VARCHAR(100),
  role VARCHAR(50) DEFAULT 'receptionist',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  is_active BOOLEAN DEFAULT true
);

-- Create index on email for faster lookups
CREATE INDEX IF NOT EXISTS idx_receptionist_users_email ON receptionist_users(email);

-- Insert default receptionist user (password: Receptionist@Sterling2026)
-- Password hash generated using bcrypt
INSERT INTO receptionist_users (email, password_hash, name, phone, department, role, is_active)
VALUES ('receptionist@sterling.com', '$2a$10$gklQi.hBlgJy7xcRf6Qi4ecABZQtBekCNKfuSSLjnh8YReDLlGkwK', 'Receptionist Sterling', '+1234567890', 'Front Desk', 'receptionist', true)
ON CONFLICT (email) DO UPDATE SET password_hash = '$2a$10$gklQi.hBlgJy7xcRf6Qi4ecABZQtBekCNKfuSSLjnh8YReDLlGkwK';

-- Add approval status columns to appointments table if they don't exist
ALTER TABLE appointments
ADD COLUMN IF NOT EXISTS approval_status VARCHAR(50) DEFAULT 'pending', -- pending, approved, rejected
ADD COLUMN IF NOT EXISTS approved_by INT,
ADD COLUMN IF NOT EXISTS approved_at TIMESTAMP,
ADD COLUMN IF NOT EXISTS rejection_reason TEXT;

-- Add patient information columns to appointments for better tracking
ALTER TABLE appointments
ADD COLUMN IF NOT EXISTS patient_first_name VARCHAR(100),
ADD COLUMN IF NOT EXISTS patient_last_name VARCHAR(100),
ADD COLUMN IF NOT EXISTS patient_email VARCHAR(255),
ADD COLUMN IF NOT EXISTS patient_phone VARCHAR(20);

-- Create patient_records table for comprehensive patient management
CREATE TABLE IF NOT EXISTS patient_records (
  id SERIAL PRIMARY KEY,
  user_id INT UNIQUE,
  first_name VARCHAR(100) NOT NULL,
  last_name VARCHAR(100) NOT NULL,
  email VARCHAR(255) NOT NULL,
  phone VARCHAR(20),
  date_of_birth DATE,
  gender VARCHAR(20),
  blood_type VARCHAR(10),
  address TEXT,
  city VARCHAR(100),
  state VARCHAR(100),
  postal_code VARCHAR(20),
  country VARCHAR(100),
  allergies TEXT,
  medical_conditions TEXT,
  current_medications TEXT,
  emergency_contact_name VARCHAR(255),
  emergency_contact_phone VARCHAR(20),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Create index on patient_records for faster lookups
CREATE INDEX IF NOT EXISTS idx_patient_records_user_id ON patient_records(user_id);
CREATE INDEX IF NOT EXISTS idx_patient_records_email ON patient_records(email);
CREATE INDEX IF NOT EXISTS idx_patient_records_phone ON patient_records(phone);

-- Create receptionist audit log table
CREATE TABLE IF NOT EXISTS receptionist_audit_logs (
  id SERIAL PRIMARY KEY,
  receptionist_id INT NOT NULL,
  action VARCHAR(255) NOT NULL,
  resource_type VARCHAR(100),
  resource_id INT,
  details JSONB,
  ip_address VARCHAR(45),
  user_agent TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (receptionist_id) REFERENCES receptionist_users(id)
);

-- Create index on receptionist audit logs
CREATE INDEX IF NOT EXISTS idx_receptionist_audit_logs_receptionist_id ON receptionist_audit_logs(receptionist_id);
CREATE INDEX IF NOT EXISTS idx_receptionist_audit_logs_created_at ON receptionist_audit_logs(created_at DESC);

-- Add foreign key constraint for approval tracking
ALTER TABLE appointments
ADD CONSTRAINT IF NOT EXISTS fk_appointments_approved_by 
  FOREIGN KEY (approved_by) REFERENCES receptionist_users(id);

-- Create index for appointment queries by receptionist
CREATE INDEX IF NOT EXISTS idx_appointments_approval_status ON appointments(approval_status);
CREATE INDEX IF NOT EXISTS idx_appointments_created_at ON appointments(created_at DESC);
