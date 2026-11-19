-- Migration: 000003_pos_and_expenses (up)
-- Tables and extensions to support Web POS, multi-item orders, expenses, photo documentation, and master data.

-- 1. Outlets
CREATE TABLE IF NOT EXISTS outlets (
    id              UUID        PRIMARY KEY,
    name            VARCHAR     NOT NULL,
    code            VARCHAR     NOT NULL UNIQUE,
    address         TEXT,
    phone           VARCHAR,
    is_active       BOOLEAN     NOT NULL DEFAULT true,
    created_at      TIMESTAMP   NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP   NOT NULL DEFAULT NOW()
);

-- 2. Perfumes
CREATE TABLE IF NOT EXISTS perfumes (
    id              UUID        PRIMARY KEY,
    name            VARCHAR     NOT NULL,
    description     TEXT,
    is_active       BOOLEAN     NOT NULL DEFAULT true,
    created_at      TIMESTAMP   NOT NULL DEFAULT NOW()
);

-- 3. Discounts
CREATE TABLE IF NOT EXISTS discounts (
    id               UUID           PRIMARY KEY,
    name             VARCHAR        NOT NULL,
    discount_type    VARCHAR        NOT NULL, -- 'PERCENTAGE', 'NOMINAL'
    value            DECIMAL(10,2)  NOT NULL,
    min_order_amount DECIMAL(12,2)  NOT NULL DEFAULT 0.00,
    is_active        BOOLEAN        NOT NULL DEFAULT true,
    created_at       TIMESTAMP      NOT NULL DEFAULT NOW()
);

-- 4. Payment Methods
CREATE TABLE IF NOT EXISTS payment_methods (
    id              UUID        PRIMARY KEY,
    code            VARCHAR     NOT NULL UNIQUE,
    name            VARCHAR     NOT NULL,
    category        VARCHAR     NOT NULL, -- 'TUNAI', 'TRANSFER', 'QRIS'
    account_number  VARCHAR,
    account_name    VARCHAR,
    is_active       BOOLEAN     NOT NULL DEFAULT true,
    created_at      TIMESTAMP   NOT NULL DEFAULT NOW()
);

-- 5. Update Orders Table (Alter & Add POS columns)
ALTER TABLE orders ALTER COLUMN service_id DROP NOT NULL;
ALTER TABLE orders ADD COLUMN IF NOT EXISTS outlet_id UUID REFERENCES outlets (id);
ALTER TABLE orders ADD COLUMN IF NOT EXISTS rack_number VARCHAR;
ALTER TABLE orders ADD COLUMN IF NOT EXISTS perfume_id UUID REFERENCES perfumes (id);
ALTER TABLE orders ADD COLUMN IF NOT EXISTS perfume_name VARCHAR;
ALTER TABLE orders ADD COLUMN IF NOT EXISTS discount_id UUID REFERENCES discounts (id);
ALTER TABLE orders ADD COLUMN IF NOT EXISTS discount_amount DECIMAL(12,2) NOT NULL DEFAULT 0.00;
ALTER TABLE orders ADD COLUMN IF NOT EXISTS paid_amount DECIMAL(12,2) NOT NULL DEFAULT 0.00;
ALTER TABLE orders ADD COLUMN IF NOT EXISTS payment_status VARCHAR NOT NULL DEFAULT 'UNPAID';
ALTER TABLE orders ADD COLUMN IF NOT EXISTS payment_method_id UUID REFERENCES payment_methods (id);
ALTER TABLE orders ADD COLUMN IF NOT EXISTS payment_method_name VARCHAR;
ALTER TABLE orders ADD COLUMN IF NOT EXISTS paid_at TIMESTAMP;
ALTER TABLE orders ADD COLUMN IF NOT EXISTS estimated_completion_at TIMESTAMP;
ALTER TABLE orders ADD COLUMN IF NOT EXISTS completed_at TIMESTAMP;
ALTER TABLE orders ADD COLUMN IF NOT EXISTS canceled_reason TEXT;
ALTER TABLE orders ADD COLUMN IF NOT EXISTS canceled_at TIMESTAMP;

-- 6. Order Items (Support multi-item / multi-layanan per order)
CREATE TABLE IF NOT EXISTS order_items (
    id              UUID           PRIMARY KEY,
    order_id        UUID           NOT NULL REFERENCES orders (id) ON DELETE CASCADE,
    service_id      UUID           REFERENCES services (id),
    service_name    VARCHAR        NOT NULL,
    category        VARCHAR        NOT NULL DEFAULT 'KILOAN',
    quantity        DECIMAL(10,2)  NOT NULL DEFAULT 1.00,
    unit            VARCHAR        NOT NULL DEFAULT 'KG',
    price_per_unit  DECIMAL(12,2)  NOT NULL DEFAULT 0.00,
    subtotal        DECIMAL(12,2)  NOT NULL DEFAULT 0.00,
    notes           TEXT,
    created_at      TIMESTAMP      NOT NULL DEFAULT NOW()
);

-- 7. Expenses (Kas Kecil / Pengeluaran Operasional)
CREATE TABLE IF NOT EXISTS expenses (
    id              UUID           PRIMARY KEY,
    outlet_id       UUID           REFERENCES outlets (id),
    expense_date    DATE           NOT NULL DEFAULT CURRENT_DATE,
    category        VARCHAR        NOT NULL, -- 'GAJI', 'BONUS', 'BAHAN_BAKU', 'LISTRIK_AIR', 'SEWA', 'GAS', 'MAINTENANCE', 'LAIN_LAIN'
    amount          DECIMAL(12,2)  NOT NULL,
    description     TEXT           NOT NULL,
    created_by      VARCHAR        NOT NULL DEFAULT 'Kasir',
    created_at      TIMESTAMP      NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP      NOT NULL DEFAULT NOW()
);

-- 8. Order Photos (Dokumentasi Masuk & Keluar)
CREATE TABLE IF NOT EXISTS order_photos (
    id              UUID        PRIMARY KEY,
    order_id        UUID        NOT NULL REFERENCES orders (id) ON DELETE CASCADE,
    photo_type      VARCHAR     NOT NULL, -- 'MASUK', 'KELUAR'
    photo_url       TEXT        NOT NULL,
    notes           TEXT,
    created_at      TIMESTAMP   NOT NULL DEFAULT NOW()
);

-- 9. Order Status History (Audit Trail)
CREATE TABLE IF NOT EXISTS order_status_history (
    id              UUID        PRIMARY KEY,
    order_id        UUID        NOT NULL REFERENCES orders (id) ON DELETE CASCADE,
    previous_status VARCHAR,
    new_status      VARCHAR     NOT NULL,
    changed_by      VARCHAR     NOT NULL DEFAULT 'SYSTEM',
    notes           TEXT,
    created_at      TIMESTAMP   NOT NULL DEFAULT NOW()
);

-- 10. Default Seed Data
INSERT INTO outlets (id, name, code, address, phone)
VALUES ('a0000000-0000-0000-0000-000000000001', 'FAN Laundry Kalibaru', 'OUTLET-01', 'Jl. Kalibaru Barat I Rt 007 Rw 008 kel. Kalibaru kec. Cilincing Jakarta Utara 14110', '085285977834')
ON CONFLICT (code) DO NOTHING;

INSERT INTO perfumes (id, name, description)
VALUES 
    ('b0000000-0000-0000-0000-000000000001', 'Akasia', 'Aroma segar bunga akasia menenangkan'),
    ('b0000000-0000-0000-0000-000000000002', 'Sakura', 'Aroma manis lembut bunga sakura khas Jepang'),
    ('b0000000-0000-0000-0000-000000000003', 'Ocean Fresh', 'Aroma segar semilir angin laut tahan lama'),
    ('b0000000-0000-0000-0000-000000000004', 'Snappy', 'Aroma harum segar dan elegan'),
    ('b0000000-0000-0000-0000-000000000005', 'Downy Black', 'Aroma mewah elegan parfum modern')
ON CONFLICT DO NOTHING;

INSERT INTO payment_methods (id, code, name, category, account_number, account_name)
VALUES
    ('c0000000-0000-0000-0000-000000000001', 'TUNAI', 'Tunai / Cash', 'TUNAI', NULL, NULL),
    ('c0000000-0000-0000-0000-000000000002', 'BCA', 'Transfer BCA', 'TRANSFER', '4140500068', 'Siti Nurjanah'),
    ('c0000000-0000-0000-0000-000000000003', 'QRIS', 'QRIS AstraPay', 'QRIS', NULL, 'FAN Laundry Kalibaru')
ON CONFLICT (code) DO NOTHING;

INSERT INTO discounts (id, name, discount_type, value, min_order_amount)
VALUES
    ('d0000000-0000-0000-0000-000000000001', 'Promo Grand Opening 10%', 'PERCENTAGE', 10.00, 20000.00),
    ('d0000000-0000-0000-0000-000000000002', 'Loyalitas Cuci 9x Gratis 1x', 'NOMINAL', 20000.00, 0.00)
ON CONFLICT DO NOTHING;
