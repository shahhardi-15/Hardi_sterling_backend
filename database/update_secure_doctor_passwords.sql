-- Update doctor passwords to secure password: SterlingDoctor@2026
-- Hash: $2a$10$aRd3hZ7Lte3J6STpZQ7Ji.aIMg4AFIauMjwbZrJdqIInOyokJ/IvK

UPDATE doctor_users SET password_hash = '$2a$10$aRd3hZ7Lte3J6STpZQ7Ji.aIMg4AFIauMjwbZrJdqIInOyokJ/IvK' WHERE email = 'dr.smith@sterling.com';
UPDATE doctor_users SET password_hash = '$2a$10$aRd3hZ7Lte3J6STpZQ7Ji.aIMg4AFIauMjwbZrJdqIInOyokJ/IvK' WHERE email = 'dr.johnson@sterling.com';
UPDATE doctor_users SET password_hash = '$2a$10$aRd3hZ7Lte3J6STpZQ7Ji.aIMg4AFIauMjwbZrJdqIInOyokJ/IvK' WHERE email = 'dr.brown@sterling.com';

-- Verify the update
SELECT email, password_hash FROM doctor_users WHERE email LIKE 'dr.%';
