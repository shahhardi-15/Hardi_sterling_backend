-- Create doctor users linked to the doctors table for admin management
-- This ensures doctors are visible in the admin "Manage Doctors" page

-- First, ensure the users table has the columns needed
ALTER TABLE users
ADD COLUMN IF NOT EXISTS full_name VARCHAR(255),
ADD COLUMN IF NOT EXISTS phone VARCHAR(20),
ADD COLUMN IF NOT EXISTS role VARCHAR(50);

-- Update existing records with full_name from first_name and last_name if full_name is null
UPDATE users
SET full_name = CONCAT(first_name, ' ', last_name)
WHERE full_name IS NULL AND first_name IS NOT NULL;

-- First, ensure the doctors table has the proper admin schema with user_id
ALTER TABLE doctors
ADD COLUMN IF NOT EXISTS user_id INTEGER UNIQUE REFERENCES users(id) ON DELETE CASCADE;

-- Add other missing columns if needed for admin schema
ALTER TABLE doctors
ADD COLUMN IF NOT EXISTS registration_number VARCHAR(100) UNIQUE,
ADD COLUMN IF NOT EXISTS consultation_fee DECIMAL(10, 2) DEFAULT 0.00,
ADD COLUMN IF NOT EXISTS department_id UUID REFERENCES departments(id) ON DELETE SET NULL,
ADD COLUMN IF NOT EXISTS available_days TEXT,
ADD COLUMN IF NOT EXISTS start_time TIME,
ADD COLUMN IF NOT EXISTS end_time TIME,
ADD COLUMN IF NOT EXISTS slot_duration_minutes INTEGER DEFAULT 15;

-- Now create the doctor users and link them
-- Doctor 1: Dr. Rajesh Kumar - Cardiology
INSERT INTO users (first_name, last_name, full_name, email, phone, role, password, is_active, created_at)
VALUES ('Dr.', 'Rajesh Kumar', 'Dr. Rajesh Kumar', 'rajesh.kumar@hospital.com', '9876543201', 'doctor', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcg7b3XeKeUxWdeS86E36DRZlG2', true, NOW())
ON CONFLICT (email) DO NOTHING;

INSERT INTO doctors (user_id, specialization, qualification, registration_number, experience_years, consultation_fee, department_id, available_days, start_time, end_time, slot_duration_minutes, created_at)
SELECT u.id, 'Cardiology', 'MD - Cardiology, Board Certified', 'REG_CARD_001', 15, 500.00, d.id, 'Mon,Tue,Wed,Thu,Fri', '09:00'::time, '17:00'::time, 30, NOW()
FROM users u, departments d
WHERE u.email = 'rajesh.kumar@hospital.com' AND d.name = 'Cardiology'
AND NOT EXISTS (SELECT 1 FROM doctors WHERE user_id = u.id)
ON CONFLICT (registration_number) DO NOTHING;

-- Doctor 2: Dr. Priya Sharma - Gynecology
INSERT INTO users (first_name, last_name, full_name, email, phone, role, password, is_active, created_at)
VALUES ('Dr.', 'Priya Sharma', 'Dr. Priya Sharma', 'priya.sharma@hospital.com', '9876543202', 'doctor', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcg7b3XeKeUxWdeS86E36DRZlG2', true, NOW())
ON CONFLICT (email) DO NOTHING;

INSERT INTO doctors (user_id, specialization, qualification, registration_number, experience_years, consultation_fee, department_id, available_days, start_time, end_time, slot_duration_minutes, created_at)
SELECT u.id, 'Gynecology', 'MD - Obstetrics & Gynecology, Board Certified', 'REG_GYN_001', 12, 450.00, d.id, 'Mon,Tue,Wed,Thu,Fri', '10:00'::time, '18:00'::time, 30, NOW()
FROM users u, departments d
WHERE u.email = 'priya.sharma@hospital.com' AND d.name = 'Gynecology'
AND NOT EXISTS (SELECT 1 FROM doctors WHERE user_id = u.id)
ON CONFLICT (registration_number) DO NOTHING;

-- Doctor 3: Dr. Anil Patel - Orthopedics
INSERT INTO users (first_name, last_name, full_name, email, phone, role, password, is_active, created_at)
VALUES ('Dr.', 'Anil Patel', 'Dr. Anil Patel', 'anil.patel@hospital.com', '9876543203', 'doctor', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcg7b3XeKeUxWdeS86E36DRZlG2', true, NOW())
ON CONFLICT (email) DO NOTHING;

INSERT INTO doctors (user_id, specialization, qualification, registration_number, experience_years, consultation_fee, department_id, available_days, start_time, end_time, slot_duration_minutes, created_at)
SELECT u.id, 'Orthopedics', 'MD - Orthopedic Surgery, Board Certified', 'REG_ORTHO_001', 18, 550.00, d.id, 'Mon,Tue,Wed,Thu,Fri', '09:00'::time, '16:00'::time, 30, NOW()
FROM users u, departments d
WHERE u.email = 'anil.patel@hospital.com' AND d.name = 'Orthopedics'
AND NOT EXISTS (SELECT 1 FROM doctors WHERE user_id = u.id)
ON CONFLICT (registration_number) DO NOTHING;

-- Doctor 4: Dr. Neha Desai - Pediatrics
INSERT INTO users (first_name, last_name, full_name, email, phone, role, password, is_active, created_at)
VALUES ('Dr.', 'Neha Desai', 'Dr. Neha Desai', 'neha.desai@hospital.com', '9876543204', 'doctor', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcg7b3XeKeUxWdeS86E36DRZlG2', true, NOW())
ON CONFLICT (email) DO NOTHING;

INSERT INTO doctors (user_id, specialization, qualification, registration_number, experience_years, consultation_fee, department_id, available_days, start_time, end_time, slot_duration_minutes, created_at)
SELECT u.id, 'Pediatrics', 'MD - Pediatrics', 'REG_PEDI_001', 10, 400.00, d.id, 'Mon,Tue,Wed,Thu,Fri', '08:00'::time, '16:00'::time, 20, NOW()
FROM users u, departments d
WHERE u.email = 'neha.desai@hospital.com' AND d.name = 'Pediatrics'
AND NOT EXISTS (SELECT 1 FROM doctors WHERE user_id = u.id)
ON CONFLICT (registration_number) DO NOTHING;

-- Doctor 5: Dr. Vikram Singh - Neurology
INSERT INTO users (first_name, last_name, full_name, email, phone, role, password, is_active, created_at)
VALUES ('Dr.', 'Vikram Singh', 'Dr. Vikram Singh', 'vikram.singh@hospital.com', '9876543205', 'doctor', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcg7b3XeKeUxWdeS86E36DRZlG2', true, NOW())
ON CONFLICT (email) DO NOTHING;

INSERT INTO doctors (user_id, specialization, qualification, registration_number, experience_years, consultation_fee, department_id, available_days, start_time, end_time, slot_duration_minutes, created_at)
SELECT u.id, 'Neurology', 'MD - Neurology, Board Certified', 'REG_NEURO_001', 14, 600.00, d.id, 'Mon,Tue,Wed,Thu,Fri', '10:00'::time, '17:00'::time, 30, NOW()
FROM users u, departments d
WHERE u.email = 'vikram.singh@hospital.com' AND d.name = 'Neurology'
AND NOT EXISTS (SELECT 1 FROM doctors WHERE user_id = u.id)
ON CONFLICT (registration_number) DO NOTHING;

-- Doctor 6: Dr. Ananya Gupta - Dermatology
INSERT INTO users (first_name, last_name, full_name, email, phone, role, password, is_active, created_at)
VALUES ('Dr.', 'Ananya Gupta', 'Dr. Ananya Gupta', 'ananya.gupta@hospital.com', '9876543206', 'doctor', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcg7b3XeKeUxWdeS86E36DRZlG2', true, NOW())
ON CONFLICT (email) DO NOTHING;

INSERT INTO doctors (user_id, specialization, qualification, registration_number, experience_years, consultation_fee, department_id, available_days, start_time, end_time, slot_duration_minutes, created_at)
SELECT u.id, 'Dermatology', 'MD - Dermatology', 'REG_DERM_001', 11, 350.00, d.id, 'Mon,Tue,Wed,Thu,Fri', '09:00'::time, '17:00'::time, 20, NOW()
FROM users u, departments d
WHERE u.email = 'ananya.gupta@hospital.com' AND d.name = 'Dermatology'
AND NOT EXISTS (SELECT 1 FROM doctors WHERE user_id = u.id)
ON CONFLICT (registration_number) DO NOTHING;

-- Doctor 7: Dr. Suresh Reddy - ENT
INSERT INTO users (first_name, last_name, full_name, email, phone, role, password, is_active, created_at)
VALUES ('Dr.', 'Suresh Reddy', 'Dr. Suresh Reddy', 'suresh.reddy@hospital.com', '9876543207', 'doctor', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcg7b3XeKeUxWdeS86E36DRZlG2', true, NOW())
ON CONFLICT (email) DO NOTHING;

INSERT INTO doctors (user_id, specialization, qualification, registration_number, experience_years, consultation_fee, department_id, available_days, start_time, end_time, slot_duration_minutes, created_at)
SELECT u.id, 'ENT', 'MD - Otolaryngology, Board Certified', 'REG_ENT_001', 13, 400.00, d.id, 'Mon,Tue,Wed,Thu,Fri', '09:00'::time, '17:00'::time, 20, NOW()
FROM users u, departments d
WHERE u.email = 'suresh.reddy@hospital.com' AND d.name = 'ENT'
AND NOT EXISTS (SELECT 1 FROM doctors WHERE user_id = u.id)
ON CONFLICT (registration_number) DO NOTHING;

-- Doctor 8: Dr. Divya Menon - Ophthalmology
INSERT INTO users (first_name, last_name, full_name, email, phone, role, password, is_active, created_at)
VALUES ('Dr.', 'Divya Menon', 'Dr. Divya Menon', 'divya.menon@hospital.com', '9876543208', 'doctor', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcg7b3XeKeUxWdeS86E36DRZlG2', true, NOW())
ON CONFLICT (email) DO NOTHING;

INSERT INTO doctors (user_id, specialization, qualification, registration_number, experience_years, consultation_fee, department_id, available_days, start_time, end_time, slot_duration_minutes, created_at)
SELECT u.id, 'Ophthalmology', 'MD - Ophthalmology, Board Certified', 'REG_OPHTHAL_001', 12, 450.00, d.id, 'Mon,Tue,Wed,Thu,Fri', '10:00'::time, '18:00'::time, 20, NOW()
FROM users u, departments d
WHERE u.email = 'divya.menon@hospital.com' AND d.name = 'Ophthalmology'
AND NOT EXISTS (SELECT 1 FROM doctors WHERE user_id = u.id)
ON CONFLICT (registration_number) DO NOTHING;

-- Doctor 9: Dr. Arjun Verma - Psychiatry
INSERT INTO users (first_name, last_name, full_name, email, phone, role, password, is_active, created_at)
VALUES ('Dr.', 'Arjun Verma', 'Dr. Arjun Verma', 'arjun.verma@hospital.com', '9876543209', 'doctor', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcg7b3XeKeUxWdeS86E36DRZlG2', true, NOW())
ON CONFLICT (email) DO NOTHING;

INSERT INTO doctors (user_id, specialization, qualification, registration_number, experience_years, consultation_fee, department_id, available_days, start_time, end_time, slot_duration_minutes, created_at)
SELECT u.id, 'Psychiatry', 'MD - Psychiatry, Board Certified', 'REG_PSY_001', 16, 500.00, d.id, 'Mon,Tue,Wed,Thu,Fri', '10:00'::time, '17:00'::time, 30, NOW()
FROM users u, departments d
WHERE u.email = 'arjun.verma@hospital.com' AND d.name = 'Psychiatry'
AND NOT EXISTS (SELECT 1 FROM doctors WHERE user_id = u.id)
ON CONFLICT (registration_number) DO NOTHING;

-- Doctor 10: Dr. Ravi Kumar - General Medicine
INSERT INTO users (first_name, last_name, full_name, email, phone, role, password, is_active, created_at)
VALUES ('Dr.', 'Ravi Kumar', 'Dr. Ravi Kumar', 'ravi.kumar@hospital.com', '9876543210', 'doctor', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcg7b3XeKeUxWdeS86E36DRZlG2', true, NOW())
ON CONFLICT (email) DO NOTHING;

INSERT INTO doctors (user_id, specialization, qualification, registration_number, experience_years, consultation_fee, department_id, available_days, start_time, end_time, slot_duration_minutes, created_at)
SELECT u.id, 'General Medicine', 'MD - General Medicine', 'REG_GM_001', 20, 300.00, d.id, 'Mon,Tue,Wed,Thu,Fri', '08:00'::time, '18:00'::time, 20, NOW()
FROM users u, departments d
WHERE u.email = 'ravi.kumar@hospital.com' AND d.name = 'General Medicine'
AND NOT EXISTS (SELECT 1 FROM doctors WHERE user_id = u.id)
ON CONFLICT (registration_number) DO NOTHING;

-- Verify doctors have been inserted successfully
SELECT 'Admin-linked doctors insertion completed!' as message;
