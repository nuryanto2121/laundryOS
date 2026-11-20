-- name: CreatePOSOrder :one
INSERT INTO orders (
    id, order_number, customer_id, service_id, status,
    total_weight, total_price, pickup_required, notes,
    outlet_id, rack_number, perfume_id, perfume_name,
    discount_id, discount_amount, paid_amount, payment_status,
    payment_method_id, payment_method_name, paid_at,
    estimated_completion_at, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5,
    $6, $7, $8, $9,
    $10, $11, $12, $13,
    $14, $15, $16, $17,
    $18, $19, $20,
    $21, NOW(), NOW()
)
RETURNING *;

-- name: UpdatePOSOrderPayment :one
UPDATE orders
SET
    paid_amount         = $2,
    payment_status      = $3,
    payment_method_id   = $4,
    payment_method_name = $5,
    paid_at             = $6,
    updated_at          = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdatePOSOrderStatus :one
UPDATE orders
SET
    status       = sqlc.arg(status)::VARCHAR,
    completed_at = CASE WHEN sqlc.arg(status)::VARCHAR = 'SELESAI' THEN NOW() ELSE completed_at END,
    updated_at   = NOW()
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: CancelPOSOrder :one
UPDATE orders
SET
    status          = 'BATAL',
    canceled_reason = $2,
    canceled_at     = NOW(),
    updated_at      = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdatePOSOrderRack :one
UPDATE orders
SET
    rack_number = $2,
    updated_at  = NOW()
WHERE id = $1
RETURNING *;

-- name: ListPOSOrders :many
SELECT 
    o.id, o.order_number, o.customer_id, o.service_id, o.status,
    o.total_weight, o.total_price, o.pickup_required, o.notes,
    o.outlet_id, o.rack_number, o.perfume_id, o.perfume_name,
    o.discount_id, o.discount_amount, o.paid_amount, o.payment_status,
    o.payment_method_id, o.payment_method_name, o.paid_at,
    o.estimated_completion_at, o.completed_at, o.canceled_reason, o.canceled_at,
    o.created_at, o.updated_at,
    c.full_name AS customer_name,
    c.phone AS customer_phone,
    COALESCE(s.service_name, '') AS primary_service_name
FROM orders o
JOIN customers c ON o.customer_id = c.id
LEFT JOIN services s ON o.service_id = s.id
WHERE (sqlc.narg('status')::VARCHAR IS NULL OR o.status = sqlc.narg('status'))
  AND (
    sqlc.narg('search')::VARCHAR IS NULL 
    OR o.order_number ILIKE '%' || sqlc.narg('search') || '%'
    OR c.full_name ILIKE '%' || sqlc.narg('search') || '%'
    OR c.phone ILIKE '%' || sqlc.narg('search') || '%'
    OR o.rack_number ILIKE '%' || sqlc.narg('search') || '%'
  )
ORDER BY o.created_at DESC;

-- name: ListPOSOrdersPaginated :many
SELECT 
    o.id, o.order_number, o.customer_id, o.service_id, o.status,
    o.total_weight, o.total_price, o.pickup_required, o.notes,
    o.outlet_id, o.rack_number, o.perfume_id, o.perfume_name,
    o.discount_id, o.discount_amount, o.paid_amount, o.payment_status,
    o.payment_method_id, o.payment_method_name, o.paid_at,
    o.estimated_completion_at, o.completed_at, o.canceled_reason, o.canceled_at,
    o.created_at, o.updated_at,
    c.full_name AS customer_name,
    c.phone AS customer_phone,
    COALESCE(s.service_name, '') AS primary_service_name
FROM orders o
JOIN customers c ON o.customer_id = c.id
LEFT JOIN services s ON o.service_id = s.id
WHERE (sqlc.narg('status')::VARCHAR IS NULL OR o.status = sqlc.narg('status'))
  AND (
    sqlc.narg('search')::VARCHAR IS NULL 
    OR o.order_number ILIKE '%' || sqlc.narg('search') || '%'
    OR c.full_name ILIKE '%' || sqlc.narg('search') || '%'
    OR c.phone ILIKE '%' || sqlc.narg('search') || '%'
    OR o.rack_number ILIKE '%' || sqlc.narg('search') || '%'
  )
ORDER BY
  CASE WHEN sqlc.narg('sort_by')::VARCHAR = 'order_number' AND sqlc.narg('sort_dir')::VARCHAR = 'asc' THEN o.order_number END ASC,
  CASE WHEN sqlc.narg('sort_by')::VARCHAR = 'order_number' AND sqlc.narg('sort_dir')::VARCHAR = 'desc' THEN o.order_number END DESC,
  CASE WHEN sqlc.narg('sort_by')::VARCHAR = 'customer_name' AND sqlc.narg('sort_dir')::VARCHAR = 'asc' THEN c.full_name END ASC,
  CASE WHEN sqlc.narg('sort_by')::VARCHAR = 'customer_name' AND sqlc.narg('sort_dir')::VARCHAR = 'desc' THEN c.full_name END DESC,
  CASE WHEN sqlc.narg('sort_by')::VARCHAR = 'total_price' AND sqlc.narg('sort_dir')::VARCHAR = 'asc' THEN o.total_price END ASC,
  CASE WHEN sqlc.narg('sort_by')::VARCHAR = 'total_price' AND sqlc.narg('sort_dir')::VARCHAR = 'desc' THEN o.total_price END DESC,
  CASE WHEN sqlc.narg('sort_by')::VARCHAR = 'status' AND sqlc.narg('sort_dir')::VARCHAR = 'asc' THEN o.status END ASC,
  CASE WHEN sqlc.narg('sort_by')::VARCHAR = 'status' AND sqlc.narg('sort_dir')::VARCHAR = 'desc' THEN o.status END DESC,
  CASE WHEN sqlc.narg('sort_by')::VARCHAR = 'created_at' AND sqlc.narg('sort_dir')::VARCHAR = 'asc' THEN o.created_at END ASC,
  CASE WHEN sqlc.narg('sort_by')::VARCHAR = 'created_at' AND sqlc.narg('sort_dir')::VARCHAR = 'desc' THEN o.created_at END DESC,
  o.created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountPOSOrders :one
SELECT COUNT(*)
FROM orders o
JOIN customers c ON o.customer_id = c.id
WHERE (sqlc.narg('status')::VARCHAR IS NULL OR o.status = sqlc.narg('status'))
  AND (
    sqlc.narg('search')::VARCHAR IS NULL 
    OR o.order_number ILIKE '%' || sqlc.narg('search') || '%'
    OR c.full_name ILIKE '%' || sqlc.narg('search') || '%'
    OR c.phone ILIKE '%' || sqlc.narg('search') || '%'
    OR o.rack_number ILIKE '%' || sqlc.narg('search') || '%'
  );

-- name: GetPOSOrderDetail :one
SELECT 
    o.id, o.order_number, o.customer_id, o.service_id, o.status,
    o.total_weight, o.total_price, o.pickup_required, o.notes,
    o.outlet_id, o.rack_number, o.perfume_id, o.perfume_name,
    o.discount_id, o.discount_amount, o.paid_amount, o.payment_status,
    o.payment_method_id, o.payment_method_name, o.paid_at,
    o.estimated_completion_at, o.completed_at, o.canceled_reason, o.canceled_at,
    o.created_at, o.updated_at,
    c.full_name AS customer_name,
    c.phone AS customer_phone,
    c.address AS customer_address,
    COALESCE(s.service_name, '') AS primary_service_name,
    COALESCE(outl.name, 'FAN Laundry Kalibaru') AS outlet_name,
    COALESCE(outl.address, '') AS outlet_address,
    COALESCE(outl.phone, '') AS outlet_phone
FROM orders o
JOIN customers c ON o.customer_id = c.id
LEFT JOIN services s ON o.service_id = s.id
LEFT JOIN outlets outl ON o.outlet_id = outl.id
WHERE o.id = $1
LIMIT 1;
