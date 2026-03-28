-- Fix doctor authentication passwords
-- These are bcrypt hashes for password "password123"

UPDATE doctor_users SET password_hash = '$2a$10$UyCEBI2xJtzb2sktH5pbsuZLJhv5RUEgU3/7J4.qZGxasc1uFwtJi' WHERE email = 'dr.smith@sterling.com';
UPDATE doctor_users SET password_hash = '$2a$10$s.M82PpMIMpxJf3vM2Z19eicz2UXW71rNU5DStgxIsdVBzhgpM5Su' WHERE email = 'dr.johnson@sterling.com';
UPDATE doctor_users SET password_hash = '$2a$10$0D9VTV2mvoUDncWTFiDELOSbz2yp7lDuyhF.or37ZVUF3Rkih5aDG' WHERE email = 'dr.brown@sterling.com';

-- Verify the update
SELECT id, email, name, specialization, is_active FROM doctor_users;
