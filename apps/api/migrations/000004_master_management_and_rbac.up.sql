-- Migration: 000004_master_management_and_rbac (up)
-- Tables and extensions to support custom expense categories, employees management, and RBAC.

-- 1. Master Expense Categories
CREATE TABLE IF NOT EXISTS expense_categories (
    id          UUID        PRIMARY KEY,
    code        VARCHAR     NOT NULL UNIQUE,
    name        VARCHAR     NOT NULL,
    is_active   BOOLEAN     NOT NULL DEFAULT true,
    created_at  TIMESTAMP   NOT NULL DEFAULT NOW()
);

-- Seed default categories
INSERT INTO expense_categories (id, code, name) VALUES
    ('e0000000-0000-0000-0000-000000000001', 'BAHAN_BAKU', 'Deterjen & Bahan Baku'),
    ('e0000000-0000-0000-0000-000000000002', 'PLASTIK', 'Plastik Packing & Kresek'),
    ('e0000000-0000-0000-0000-000000000003', 'LISTRIK_AIR', 'Listrik PLN & Air PAM'),
    ('e0000000-0000-0000-0000-000000000004', 'GAS', 'Gas LPG Mesin Pengering'),
    ('e0000000-0000-0000-0000-000000000005', 'GAJI', 'Gaji & Uang Makan Staf'),
    ('e0000000-0000-0000-0000-000000000006', 'BONUS', 'Bonus & Insentif Kerja'),
    ('e0000000-0000-0000-0000-000000000007', 'SEWA', 'Sewa Tempat / Ruko'),
    ('e0000000-0000-0000-0000-000000000008', 'MAINTENANCE', 'Perawatan & Servis Mesin'),
    ('e0000000-0000-0000-0000-000000000009', 'TRANSPORT', 'Bensin & Operasional Kurir'),
    ('e0000000-0000-0000-0000-000000000010', 'LAIN_LAIN', 'Biaya Lain-lain')
ON CONFLICT (code) DO NOTHING;

-- 2. Master Employees & RBAC
CREATE TABLE IF NOT EXISTS employees (
    id          UUID        PRIMARY KEY,
    username    VARCHAR     NOT NULL UNIQUE,
    full_name   VARCHAR     NOT NULL,
    phone       VARCHAR,
    role        VARCHAR     NOT NULL DEFAULT 'KASIR', -- 'OWNER', 'KASIR', 'KURIR', 'OPERATOR'
    pin         VARCHAR     NOT NULL DEFAULT '1234',
    is_active   BOOLEAN     NOT NULL DEFAULT true,
    created_at  TIMESTAMP   NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP   NOT NULL DEFAULT NOW()
);

-- Seed default staff
INSERT INTO employees (id, username, full_name, phone, role, pin) VALUES
    ('f0000000-0000-0000-0000-000000000001', 'owner', 'Owner Utama (FAN Laundry)', '085285977834', 'OWNER', '1234'),
    ('f0000000-0000-0000-0000-000000000002', 'kasir1', 'Kasir Utama Kalibaru', '081234567890', 'KASIR', '1234')
ON CONFLICT (username) DO NOTHING;
