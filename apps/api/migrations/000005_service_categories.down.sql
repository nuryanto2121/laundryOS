-- 000005_service_categories.down.sql
ALTER TABLE services DROP COLUMN IF EXISTS unit;
ALTER TABLE services DROP COLUMN IF EXISTS category_name;
ALTER TABLE services DROP COLUMN IF EXISTS category_id;
DROP TABLE IF EXISTS service_categories;
