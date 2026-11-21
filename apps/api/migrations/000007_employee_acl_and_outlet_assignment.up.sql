-- 000007_employee_acl_and_outlet_assignment.up.sql
-- Add password, outlet assignment, email, and granular ACL permissions to employees.

ALTER TABLE employees ADD COLUMN IF NOT EXISTS password VARCHAR(255) DEFAULT '123456';
ALTER TABLE employees ADD COLUMN IF NOT EXISTS outlet_id UUID REFERENCES outlets(id);
ALTER TABLE employees ADD COLUMN IF NOT EXISTS email VARCHAR(255);
ALTER TABLE employees ADD COLUMN IF NOT EXISTS permissions JSONB NOT NULL DEFAULT '{}'::jsonb;

-- Sync initial password from PIN if available
UPDATE employees SET password = pin WHERE pin IS NOT NULL AND (password IS NULL OR password = '123456');

-- Ensure default Super Admin account exists
INSERT INTO employees (id, username, full_name, phone, email, role, password, pin, is_active, permissions)
VALUES (
    'f0000000-0000-0000-0000-000000000000',
    'superadmin',
    'Super Administrator (LaundryOS)',
    '081200000000',
    'superadmin@laundryos.id',
    'SUPER_ADMIN',
    'admin123',
    '9999',
    true,
    '{"all": true}'::jsonb
)
ON CONFLICT (username) DO UPDATE
SET role = 'SUPER_ADMIN',
    password = 'admin123',
    is_active = true,
    permissions = '{"all": true}'::jsonb;

-- Update Owner default account
UPDATE employees
SET password = 'owner123',
    email = 'owner@fanlaundry.com',
    outlet_id = 'a0000000-0000-0000-0000-000000000001',
    permissions = '{"all": true}'::jsonb
WHERE username = 'owner';

-- Update Kasir1 default account with standard Kasir ACL permissions
UPDATE employees
SET password = 'kasir123',
    email = 'kasir@fanlaundry.com',
    outlet_id = 'a0000000-0000-0000-0000-000000000001',
    permissions = '{
        "can_create_order": true,
        "can_backdate_order": false,
        "can_edit_order": false,
        "can_cancel_order": false,
        "can_create_expense": true,
        "can_process_wash": true,
        "can_process_dry": true,
        "can_process_iron": true,
        "can_process_pack": true,
        "can_manage_services": false,
        "can_manage_customers": true,
        "can_manage_employees": false,
        "can_view_financials": false,
        "can_view_order_reports": false,
        "can_view_financial_reports": false,
        "can_view_performance_reports": false,
        "can_view_customer_reports": false
    }'::jsonb
WHERE username = 'kasir1';
