-- Reset postgres user password
ALTER USER postgres WITH PASSWORD 'admin';

-- Create database if it doesn't exist
CREATE DATABASE sterling;

-- Connect to sterling database (this can't be done in batch, so we'll do it separately)
