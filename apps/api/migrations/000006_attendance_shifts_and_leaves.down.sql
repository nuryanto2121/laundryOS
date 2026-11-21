-- 000006_attendance_shifts_and_leaves.down.sql
DROP TABLE IF EXISTS leave_requests;
DROP TABLE IF EXISTS employee_attendances;
DROP TABLE IF EXISTS shifts;
ALTER TABLE outlets DROP COLUMN IF EXISTS geofence_radius_meters;
ALTER TABLE outlets DROP COLUMN IF EXISTS longitude;
ALTER TABLE outlets DROP COLUMN IF EXISTS latitude;
