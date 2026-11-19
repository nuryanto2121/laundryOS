-- name: CreateOrder :one
INSERT INTO orders (
    id, order_number, customer_id, service_id, status,
    total_weight, total_price, pickup_required, notes, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW()
)
RETURNING *;

-- name: GetOrder :one
SELECT * FROM orders
WHERE id = $1
LIMIT 1;

-- name: ListOrders :many
SELECT * FROM orders
ORDER BY created_at DESC;

-- name: UpdateOrder :one
UPDATE orders
SET
    customer_id     = $2,
    service_id      = $3,
    total_weight    = $4,
    total_price     = $5,
    pickup_required = $6,
    notes           = $7,
    updated_at      = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateOrderStatus :one
UPDATE orders
SET
    status     = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteOrder :exec
DELETE FROM orders
WHERE id = $1;

-- name: GetLatestOrderNumberToday :one
SELECT order_number FROM orders
WHERE created_at >= CURRENT_DATE
ORDER BY order_number DESC
LIMIT 1;

-- name: GetLatestOrderByCustomerPhone :one
SELECT 
    o.id, o.order_number, o.customer_id, o.service_id, o.status,
    o.total_weight, o.total_price, o.pickup_required, o.notes,
    o.created_at, o.updated_at,
    s.service_name, s.duration_day,
    c.full_name AS customer_name
FROM orders o
JOIN customers c ON o.customer_id = c.id
JOIN services s ON o.service_id = s.id
WHERE c.phone = $1
ORDER BY o.created_at DESC
LIMIT 1;

-- name: GetOrderByOrderNumber :one
SELECT 
    o.id, o.order_number, o.customer_id, o.service_id, o.status,
    o.total_weight, o.total_price, o.pickup_required, o.notes,
    o.created_at, o.updated_at,
    s.service_name, s.duration_day,
    c.full_name AS customer_name
FROM orders o
JOIN customers c ON o.customer_id = c.id
JOIN services s ON o.service_id = s.id
WHERE o.order_number = $1
LIMIT 1;
