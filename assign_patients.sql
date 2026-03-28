-- Assign first 3 patients to doctor 1 (dr.smith@sterling.com)
INSERT INTO doctor_patient_assignment (doctor_id, patient_id, is_active)
SELECT 1, u.id, true 
FROM users u 
LIMIT 3
ON CONFLICT DO NOTHING;

-- Verify assignments
SELECT d.id, d.email, COUNT(dpa.patient_id) as assigned_patients
FROM doctor_users d
LEFT JOIN doctor_patient_assignment dpa ON d.id = dpa.doctor_id AND dpa.is_active = true
WHERE d.email = 'dr.smith@sterling.com'
GROUP BY d.id, d.email;

-- Show assigned patients
SELECT dpa.doctor_id, d.email, u.id as patient_id, u.first_name, u.last_name, u.email as patient_email
FROM doctor_patient_assignment dpa
JOIN doctor_users d ON dpa.doctor_id = d.id
JOIN users u ON dpa.patient_id = u.id
WHERE d.email = 'dr.smith@sterling.com' AND dpa.is_active = true;
