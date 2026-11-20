-- 000005_service_categories.up.sql
-- Create service_categories table
CREATE TABLE IF NOT EXISTS service_categories (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code          VARCHAR(64) NOT NULL UNIQUE,
    name          VARCHAR(128) NOT NULL,
    default_unit  VARCHAR(32) NOT NULL DEFAULT 'Kg',
    display_order INT DEFAULT 0,
    is_active     BOOLEAN DEFAULT true,
    created_at    TIMESTAMP DEFAULT NOW(),
    updated_at    TIMESTAMP DEFAULT NOW()
);

-- Seed core laundry service categories
INSERT INTO service_categories (id, code, name, default_unit, display_order, is_active) VALUES
    ('b0000000-0000-0000-0000-000000000001', 'SELF_SERVICE', 'Self Service (Cuci / Kering)', 'Tabung', 1, true),
    ('b0000000-0000-0000-0000-000000000002', 'CUCI_SETRIKA', 'Cuci Setrika Reguler', 'Kg', 2, true),
    ('b0000000-0000-0000-0000-000000000003', 'SETRIKA', 'Setrika Saja', 'Kg', 3, true),
    ('b0000000-0000-0000-0000-000000000004', 'CUCI_SATUAN', 'Cuci Satuan & Pakaian Khusus', 'Pcs', 4, true),
    ('b0000000-0000-0000-0000-000000000005', 'KARPET', 'Karpet & Permadani', 'Meter', 5, true),
    ('b0000000-0000-0000-0000-000000000006', 'BONEKA', 'Boneka & Bantal', 'Pcs', 6, true),
    ('b0000000-0000-0000-0000-000000000007', 'SEPATU_TAS', 'Sepatu & Tas', 'Pasang', 7, true),
    ('b0000000-0000-0000-0000-000000000008', 'BED_COVER', 'Bed Cover & Selimut', 'Pcs', 8, true)
ON CONFLICT (code) DO NOTHING;

-- Add category_id, category_name, and unit to services table
ALTER TABLE services ADD COLUMN IF NOT EXISTS category_id UUID REFERENCES service_categories(id);
ALTER TABLE services ADD COLUMN IF NOT EXISTS category_name VARCHAR(128) DEFAULT 'Cuci Setrika Reguler';
ALTER TABLE services ADD COLUMN IF NOT EXISTS unit VARCHAR(32) DEFAULT 'Kg';

-- Map existing services to categories
UPDATE services
SET category_id = 'b0000000-0000-0000-0000-000000000001', category_name = 'Self Service (Cuci / Kering)', unit = 'Tabung'
WHERE service_name ILIKE '%self service%' OR service_name ILIKE '%tabung%';

UPDATE services
SET category_id = 'b0000000-0000-0000-0000-000000000002', category_name = 'Cuci Setrika Reguler', unit = 'Kg'
WHERE category_id IS NULL AND (service_name ILIKE '%setrika%' AND service_name ILIKE '%cuci%');

UPDATE services
SET category_id = 'b0000000-0000-0000-0000-000000000003', category_name = 'Setrika Saja', unit = 'Kg'
WHERE category_id IS NULL AND (service_name ILIKE '%setrika%');

UPDATE services
SET category_id = 'b0000000-0000-0000-0000-000000000002', category_name = 'Cuci Setrika Reguler', unit = 'Kg'
WHERE category_id IS NULL;
