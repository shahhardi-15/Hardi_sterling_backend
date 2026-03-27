-- Add new columns to doctors table for enhanced doctor information
ALTER TABLE doctors
ADD COLUMN IF NOT EXISTS experience_years INTEGER DEFAULT 0,
ADD COLUMN IF NOT EXISTS qualification VARCHAR(255),
ADD COLUMN IF NOT EXISTS address TEXT;

-- Update sample doctors with realistic data
UPDATE doctors SET 
  experience_years = 15,
  qualification = 'MD - General Medicine',
  address = '123 Medical Center, City Hospital, Street 1, Main City'
WHERE name = 'Dr. John Smith';

UPDATE doctors SET 
  experience_years = 12,
  qualification = 'MD - Cardiology, Board Certified',
  address = '456 Cardiology Wing, City Hospital, Heart Street, Main City'
WHERE name = 'Dr. Sarah Johnson';

UPDATE doctors SET 
  experience_years = 10,
  qualification = 'MD - Dermatology',
  address = '789 Dermatology Clinic, City Hospital, Skin Lane, Main City'
WHERE name = 'Dr. Michael Brown';

UPDATE doctors SET 
  experience_years = 8,
  qualification = 'MD - Neurology',
  address = '321 Neurology Department, City Hospital, Brain Avenue, Main City'
WHERE name = 'Dr. Emily Davis';

UPDATE doctors SET 
  experience_years = 18,
  qualification = 'MD - Orthopedic Surgery, Board Certified',
  address = '654 Orthopedic Center, City Hospital, Bone Street, Main City'
WHERE name = 'Dr. Robert Wilson';

-- Insert new doctors with additional specializations
INSERT INTO doctors (name, specialization, email, phone, experience_years, qualification, address, is_available) VALUES
  ('Dr. Lisa Anderson', 'Pediatrician', 'lisa.anderson@hospital.com', '+1234567895', 11, 'MD - Pediatrics', '111 Children''s Hospital, Pediatric Wing, Main City', true),
  ('Dr. Jennifer Martinez', 'Gynecologist', 'jennifer.martinez@hospital.com', '+1234567896', 13, 'MD - Obstetrics & Gynecology, Board Certified', '222 Women''s Health Center, Main City', true),
  ('Dr. James Wilson', 'Psychiatrist', 'james.wilson@hospital.com', '+1234567897', 16, 'MD - Psychiatry, Board Certified', '333 Mental Health Institute, Main City', true),
  ('Dr. Patricia Lee', 'Ophthalmologist', 'patricia.lee@hospital.com', '+1234567898', 14, 'MD - Ophthalmology, Board Certified', '444 Eye Care Center, Main City', true),
  ('Dr. David Kim', 'Dentist', 'david.kim@hospital.com', '+1234567899', 9, 'DDS - Dentistry', '555 Dental Clinic, Main City', true),
  ('Dr. Thomas Garcia', 'Emergency Medicine Specialist', 'thomas.garcia@hospital.com', '+1234567900', 12, 'MD - Emergency Medicine, Board Certified', '666 Emergency Department, City Hospital, Main City', true),
  ('Dr. Maria Rodriguez', 'Critical Care Specialist', 'maria.rodriguez@hospital.com', '+1234567901', 15, 'MD - Critical Care Medicine, Board Certified', '777 ICU Unit, City Hospital, Main City', true),
  ('Dr. Christopher Taylor', 'ENT Specialist (Otolaryngologist)', 'christopher.taylor@hospital.com', '+1234567902', 13, 'MD - Otolaryngology, Board Certified', '888 ENT Center, Main City', true),
  ('Dr. Susan White', 'Surgeon (General Surgeon)', 'susan.white@hospital.com', '+1234567903', 18, 'MD - General Surgery, Board Certified', '999 Surgical Suite, City Hospital, Main City', true),
  ('Dr. Andrew Miller', 'Plastic Surgeon', 'andrew.miller@hospital.com', '+1234567904', 14, 'MD - Plastic Surgery, Board Certified', '1010 Cosmetic Surgery Center, Main City', true)
ON CONFLICT DO NOTHING;

-- Create index on specialization for faster queries
CREATE INDEX IF NOT EXISTS idx_doctors_specialization ON doctors(specialization);
