-- Add UHID column to users table
ALTER TABLE users
ADD COLUMN IF NOT EXISTS uhid VARCHAR(20) UNIQUE;

-- Create index on UHID for faster lookups
CREATE INDEX IF NOT EXISTS idx_users_uhid ON users(uhid);

-- Update patient_records table to include registration_date if missing
ALTER TABLE patient_records
ADD COLUMN IF NOT EXISTS registration_date TIMESTAMP;

-- Set registration_date to created_at for existing records
UPDATE patient_records 
SET registration_date = created_at 
WHERE registration_date IS NULL;

-- Make registration_date default to current timestamp for new records
ALTER TABLE patient_records 
ALTER COLUMN registration_date SET DEFAULT CURRENT_TIMESTAMP;
