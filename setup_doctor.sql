-- Check if doctor_users table exists
SELECT EXISTS (
    SELECT FROM information_schema.tables 
    WHERE table_name = 'doctor_users'
) AS table_exists;

-- If it exists, check current doctors
SELECT * FROM doctor_users;

-- Insert new doctor if table exists
INSERT INTO doctor_users (email, name, password_hash, specialization, phone, role, is_active)
VALUES (
    'dr.smith@sterling.com',
    'Dr. John Smith',
    '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcg7b3XeKeUxWdeS86E36DRZlG2',
    'General Medicine',
    '+1-555-0001',
    'doctor',
    true
)
ON CONFLICT (email) DO UPDATE SET
    password_hash = EXCLUDED.password_hash,
    name = EXCLUDED.name,
    is_active = true;

-- Verify doctor was created
SELECT id, email, name, specialization, role, is_active FROM doctor_users WHERE email = 'dr.smith@sterling.com';
