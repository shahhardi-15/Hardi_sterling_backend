-- Set up doctor credentials for testpat@example.com
-- Password: TestPass@123
-- Hash: $2a$10$zubNRCj59eydSYdoLy6ReOxVUv/gdaHCFDtqeT5uLv8C1W/0CwpTu

-- Check if doctor exists
SELECT id, email FROM doctor_users WHERE email = 'testpat@example.com';

-- Insert or update doctor user
INSERT INTO doctor_users (email, name, password_hash, specialization, phone, role, is_active)
VALUES ('testpat@example.com', 'Test Patient', '$2a$10$zubNRCj59eydSYdoLy6ReOxVUv/gdaHCFDtqeT5uLv8C1W/0CwpTu', 'General Medicine', '555-0001', 'doctor', true)
ON CONFLICT (email) DO UPDATE SET 
  password_hash = '$2a$10$zubNRCj59eydSYdoLy6ReOxVUv/gdaHCFDtqeT5uLv8C1W/0CwpTu',
  is_active = true;

-- Verify
SELECT email, password_hash, is_active FROM doctor_users WHERE email = 'testpat@example.com';
