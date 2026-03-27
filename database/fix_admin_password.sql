UPDATE admin_users 
SET password_hash = '$2a$10$Xoqs9umkeTkAhlBGSEu/p.ZKpVPNRdme6vE17vTbjY78xPWh4Kvdi'
WHERE email = 'adminsterling@gmail.com';

-- Verify the update
SELECT email, password_hash FROM admin_users WHERE email = 'adminsterling@gmail.com';
