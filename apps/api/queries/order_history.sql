-- name: CreateOrderStatusHistory :one
INSERT INTO order_status_history (
    id, order_id, previous_status, new_status, changed_by, notes, created_at
) VALUES (
    $1, $2, $3, $4, $5, $6, NOW()
)
RETURNING *;

-- name: ListOrderStatusHistoryByOrderId :many
SELECT * FROM order_status_history
WHERE order_id = $1
ORDER BY created_at ASC;
