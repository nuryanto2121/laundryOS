-- Migration: 000001_init_schema (down)
-- Reverts the initial database schema for LaundryOS.
-- Drops in reverse dependency order.

DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS services;
DROP TABLE IF EXISTS customers;
