-- 000007_employee_acl_and_outlet_assignment.down.sql

DELETE FROM employees WHERE username = 'superadmin';
ALTER TABLE employees DROP COLUMN IF EXISTS permissions;
ALTER TABLE employees DROP COLUMN IF EXISTS email;
ALTER TABLE employees DROP COLUMN IF EXISTS outlet_id;
ALTER TABLE employees DROP COLUMN IF EXISTS password;
