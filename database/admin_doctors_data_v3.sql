-- Smart migration to link doctors with users (one doctor per specialization)
-- This version avoids UNIQUE constraint violations

-- Step 1: Ensure full_name column exists and is populated
ALTER TABLE users ADD COLUMN IF NOT EXISTS full_name VARCHAR(255);
UPDATE users SET full_name = CONCAT(first_name, ' ', last_name) WHERE full_name IS NULL AND first_name IS NOT NULL;

-- Step 2: Ensure doctors table has necessary columns and user_id foreign key
ALTER TABLE doctors 
ADD COLUMN IF NOT EXISTS user_id INTEGER UNIQUE REFERENCES users(id) ON DELETE CASCADE,
ADD COLUMN IF NOT EXISTS registration_number VARCHAR(100) UNIQUE,
ADD COLUMN IF NOT EXISTS consultation_fee DECIMAL(10, 2) DEFAULT 0.00,
ADD COLUMN IF NOT EXISTS department_id UUID REFERENCES departments(id) ON DELETE SET NULL,
ADD COLUMN IF NOT EXISTS available_days TEXT,
ADD COLUMN IF NOT EXISTS start_time TIME,
ADD COLUMN IF NOT EXISTS end_time TIME,
ADD COLUMN IF NOT EXISTS slot_duration_minutes INTEGER DEFAULT 15;

-- Step 3: Create doctor users (one per specialization)
INSERT INTO users (first_name, last_name, full_name, email, phone, role, password, is_active, created_at)
VALUES 
('Dr.', 'Rajesh Kumar', 'Dr. Rajesh Kumar', 'rajesh.kumar@hospital.com', '9876543201', 'doctor', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcg7b3XeKeUxWdeS86E36DRZlG2', true, NOW()),
('Dr.', 'Priya Sharma', 'Dr. Priya Sharma', 'priya.sharma@hospital.com', '9876543202', 'doctor', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcg7b3XeKeUxWdeS86E36DRZlG2', true, NOW()),
('Dr.', 'Anil Patel', 'Dr. Anil Patel', 'anil.patel@hospital.com', '9876543203', 'doctor', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcg7b3XeKeUxWdeS86E36DRZlG2', true, NOW()),
('Dr.', 'Neha Desai', 'Dr. Neha Desai', 'neha.desai@hospital.com', '9876543204', 'doctor', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcg7b3XeKeUxWdeS86E36DRZlG2', true, NOW()),
('Dr.', 'Vikram Singh', 'Dr. Vikram Singh', 'vikram.singh@hospital.com', '9876543205', 'doctor', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcg7b3XeKeUxWdeS86E36DRZlG2', true, NOW()),
('Dr.', 'Ananya Gupta', 'Dr. Ananya Gupta', 'ananya.gupta@hospital.com', '9876543206', 'doctor', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcg7b3XeKeUxWdeS86E36DRZlG2', true, NOW()),
('Dr.', 'Suresh Reddy', 'Dr. Suresh Reddy', 'suresh.reddy@hospital.com', '9876543207', 'doctor', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcg7b3XeKeUxWdeS86E36DRZlG2', true, NOW()),
('Dr.', 'Divya Menon', 'Dr. Divya Menon', 'divya.menon@hospital.com', '9876543208', 'doctor', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcg7b3XeKeUxWdeS86E36DRZlG2', true, NOW()),
('Dr.', 'Arjun Verma', 'Dr. Arjun Verma', 'arjun.verma@hospital.com', '9876543209', 'doctor', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcg7b3XeKeUxWdeS86E36DRZlG2', true, NOW()),
('Dr.', 'Ravi Kumar', 'Dr. Ravi Kumar', 'ravi.kumar@hospital.com', '9876543210', 'doctor', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcg7b3XeKeUxWdeS86E36DRZlG2', true, NOW())
ON CONFLICT (email) DO NOTHING;

-- Step 4: Link FIRST doctor of each specialization with corresponding doctor user (one-to-one)
UPDATE doctors d
SET user_id = (SELECT id FROM users WHERE email = CASE 
    WHEN d.specialization = 'Cardiology' THEN 'rajesh.kumar@hospital.com'
    WHEN d.specialization = 'Gynecology' THEN 'priya.sharma@hospital.com'
    WHEN d.specialization = 'Orthopedics' THEN 'anil.patel@hospital.com'
    WHEN d.specialization = 'Pediatrics' THEN 'neha.desai@hospital.com'
    WHEN d.specialization = 'Neurology' THEN 'vikram.singh@hospital.com'
    WHEN d.specialization = 'Dermatology' THEN 'ananya.gupta@hospital.com'
    WHEN d.specialization = 'ENT' THEN 'suresh.reddy@hospital.com'
    WHEN d.specialization = 'Ophthalmology' THEN 'divya.menon@hospital.com'
    WHEN d.specialization = 'Psychiatry' THEN 'arjun.verma@hospital.com'
    WHEN d.specialization = 'General Medicine' THEN 'ravi.kumar@hospital.com'
END)
WHERE user_id IS NULL 
  AND specialization IN ('Cardiology', 'Gynecology', 'Orthopedics', 'Pediatrics', 'Neurology', 'Dermatology', 'ENT', 'Ophthalmology', 'Psychiatry', 'General Medicine')
  AND id IN (
    SELECT DISTINCT ON (specialization) id 
    FROM doctors 
    WHERE specialization IN ('Cardiology', 'Gynecology', 'Orthopedics', 'Pediatrics', 'Neurology', 'Dermatology', 'ENT', 'Ophthalmology', 'Psychiatry', 'General Medicine')
    ORDER BY specialization, id
  );

-- Step 5: Update registration numbers for doctors without them
UPDATE doctors
SET registration_number = CONCAT('REG_', UPPER(SUBSTRING(specialization, 1, 3)), '_', id)
WHERE registration_number IS NULL;

-- Step 6: Update default consultation fees
UPDATE doctors
SET consultation_fee = 450.00
WHERE consultation_fee = 0 OR consultation_fee IS NULL;

-- Step 7: Set availability defaults
UPDATE doctors
SET available_days = 'Mon,Tue,Wed,Thu,Fri',
    start_time = '09:00'::time,
    end_time = '17:00'::time,
    slot_duration_minutes = 30
WHERE available_days IS NULL;

-- Step 8: Update department_id for all doctors that have a specialization matching a department
UPDATE doctors d
SET department_id = (SELECT id FROM departments WHERE name = d.specialization)
WHERE department_id IS NULL AND specialization IN (SELECT name FROM departments);

-- Confirmation
SELECT 'Doctor-User relationship setup completed!' as status;
SELECT COUNT(*) as doctors_with_user_id FROM doctors WHERE user_id IS NOT NULL;
SELECT COUNT(*) as total_doctors FROM doctors;
