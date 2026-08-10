-- Migration: 000002_customer_status (down)
-- Remove status column from customers table.

ALTER TABLE customers DROP COLUMN status;
