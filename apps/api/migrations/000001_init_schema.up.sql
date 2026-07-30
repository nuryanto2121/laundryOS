-- Migration: 000001_init_schema (up)
-- Creates the initial database schema for LaundryOS.
-- Tables: customers, services, orders

CREATE TABLE customers (
    id              UUID        PRIMARY KEY,
    customer_code   VARCHAR     NOT NULL UNIQUE,
    full_name       VARCHAR     NOT NULL,
    phone           VARCHAR     NOT NULL UNIQUE,
    email           VARCHAR,
    address         TEXT,
    latitude        DECIMAL(10,7),
    longitude       DECIMAL(10,7),
    created_at      TIMESTAMP,
    updated_at      TIMESTAMP,
    deleted_at      TIMESTAMP
);

CREATE TABLE services (
    id              UUID        PRIMARY KEY,
    service_code    VARCHAR     NOT NULL UNIQUE,
    service_name    VARCHAR     NOT NULL,
    duration_day    INTEGER,
    price_per_kg    DECIMAL(12,2),
    active          BOOLEAN
);

CREATE TABLE orders (
    id               UUID        PRIMARY KEY,
    order_number     VARCHAR     NOT NULL UNIQUE,
    customer_id      UUID        NOT NULL REFERENCES customers (id),
    service_id       UUID        NOT NULL REFERENCES services (id),
    status           VARCHAR,
    total_weight     DECIMAL(10,2),
    total_price      DECIMAL(12,2),
    pickup_required  BOOLEAN,
    notes            TEXT,
    created_at       TIMESTAMP,
    updated_at       TIMESTAMP
);
