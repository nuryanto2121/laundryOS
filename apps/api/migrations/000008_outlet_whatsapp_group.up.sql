-- Migration: 000008_outlet_whatsapp_group (up)
-- Adds WhatsApp group routing per-outlet and customer branch preference.

-- 1. Add WhatsApp group fields to outlets table
ALTER TABLE outlets ADD COLUMN IF NOT EXISTS whatsapp_group_id VARCHAR;
ALTER TABLE outlets ADD COLUMN IF NOT EXISTS whatsapp_group_name VARCHAR;

-- 2. Add default outlet preference to customers table
ALTER TABLE customers ADD COLUMN IF NOT EXISTS outlet_id UUID REFERENCES outlets(id);

-- 3. Set default outlet for existing Kalibaru outlet if present
UPDATE outlets
SET whatsapp_group_name = 'Tim Kurir & Operasional Kalibaru'
WHERE id = 'a0000000-0000-0000-0000-000000000001' AND whatsapp_group_name IS NULL;
