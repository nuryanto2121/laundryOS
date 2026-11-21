-- 000006_attendance_shifts_and_leaves.up.sql

-- 1. Add coordinates and geofence radius to outlets table
ALTER TABLE outlets ADD COLUMN IF NOT EXISTS latitude NUMERIC(10, 7) DEFAULT -6.113425;
ALTER TABLE outlets ADD COLUMN IF NOT EXISTS longitude NUMERIC(10, 7) DEFAULT 106.914230;
ALTER TABLE outlets ADD COLUMN IF NOT EXISTS geofence_radius_meters INT DEFAULT 25;

-- 2. Master Shifts Table
CREATE TABLE IF NOT EXISTS shifts (
    id                             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code                           VARCHAR(64) NOT NULL UNIQUE,
    shift_name                     VARCHAR(128) NOT NULL,
    start_time                     TIME NOT NULL,
    end_time                       TIME NOT NULL,
    break_start                    TIME,
    break_end                      TIME,
    late_tolerance_minutes         INT NOT NULL DEFAULT 15,
    early_leave_tolerance_minutes  INT NOT NULL DEFAULT 10,
    is_active                      BOOLEAN NOT NULL DEFAULT true,
    created_at                     TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at                     TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Seed FAN Laundry Shift 1 and Shift 2
INSERT INTO shifts (id, code, shift_name, start_time, end_time, break_start, break_end, late_tolerance_minutes, early_leave_tolerance_minutes, is_active)
VALUES
    ('d0000000-0000-0000-0000-000000000001', 'SHIFT_1', 'Shift 1 (Pagi - Sore)', '08:00:00', '17:00:00', '12:00:00', '13:00:00', 15, 10, true),
    ('d0000000-0000-0000-0000-000000000002', 'SHIFT_2', 'Shift 2 (Siang - Malam)', '12:00:00', '21:00:00', '18:00:00', '19:00:00', 15, 10, true)
ON CONFLICT (code) DO NOTHING;

-- 3. Employee Attendances (Log Clock-In, Clock-Out, Selfie Photo, GPS Geofencing, Cashier Drawer)
CREATE TABLE IF NOT EXISTS employee_attendances (
    id                          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    employee_name               VARCHAR(128) NOT NULL,
    employee_phone              VARCHAR(64),
    shift_id                    UUID REFERENCES shifts(id),
    shift_name                  VARCHAR(128) NOT NULL,
    outlet_id                   UUID REFERENCES outlets(id),
    clock_in_at                 TIMESTAMP NOT NULL DEFAULT NOW(),
    clock_in_photo_url          TEXT,
    clock_in_lat                NUMERIC(10, 7),
    clock_in_lng                NUMERIC(10, 7),
    clock_in_distance_meters    NUMERIC(8, 2),
    clock_out_at                TIMESTAMP,
    clock_out_photo_url         TEXT,
    clock_out_lat               NUMERIC(10, 7),
    clock_out_lng               NUMERIC(10, 7),
    clock_out_distance_meters   NUMERIC(8, 2),
    starting_cash               NUMERIC(14, 2) DEFAULT 0,
    ending_cash                 NUMERIC(14, 2) DEFAULT 0,
    status                      VARCHAR(32) NOT NULL DEFAULT 'ON_DUTY', -- ON_DUTY, COMPLETED
    punctuality                 VARCHAR(32) NOT NULL DEFAULT 'ON_TIME', -- ON_TIME, LATE, EARLY_LEAVE, OVERTIME
    late_minutes                INT NOT NULL DEFAULT 0,
    early_leave_minutes         INT NOT NULL DEFAULT 0,
    notes                       TEXT,
    created_at                  TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at                  TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_attendances_clock_in_at ON employee_attendances (clock_in_at DESC);
CREATE INDEX IF NOT EXISTS idx_attendances_employee ON employee_attendances (employee_name, clock_in_at DESC);
CREATE INDEX IF NOT EXISTS idx_attendances_status ON employee_attendances (status);

-- 4. Leave Requests (Pengajuan Libur Karyawan: Max 2x Sebulan)
CREATE TABLE IF NOT EXISTS leave_requests (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    employee_name       VARCHAR(128) NOT NULL,
    employee_phone      VARCHAR(64),
    leave_date          DATE NOT NULL,
    reason              TEXT NOT NULL,
    status              VARCHAR(32) NOT NULL DEFAULT 'PENDING', -- PENDING, APPROVED, REJECTED, CANCELLED
    reviewed_by         VARCHAR(128),
    reviewed_at         TIMESTAMP,
    review_notes        TEXT,
    created_at          TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_leave_requests_date ON leave_requests (leave_date DESC);
CREATE INDEX IF NOT EXISTS idx_leave_requests_emp ON leave_requests (employee_name, leave_date);
CREATE INDEX IF NOT EXISTS idx_leave_requests_status ON leave_requests (status);
