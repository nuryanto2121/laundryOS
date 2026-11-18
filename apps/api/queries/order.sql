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

-- name: DeleteOrder :exec
DELETE FROM orders
WHERE id = $1;

-- name: GetLatestOrderNumberToday :one
SELECT order_number FROM orders
WHERE created_at >= CURRENT_DATE
ORDER BY order_number DESC
LIMIT 1;
