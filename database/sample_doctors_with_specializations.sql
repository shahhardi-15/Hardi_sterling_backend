-- Insert sample doctors with all specializations for the appointment booking system
-- This ensures all doctors are available in both patient booking and admin dashboard

INSERT INTO doctors (name, specialization, email, phone, experience_years, qualification, address, is_available, created_at)
VALUES
  ('Dr. Rajesh Kumar', 'Cardiology', 'rajesh.kumar@hospital.com', '9876543201', 15, 'MD - Cardiology, Board Certified', '123 Hearts Hospital, Cardiology Wing, Main City', true, NOW())
ON CONFLICT DO NOTHING;

INSERT INTO doctors (name, specialization, email, phone, experience_years, qualification, address, is_available, created_at)
VALUES
  ('Dr. Priya Sharma', 'Gynecology', 'priya.sharma@hospital.com', '9876543202', 12, 'MD - Obstetrics & Gynecology, Board Certified', '456 Women Health Center, Main City', true, NOW())
ON CONFLICT DO NOTHING;

INSERT INTO doctors (name, specialization, email, phone, experience_years, qualification, address, is_available, created_at)
VALUES
  ('Dr. Anil Patel', 'Orthopedics', 'anil.patel@hospital.com', '9876543203', 18, 'MD - Orthopedic Surgery, Board Certified', '789 Bone & Joint Center, Main City', true, NOW())
ON CONFLICT DO NOTHING;

INSERT INTO doctors (name, specialization, email, phone, experience_years, qualification, address, is_available, created_at)
VALUES
  ('Dr. Neha Desai', 'Pediatrics', 'neha.desai@hospital.com', '9876543204', 10, 'MD - Pediatrics', '321 Children Hospital, Pediatric Wing, Main City', true, NOW())
ON CONFLICT DO NOTHING;

INSERT INTO doctors (name, specialization, email, phone, experience_years, qualification, address, is_available, created_at)
VALUES
  ('Dr. Vikram Singh', 'Neurology', 'vikram.singh@hospital.com', '9876543205', 14, 'MD - Neurology, Board Certified', '654 Brain & Spine Institute, Main City', true, NOW())
ON CONFLICT DO NOTHING;

INSERT INTO doctors (name, specialization, email, phone, experience_years, qualification, address, is_available, created_at)
VALUES
  ('Dr. Ananya Gupta', 'Dermatology', 'ananya.gupta@hospital.com', '9876543206', 11, 'MD - Dermatology', '147 Skin Care Clinic, Main City', true, NOW())
ON CONFLICT DO NOTHING;

INSERT INTO doctors (name, specialization, email, phone, experience_years, qualification, address, is_available, created_at)
VALUES
  ('Dr. Suresh Reddy', 'ENT', 'suresh.reddy@hospital.com', '9876543207', 13, 'MD - Otolaryngology, Board Certified', '258 ENT & Audiology Center, Main City', true, NOW())
ON CONFLICT DO NOTHING;

INSERT INTO doctors (name, specialization, email, phone, experience_years, qualification, address, is_available, created_at)
VALUES
  ('Dr. Divya Menon', 'Ophthalmology', 'divya.menon@hospital.com', '9876543208', 12, 'MD - Ophthalmology, Board Certified', '369 Eye Hospital, Main City', true, NOW())
ON CONFLICT DO NOTHING;

INSERT INTO doctors (name, specialization, email, phone, experience_years, qualification, address, is_available, created_at)
VALUES
  ('Dr. Arjun Verma', 'Psychiatry', 'arjun.verma@hospital.com', '9876543209', 16, 'MD - Psychiatry, Board Certified', '741 Mental Health Institute, Main City', true, NOW())
ON CONFLICT DO NOTHING;

INSERT INTO doctors (name, specialization, email, phone, experience_years, qualification, address, is_available, created_at)
VALUES
  ('Dr. Ravi Kumar', 'General Medicine', 'ravi.kumar@hospital.com', '9876543210', 20, 'MD - General Medicine', '852 General Medicine Clinic, Main City', true, NOW())
ON CONFLICT DO NOTHING;

-- Verify all doctors have been inserted with their specializations
SELECT 'Sample doctors with specializations inserted successfully!' as message;

-- Display all doctors with their specializations
SELECT 
    name,
    specialization,
    qualification,
    email,
    phone,
    experience_years,
    address,
    is_available
FROM doctors
ORDER BY specialization, name;
