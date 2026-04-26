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

WITH seed AS (
    SELECT
        gs AS seq,
        CASE WHEN gs <= 50 THEN 'Hospital A' ELSE 'Hospital B' END AS hospital,
        CASE WHEN gs <= 50 THEN gs ELSE gs - 50 END AS hn_seq
    FROM generate_series(1, 100) AS gs
)
INSERT INTO patients (
    hospital,
    patient_hn,
    national_id,
    passport_id,
    first_name_th,
    middle_name_th,
    last_name_th,
    first_name_en,
    middle_name_en,
    last_name_en,
    date_of_birth,
    phone_number,
    email,
    gender
)
SELECT
    s.hospital,
    'HN' || LPAD(s.hn_seq::text, 3, '0') AS patient_hn,
    CASE
        WHEN s.seq % 3 IN (0, 1) THEN LPAD((1000000000000 + s.seq)::text, 13, '0')
        ELSE NULL
    END AS national_id,
    CASE
        WHEN s.seq % 3 IN (0, 2) THEN 'P' || LPAD((88000000 + s.seq)::text, 8, '0')
        ELSE NULL
    END AS passport_id,
    CASE (s.seq % 8)
        WHEN 0 THEN 'สมชาย'
        WHEN 1 THEN 'สมศรี'
        WHEN 2 THEN 'อนันต์'
        WHEN 3 THEN 'สุรีย์'
        WHEN 4 THEN 'ธนา'
        WHEN 5 THEN 'กานต์'
        WHEN 6 THEN 'ปิยะ'
        ELSE 'ศิริพร'
    END AS first_name_th,
    CASE
        WHEN s.seq % 4 = 0 THEN 'กลาง'
        ELSE NULL
    END AS middle_name_th,
    CASE (s.seq % 8)
        WHEN 0 THEN 'ใจดี'
        WHEN 1 THEN 'มีสุข'
        WHEN 2 THEN 'แซ่ลิ้ม'
        WHEN 3 THEN 'ทองคำ'
        WHEN 4 THEN 'ศรีสุข'
        WHEN 5 THEN 'พูนทรัพย์'
        WHEN 6 THEN 'บุญมา'
        ELSE 'พงษ์ดี'
    END AS last_name_th,
    CASE
        WHEN s.seq % 2 = 0 THEN 'Name' || s.seq::text
        ELSE NULL
    END AS first_name_en,
    CASE
        WHEN s.seq % 10 = 0 THEN 'Mid' || s.seq::text
        ELSE NULL
    END AS middle_name_en,
    CASE
        WHEN s.seq % 2 = 0 THEN 'Surname' || s.seq::text
        ELSE NULL
    END AS last_name_en,
    DATE '1980-01-01' + ((s.seq * 37) % 12000) AS date_of_birth,
    CASE
        WHEN s.seq % 5 = 0 THEN NULL
        ELSE '08' || LPAD((10000000 + s.seq)::text, 8, '0')
    END AS phone_number,
    CASE
        WHEN s.seq % 6 = 0 THEN NULL
        ELSE 'patient' || s.seq::text || '@example.com'
    END AS email,
    CASE
        WHEN s.seq % 2 = 0 THEN 'F'
        ELSE 'M'
    END AS gender
FROM seed AS s
ON CONFLICT DO NOTHING;