CREATE TABLE patients
(
    patient_id    SERIAL PRIMARY KEY,
    nhs_number    VARCHAR(15) UNIQUE NOT NULL,
    first_name    VARCHAR(50)        NOT NULL,
    last_name     VARCHAR(50)        NOT NULL,
    date_of_birth DATE,
    address       TEXT,
    phone_number  VARCHAR(20),
    email         VARCHAR(100),
    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE medical_records
(
    record_id    SERIAL PRIMARY KEY,
    patient_id   INT REFERENCES patients (patient_id) ON DELETE CASCADE,
    conditions   TEXT,
    medications  TEXT,
    allergies    TEXT,
    notes        TEXT,
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE emergency_calls
(
    call_id               SERIAL PRIMARY KEY,
    patient_id            INT REFERENCES patients (patient_id) ON DELETE SET NULL,
    nhs_number            VARCHAR(15),
    caller_name           VARCHAR(100),
    caller_phone          VARCHAR(20),
    call_time             TIMESTAMP   DEFAULT CURRENT_TIMESTAMP,
    medical_condition     TEXT,
    location              TEXT,
    status                VARCHAR(20) DEFAULT 'Pending', -- e.g., 'Pending', 'Dispatched', 'Completed'
    assigned_ambulance_id INT REFERENCES ambulances (ambulance_id),
    assigned_hospital_id  INT REFERENCES regional_hospitals (hospital_id)
);

CREATE TABLE ambulances
(
    ambulance_id         SERIAL PRIMARY KEY,
    ambulance_number     VARCHAR(20) UNIQUE NOT NULL,
    current_location     GEOGRAPHY(POINT, 4326),          -- Using PostGIS for GPS data
    status               VARCHAR(20) DEFAULT 'Available', -- e.g., 'Available', 'On Call', 'Maintenance'
    regional_hospital_id INT REFERENCES regional_hospitals (hospital_id)
);

CREATE TABLE ambulance_staff
(
    staff_id     SERIAL PRIMARY KEY,
    first_name   VARCHAR(50) NOT NULL,
    last_name    VARCHAR(50) NOT NULL,
    phone_number VARCHAR(20),
    email        VARCHAR(100),
    role         VARCHAR(50), -- e.g., 'Paramedic', 'Driver'
    ambulance_id INT REFERENCES ambulances (ambulance_id),
    is_active    BOOLEAN DEFAULT TRUE
);

CREATE TABLE regional_hospitals
(
    hospital_id         SERIAL PRIMARY KEY,
    name                VARCHAR(100) NOT NULL,
    address             TEXT,
    phone_number        VARCHAR(20),
    email               VARCHAR(100),
    location            GEOGRAPHY(POINT, 4326), -- Using PostGIS for GPS data
    capacity            INT,                    -- Number of beds or patients that can be handled
    available_resources TEXT,                   -- Equipment or specializations
    created_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE call_out_details
(
    detail_id    SERIAL PRIMARY KEY,
    call_id      INT REFERENCES emergency_calls (call_id) ON DELETE CASCADE,
    ambulance_id INT REFERENCES ambulances (ambulance_id),
    action_taken TEXT,
    time_spent INTERVAL,
    notes        TEXT,
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE gps_data
(
    gps_id       SERIAL PRIMARY KEY,
    ambulance_id INT REFERENCES ambulances (ambulance_id) ON DELETE CASCADE,
    timestamp    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    location     GEOGRAPHY(POINT, 4326)
);

CREATE TABLE users
(
    user_id             SERIAL PRIMARY KEY,
    username            VARCHAR(50) UNIQUE NOT NULL,
    password_hash       TEXT               NOT NULL,
    role                VARCHAR(50), -- e.g., 'Operator', 'Hospital Staff', 'Ambulance Staff'
    associated_staff_id INT,         -- References ambulance_staff(staff_id) if applicable
    is_active           BOOLEAN   DEFAULT TRUE,
    created_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);