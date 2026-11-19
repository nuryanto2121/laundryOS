-- name: GetHomeMetricsToday :one
SELECT 
    COALESCE(SUM(CASE WHEN created_at >= CURRENT_DATE AND status != 'BATAL' THEN total_price ELSE 0 END), 0)::DECIMAL(12,2) AS omset_hari_ini,
    COUNT(CASE WHEN created_at >= CURRENT_DATE AND status != 'BATAL' THEN 1 END)::BIGINT AS order_masuk_hari_ini,
    COUNT(CASE WHEN estimated_completion_at::DATE = CURRENT_DATE AND status NOT IN ('SELESAI', 'BATAL') THEN 1 END)::BIGINT AS harus_selesai_hari_ini,
    COUNT(CASE WHEN estimated_completion_at < NOW() AND status NOT IN ('SELESAI', 'BATAL') THEN 1 END)::BIGINT AS terlambat_hari_ini
FROM orders;

-- name: GetDailyRevenueChart :many
SELECT 
    DATE(created_at) AS report_date,
    COALESCE(SUM(total_price), 0)::DECIMAL(12,2) AS omset,
    COALESCE(SUM(paid_amount), 0)::DECIMAL(12,2) AS pendapatan_kas,
    COUNT(*)::BIGINT AS total_transaksi
FROM orders
WHERE created_at::DATE BETWEEN $1 AND $2
  AND status != 'BATAL'
GROUP BY DATE(created_at)
ORDER BY report_date ASC;

-- name: ListUnpaidOrders :many
SELECT 
    o.id, o.order_number, o.customer_id, o.status,
    o.total_price, o.paid_amount,
    (o.total_price - o.paid_amount)::DECIMAL(12,2) AS sisa_tagihan,
    o.payment_status, o.rack_number, o.created_at,
    c.full_name AS customer_name,
    c.phone AS customer_phone
FROM orders o
JOIN customers c ON o.customer_id = c.id
WHERE o.payment_status != 'PAID'
  AND o.status != 'BATAL'
ORDER BY o.created_at ASC;

-- name: GetServiceVolumeBreakdown :many
SELECT 
    oi.category,
    oi.unit,
    COALESCE(SUM(oi.quantity), 0)::DECIMAL(10,2) AS total_volume,
    COALESCE(SUM(oi.subtotal), 0)::DECIMAL(12,2) AS total_pendapatan,
    COUNT(DISTINCT oi.order_id)::BIGINT AS total_order
FROM order_items oi
JOIN orders o ON oi.order_id = o.id
WHERE o.created_at::DATE BETWEEN $1 AND $2
  AND o.status != 'BATAL'
GROUP BY oi.category, oi.unit
ORDER BY total_pendapatan DESC;

-- name: GetTopCustomersLeaderboard :many
SELECT 
    c.id, c.customer_code, c.full_name, c.phone,
    COUNT(o.id)::BIGINT AS total_orders,
    COALESCE(SUM(o.total_price), 0)::DECIMAL(12,2) AS total_spend,
    COALESCE(SUM(o.total_weight), 0)::DECIMAL(10,2) AS total_kg,
    MAX(o.created_at) AS last_order_date
FROM customers c
JOIN orders o ON c.id = o.customer_id
WHERE o.status != 'BATAL'
GROUP BY c.id, c.customer_code, c.full_name, c.phone
ORDER BY total_spend DESC
LIMIT $1;

-- name: GetCustomerCRMStats :one
SELECT 
    c.id, c.customer_code, c.full_name, c.phone, c.address, c.created_at,
    COUNT(o.id)::BIGINT AS total_orders,
    COALESCE(SUM(CASE WHEN o.status != 'BATAL' THEN o.total_price ELSE 0 END), 0)::DECIMAL(12,2) AS lifetime_spend,
    COALESCE(SUM(CASE WHEN o.status != 'BATAL' THEN o.total_weight ELSE 0 END), 0)::DECIMAL(10,2) AS lifetime_weight,
    COALESCE(SUM(CASE WHEN o.status != 'BATAL' AND o.payment_status != 'PAID' THEN (o.total_price - o.paid_amount) ELSE 0 END), 0)::DECIMAL(12,2) AS total_unpaid
FROM customers c
LEFT JOIN orders o ON c.id = o.customer_id
WHERE c.id = $1
GROUP BY c.id, c.customer_code, c.full_name, c.phone, c.address, c.created_at;
