-- Migration: 000003_pos_and_expenses (down)

DROP TABLE IF EXISTS order_status_history;
DROP TABLE IF EXISTS order_photos;
DROP TABLE IF EXISTS expenses;
DROP TABLE IF EXISTS order_items;

ALTER TABLE orders DROP COLUMN IF EXISTS canceled_at;
ALTER TABLE orders DROP COLUMN IF EXISTS canceled_reason;
ALTER TABLE orders DROP COLUMN IF EXISTS completed_at;
ALTER TABLE orders DROP COLUMN IF EXISTS estimated_completion_at;
ALTER TABLE orders DROP COLUMN IF EXISTS paid_at;
ALTER TABLE orders DROP COLUMN IF EXISTS payment_method_name;
ALTER TABLE orders DROP COLUMN IF EXISTS payment_method_id;
ALTER TABLE orders DROP COLUMN IF EXISTS payment_status;
ALTER TABLE orders DROP COLUMN IF EXISTS paid_amount;
ALTER TABLE orders DROP COLUMN IF EXISTS discount_amount;
ALTER TABLE orders DROP COLUMN IF EXISTS discount_id;
ALTER TABLE orders DROP COLUMN IF EXISTS perfume_name;
ALTER TABLE orders DROP COLUMN IF EXISTS perfume_id;
ALTER TABLE orders DROP COLUMN IF EXISTS rack_number;
ALTER TABLE orders DROP COLUMN IF EXISTS outlet_id;

DROP TABLE IF EXISTS payment_methods;
DROP TABLE IF EXISTS discounts;
DROP TABLE IF EXISTS perfumes;
DROP TABLE IF EXISTS outlets;
