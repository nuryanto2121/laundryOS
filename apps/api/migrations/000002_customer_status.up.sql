-- Migration: 000002_customer_status (up)
-- Add lifecycle status column to customers table with safety net default 'ACTIVE'.
-- Note: unique constraint on phone already exists from 000001_init_schema.

ALTER TABLE customers ADD COLUMN status VARCHAR NOT NULL DEFAULT 'ACTIVE';
