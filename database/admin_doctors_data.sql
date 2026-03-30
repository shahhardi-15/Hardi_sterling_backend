-- Create doctor users linked to the doctors table for admin management
-- This ensures doctors are visible in the admin "Manage Doctors" page

-- First, clean up any existing doctors with registration numbers to avoid conflicts
DELETE FROM doctors WHERE registration_number IN ('REG001', 'REG002', 'REG003', 'REG004', 'REG005', 'REG006', 'REG007', 'REG008', 'REG009', 'REG010');

-- Delete users that we're about to create
DELETE FROM users WHERE email IN (
  'rajesh.kumar@hospital.com',
  'priya.sharma@hospital.com',
  'anil.patel@hospital.com',
  'neha.desai@hospital.com',
  'vikram.singh@hospital.com',
  'ananya.gupta@hospital.com',
  'suresh.reddy@hospital.com',
  'divya.menon@hospital.com',
  'arjun.verma@hospital.com',
  'ravi.kumar@hospital.com'
);

-- Get department IDs for insertion (to be used in doctor records)
-- Doctor 1: Dr. Rajesh Kumar - Cardiology
INSERT INTO users (full_name, email, phone, role, password_hash, is_active, created_at)
VALUES ('Dr. Rajesh Kumar', 'rajesh.kumar@hospital.com', '9876543201', 'doctor', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcg7b3XeKeUxWdeS86E36DRZlG2', true, NOW())
ON CONFLICT (email) DO NOTHING;

INSERT INTO doctors (user_id, specialization, qualification, registration_number, experience_years, consultation_fee, department_id, available_days, start_time, end_time, slot_duration_minutes, created_at)
SELECT u.id, 'Cardiology', 'MD - Cardiology, Board Certified', 'REG001', 15, 500.00, d.id, 'Mon,Tue,Wed,Thu,Fri', '09:00'::time, '17:00'::time, 30, NOW()
FROM users u, departments d
WHERE u.email = 'rajesh.kumar@hospital.com' AND d.name = 'Cardiology'
AND NOT EXISTS (SELECT 1 FROM doctors WHERE user_id = u.id)
ON CONFLICT (registration_number) DO NOTHING;

-- Doctor 2: Dr. Priya Sharma - Gynecology
INSERT INTO users (full_name, email, phone, role, password_hash, is_active, created_at)
VALUES ('Dr. Priya Sharma', 'priya.sharma@hospital.com', '9876543202', 'doctor', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcg7b3XeKeUxWdeS86E36DRZlG2', true, NOW())
ON CONFLICT (email) DO NOTHING;

INSERT INTO doctors (user_id, specialization, qualification, registration_number, experience_years, consultation_fee, department_id, available_days, start_time, end_time, slot_duration_minutes, created_at)
SELECT u.id, 'Gynecology', 'MD - Obstetrics & Gynecology, Board Certified', 'REG002', 12, 450.00, d.id, 'Mon,Tue,Wed,Thu,Fri', '10:00'::time, '18:00'::time, 30, NOW()
FROM users u, departments d
WHERE u.email = 'priya.sharma@hospital.com' AND d.name = 'Gynecology'
AND NOT EXISTS (SELECT 1 FROM doctors WHERE user_id = u.id)
ON CONFLICT (registration_number) DO NOTHING;

-- Doctor 3: Dr. Anil Patel - Orthopedics
INSERT INTO users (full_name, email, phone, role, password_hash, is_active, created_at)
VALUES ('Dr. Anil Patel', 'anil.patel@hospital.com', '9876543203', 'doctor', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcg7b3XeKeUxWdeS86E36DRZlG2', true, NOW())
ON CONFLICT (email) DO NOTHING;

INSERT INTO doctors (user_id, specialization, qualification, registration_number, experience_years, consultation_fee, department_id, available_days, start_time, end_time, slot_duration_minutes, created_at)
SELECT u.id, 'Orthopedics', 'MD - Orthopedic Surgery, Board Certified', 'REG003', 18, 550.00, d.id, 'Mon,Tue,Wed,Thu,Fri', '09:00'::time, '16:00'::time, 30, NOW()
FROM users u, departments d
WHERE u.email = 'anil.patel@hospital.com' AND d.name = 'Orthopedics'
AND NOT EXISTS (SELECT 1 FROM doctors WHERE user_id = u.id)
ON CONFLICT (registration_number) DO NOTHING;

-- Doctor 4: Dr. Neha Desai - Pediatrics
INSERT INTO users (full_name, email, phone, role, password_hash, is_active, created_at)
VALUES ('Dr. Neha Desai', 'neha.desai@hospital.com', '9876543204', 'doctor', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcg7b3XeKeUxWdeS86E36DRZlG2', true, NOW())
ON CONFLICT (email) DO NOTHING;

INSERT INTO doctors (user_id, specialization, qualification, registration_number, experience_years, consultation_fee, department_id, available_days, start_time, end_time, slot_duration_minutes, created_at)
SELECT u.id, 'Pediatrics', 'MD - Pediatrics', 'REG004', 10, 400.00, d.id, 'Mon,Tue,Wed,Thu,Fri', '08:00'::time, '16:00'::time, 20, NOW()
FROM users u, departments d
WHERE u.email = 'neha.desai@hospital.com' AND d.name = 'Pediatrics'
AND NOT EXISTS (SELECT 1 FROM doctors WHERE user_id = u.id)
ON CONFLICT (registration_number) DO NOTHING;

-- Doctor 5: Dr. Vikram Singh - Neurology
INSERT INTO users (full_name, email, phone, role, password_hash, is_active, created_at)
VALUES ('Dr. Vikram Singh', 'vikram.singh@hospital.com', '9876543205', 'doctor', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcg7b3XeKeUxWdeS86E36DRZlG2', true, NOW())
ON CONFLICT (email) DO NOTHING;

INSERT INTO doctors (user_id, specialization, qualification, registration_number, experience_years, consultation_fee, department_id, available_days, start_time, end_time, slot_duration_minutes, created_at)
SELECT u.id, 'Neurology', 'MD - Neurology, Board Certified', 'REG005', 14, 600.00, d.id, 'Mon,Tue,Wed,Thu,Fri', '10:00'::time, '17:00'::time, 30, NOW()
FROM users u, departments d
WHERE u.email = 'vikram.singh@hospital.com' AND d.name = 'Neurology'
AND NOT EXISTS (SELECT 1 FROM doctors WHERE user_id = u.id)
ON CONFLICT (registration_number) DO NOTHING;

-- Doctor 6: Dr. Ananya Gupta - Dermatology
INSERT INTO users (full_name, email, phone, role, password_hash, is_active, created_at)
VALUES ('Dr. Ananya Gupta', 'ananya.gupta@hospital.com', '9876543206', 'doctor', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcg7b3XeKeUxWdeS86E36DRZlG2', true, NOW())
ON CONFLICT (email) DO NOTHING;

INSERT INTO doctors (user_id, specialization, qualification, registration_number, experience_years, consultation_fee, department_id, available_days, start_time, end_time, slot_duration_minutes, created_at)
SELECT u.id, 'Dermatology', 'MD - Dermatology', 'REG006', 11, 350.00, d.id, 'Mon,Tue,Wed,Thu,Fri', '09:00'::time, '17:00'::time, 20, NOW()
FROM users u, departments d
WHERE u.email = 'ananya.gupta@hospital.com' AND d.name = 'Dermatology'
AND NOT EXISTS (SELECT 1 FROM doctors WHERE user_id = u.id)
ON CONFLICT (registration_number) DO NOTHING;

-- Doctor 7: Dr. Suresh Reddy - ENT
INSERT INTO users (full_name, email, phone, role, password_hash, is_active, created_at)
VALUES ('Dr. Suresh Reddy', 'suresh.reddy@hospital.com', '9876543207', 'doctor', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcg7b3XeKeUxWdeS86E36DRZlG2', true, NOW())
ON CONFLICT (email) DO NOTHING;

INSERT INTO doctors (user_id, specialization, qualification, registration_number, experience_years, consultation_fee, department_id, available_days, start_time, end_time, slot_duration_minutes, created_at)
SELECT u.id, 'ENT', 'MD - Otolaryngology, Board Certified', 'REG007', 13, 400.00, d.id, 'Mon,Tue,Wed,Thu,Fri', '09:00'::time, '17:00'::time, 20, NOW()
FROM users u, departments d
WHERE u.email = 'suresh.reddy@hospital.com' AND d.name = 'ENT'
AND NOT EXISTS (SELECT 1 FROM doctors WHERE user_id = u.id)
ON CONFLICT (registration_number) DO NOTHING;

-- Doctor 8: Dr. Divya Menon - Ophthalmology
INSERT INTO users (full_name, email, phone, role, password_hash, is_active, created_at)
VALUES ('Dr. Divya Menon', 'divya.menon@hospital.com', '9876543208', 'doctor', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcg7b3XeKeUxWdeS86E36DRZlG2', true, NOW())
ON CONFLICT (email) DO NOTHING;

INSERT INTO doctors (user_id, specialization, qualification, registration_number, experience_years, consultation_fee, department_id, available_days, start_time, end_time, slot_duration_minutes, created_at)
SELECT u.id, 'Ophthalmology', 'MD - Ophthalmology, Board Certified', 'REG008', 12, 450.00, d.id, 'Mon,Tue,Wed,Thu,Fri', '10:00'::time, '18:00'::time, 20, NOW()
FROM users u, departments d
WHERE u.email = 'divya.menon@hospital.com' AND d.name = 'Ophthalmology'
AND NOT EXISTS (SELECT 1 FROM doctors WHERE user_id = u.id)
ON CONFLICT (registration_number) DO NOTHING;

-- Doctor 9: Dr. Arjun Verma - Psychiatry
INSERT INTO users (full_name, email, phone, role, password_hash, is_active, created_at)
VALUES ('Dr. Arjun Verma', 'arjun.verma@hospital.com', '9876543209', 'doctor', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcg7b3XeKeUxWdeS86E36DRZlG2', true, NOW())
ON CONFLICT (email) DO NOTHING;

INSERT INTO doctors (user_id, specialization, qualification, registration_number, experience_years, consultation_fee, department_id, available_days, start_time, end_time, slot_duration_minutes, created_at)
SELECT u.id, 'Psychiatry', 'MD - Psychiatry, Board Certified', 'REG009', 16, 500.00, d.id, 'Mon,Tue,Wed,Thu,Fri', '10:00'::time, '17:00'::time, 30, NOW()
FROM users u, departments d
WHERE u.email = 'arjun.verma@hospital.com' AND d.name = 'Psychiatry'
AND NOT EXISTS (SELECT 1 FROM doctors WHERE user_id = u.id)
ON CONFLICT (registration_number) DO NOTHING;

-- Doctor 10: Dr. Ravi Kumar - General Medicine
INSERT INTO users (full_name, email, phone, role, password_hash, is_active, created_at)
VALUES ('Dr. Ravi Kumar', 'ravi.kumar@hospital.com', '9876543210', 'doctor', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcg7b3XeKeUxWdeS86E36DRZlG2', true, NOW())
ON CONFLICT (email) DO NOTHING;

INSERT INTO doctors (user_id, specialization, qualification, registration_number, experience_years, consultation_fee, department_id, available_days, start_time, end_time, slot_duration_minutes, created_at)
SELECT u.id, 'General Medicine', 'MD - General Medicine', 'REG010', 20, 300.00, d.id, 'Mon,Tue,Wed,Thu,Fri', '08:00'::time, '18:00'::time, 20, NOW()
FROM users u, departments d
WHERE u.email = 'ravi.kumar@hospital.com' AND d.name = 'General Medicine'
AND NOT EXISTS (SELECT 1 FROM doctors WHERE user_id = u.id)
ON CONFLICT (registration_number) DO NOTHING;

-- Verify doctors have been inserted successfully
SELECT 'Admin-linked doctors inserted successfully!' as message;

-- Display all doctors with proper admin schema
SELECT 
    d.id,
    u.full_name,
    d.specialization,
    d.qualification,
    d.registration_number,
    d.experience_years,
    d.consultation_fee,
    dept.name as department_name,
    u.email,
    u.phone,
    u.is_active
FROM doctors d
JOIN users u ON d.user_id = u.id
LEFT JOIN departments dept ON d.department_id = dept.id
ORDER BY u.full_name;
