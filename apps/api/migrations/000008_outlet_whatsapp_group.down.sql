-- Migration: 000008_outlet_whatsapp_group (down)

ALTER TABLE customers DROP COLUMN IF EXISTS outlet_id;
ALTER TABLE outlets DROP COLUMN IF EXISTS whatsapp_group_name;
ALTER TABLE outlets DROP COLUMN IF EXISTS whatsapp_group_id;
