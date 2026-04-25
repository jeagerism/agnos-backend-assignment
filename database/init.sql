-- enable extension to generate UUID automatically (for optional fields)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ==========================================
-- 1. create staffs table
-- ==========================================
CREATE TABLE IF NOT EXISTS staffs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    hospital VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- ==========================================
-- 2. create patients table
-- ==========================================
CREATE TABLE IF NOT EXISTS patients (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    hospital VARCHAR(255) NOT NULL,
    patient_hn VARCHAR(255) NOT NULL,
    
    national_id VARCHAR(20),
    passport_id VARCHAR(50),
    
    first_name_th VARCHAR(100) NOT NULL,
    middle_name_th VARCHAR(100),
    last_name_th VARCHAR(100) NOT NULL,
    
    first_name_en VARCHAR(100),
    middle_name_en VARCHAR(100),
    last_name_en VARCHAR(100),
    
    date_of_birth DATE NOT NULL,
    phone_number VARCHAR(20),
    email VARCHAR(100),
    gender VARCHAR(1) NOT NULL CHECK (gender IN ('M', 'F')),
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- create composite unique index (prevent HN from being the same in the same hospital)
    CONSTRAINT idx_hosp_hn UNIQUE (hospital, patient_hn),
    -- require at least one identity field
    CONSTRAINT chk_patient_identity CHECK (national_id IS NOT NULL OR passport_id IS NOT NULL),
    -- english names must either be both provided or both omitted
    CONSTRAINT chk_en_name_pair CHECK (
        (first_name_en IS NULL AND last_name_en IS NULL) OR
        (first_name_en IS NOT NULL AND last_name_en IS NOT NULL)
    )
);

-- create index for faster search by hospital
CREATE INDEX idx_patients_hospital ON patients(hospital);

-- ==========================================
-- 3. seed data (mock data)
-- ==========================================
-- password is 'password123' that has been hashed (Bcrypt)
INSERT INTO staffs (username, password, hospital) VALUES 
('staff_a', '$2a$10$NcTGHE6lygcBVp8mDPpXvOKBvacjklSIAxgjQrFI6fR7hcC2tSskS', 'Hospital A'),
('staff_b', '$2a$10$NcTGHE6lygcBVp8mDPpXvOKBvacjklSIAxgjQrFI6fR7hcC2tSskS', 'Hospital B')
ON CONFLICT DO NOTHING;

INSERT INTO patients (hospital, patient_hn, national_id, passport_id, first_name_th, last_name_th, date_of_birth, gender) VALUES
('Hospital A', 'HN001', '1111111111111', NULL,         'สมชาย', 'ใจดี', '1990-05-15', 'M'),
('Hospital A', 'HN002', NULL,            'P88888888',  'สมศรี', 'มีสุข', '1992-08-01', 'F')
ON CONFLICT DO NOTHING;

INSERT INTO patients (hospital, patient_hn, passport_id, first_name_th, last_name_th, first_name_en, last_name_en, date_of_birth, gender) VALUES
('Hospital B', 'HN001', 'P12345678', 'จอห์น', 'ดีแลน', 'John', 'Dylan', '1985-11-20', 'M')
ON CONFLICT DO NOTHING;