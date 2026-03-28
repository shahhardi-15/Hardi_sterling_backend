UPDATE doctor_users 
SET password_hash = '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcg7b3XeKeUxWdeS86E36DRZlG2'
WHERE email = 'dr.smith@sterling.com';

SELECT id, email, name, password_hash FROM doctor_users WHERE email = 'dr.smith@sterling.com';
