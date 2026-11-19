-- name: CreateOrderPhoto :one
INSERT INTO order_photos (
    id, order_id, photo_type, photo_url, notes, created_at
) VALUES (
    $1, $2, $3, $4, $5, NOW()
)
RETURNING *;

-- name: ListOrderPhotosByOrderId :many
SELECT * FROM order_photos
WHERE order_id = $1
ORDER BY created_at ASC;
