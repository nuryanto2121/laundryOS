-- name: CreateOrderItem :one
INSERT INTO order_items (
    id, order_id, service_id, service_name, category,
    quantity, unit, price_per_unit, subtotal, notes, created_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW()
)
RETURNING *;

-- name: ListOrderItemsByOrderId :many
SELECT * FROM order_items
WHERE order_id = $1
ORDER BY created_at ASC;

-- name: DeleteOrderItemsByOrderId :exec
DELETE FROM order_items
WHERE order_id = $1;
