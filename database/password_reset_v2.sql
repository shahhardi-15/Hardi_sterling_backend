-- Updated Password Reset Schema (Simplified without OTP)
-- Run this to replace the previous migration

-- Drop old tables if they exist
DROP TABLE IF EXISTS otp_lockouts CASCADE;
DROP TABLE IF EXISTS password_reset_tokens CASCADE;
DROP TABLE IF EXISTS password_reset_logs CASCADE;

-- Password reset tokens table (simplified)
CREATE TABLE IF NOT EXISTS password_reset_tokens (
  id SERIAL PRIMARY KEY,
  user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  token_hash VARCHAR(255) NOT NULL UNIQUE,
  expires_at TIMESTAMP NOT NULL,
  used_at TIMESTAMP,
  is_used BOOLEAN DEFAULT false,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_password_reset_tokens_user_id ON password_reset_tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_password_reset_tokens_token_hash ON password_reset_tokens(token_hash);
CREATE INDEX IF NOT EXISTS idx_password_reset_tokens_expires_at ON password_reset_tokens(expires_at);
CREATE INDEX IF NOT EXISTS idx_password_reset_tokens_is_used ON password_reset_tokens(is_used);

-- Password reset attempts/logs table for tracking and audit
CREATE TABLE IF NOT EXISTS password_reset_logs (
  id SERIAL PRIMARY KEY,
  user_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
  email VARCHAR(255),
  action VARCHAR(50) NOT NULL, -- 'forgot_password_request', 'reset_success', 'reset_failed'
  ip_address VARCHAR(45),
  user_agent TEXT,
  success BOOLEAN DEFAULT false,
  error_message TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for audit trail and rate limiting
CREATE INDEX IF NOT EXISTS idx_password_reset_logs_user_id ON password_reset_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_password_reset_logs_email ON password_reset_logs(email);
CREATE INDEX IF NOT EXISTS idx_password_reset_logs_email_created ON password_reset_logs(email, created_at);
CREATE INDEX IF NOT EXISTS idx_password_reset_logs_action ON password_reset_logs(action);
CREATE INDEX IF NOT EXISTS idx_password_reset_logs_created_at ON password_reset_logs(created_at DESC);
