-- Add new columns to appointments table for patient booking workflow
-- This migration adds fields to track who booked the appointment and capture disapproval reasons

-- Add booked_by_role column
ALTER TABLE appointments
ADD COLUMN IF NOT EXISTS booked_by_role VARCHAR(20)
DEFAULT 'admin';

-- Add disapproval_reason column
ALTER TABLE appointments
ADD COLUMN IF NOT EXISTS disapproval_reason TEXT;

-- Add booked_by column (user ID who booked the appointment)
ALTER TABLE appointments
ADD COLUMN IF NOT EXISTS booked_by INTEGER;

-- Update the status CHECK constraint to include 'pending' status
-- First, we need to drop the old CHECK constraint and create a new one
-- PostgreSQL doesn't allow ALTER CONSTRAINT directly, so we reconstruct it
ALTER TABLE appointments 
DROP CONSTRAINT IF EXISTS appointments_status_check;

ALTER TABLE appointments
ADD CONSTRAINT appointments_status_check 
CHECK (status IN ('scheduled', 'completed', 'cancelled', 'no-show', 'pending'));

-- Verify the migration
SELECT column_name, data_type, column_default 
FROM information_schema.columns 
WHERE table_name = 'appointments' 
ORDER BY ordinal_position;
