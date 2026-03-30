-- Create departments table if it doesn't exist
CREATE TABLE IF NOT EXISTS departments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(150) NOT NULL UNIQUE,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Create doctors table if it doesn't exist
CREATE TABLE IF NOT EXISTS doctors (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id INTEGER NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    specialization VARCHAR(255) NOT NULL,
    qualification VARCHAR(255),
    registration_number VARCHAR(100) UNIQUE NOT NULL,
    experience_years INTEGER DEFAULT 0,
    consultation_fee DECIMAL(10, 2) DEFAULT 0.00,
    department_id UUID REFERENCES departments(id) ON DELETE SET NULL,
    available_days TEXT, -- Stored as comma-separated string: "Mon,Tue,Wed,Thu,Fri"
    start_time TIME,
    end_time TIME,
    slot_duration_minutes INTEGER DEFAULT 15,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_doctors_user_id ON doctors(user_id);
CREATE INDEX IF NOT EXISTS idx_doctors_registration_number ON doctors(registration_number);
CREATE INDEX IF NOT EXISTS idx_doctors_department_id ON doctors(department_id);
CREATE INDEX IF NOT EXISTS idx_doctors_specialization ON doctors(specialization);

-- Create indexes on departments
CREATE INDEX IF NOT EXISTS idx_departments_name ON departments(name);
CREATE INDEX IF NOT EXISTS idx_departments_is_active ON departments(is_active);

-- Insert default departments if the table is empty
INSERT INTO departments (name) VALUES
    ('General Medicine'),
    ('Cardiology'),
    ('Orthopedics'),
    ('Pediatrics'),
    ('Gynecology'),
    ('Neurology'),
    ('Dermatology'),
    ('ENT'),
    ('Ophthalmology'),
    ('Psychiatry'),
    ('Oncology'),
    ('Radiology'),
    ('Pathology'),
    ('Emergency')
ON CONFLICT (name) DO NOTHING;
