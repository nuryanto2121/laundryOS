-- name: ListActivePerfumes :many
SELECT * FROM perfumes
WHERE is_active = true
ORDER BY name ASC;

-- name: ListActiveDiscounts :many
SELECT * FROM discounts
WHERE is_active = true
ORDER BY name ASC;

-- name: ListActivePaymentMethods :many
SELECT * FROM payment_methods
WHERE is_active = true
ORDER BY category ASC, name ASC;

-- name: ListOutlets :many
SELECT * FROM outlets
WHERE is_active = true
ORDER BY name ASC;

-- name: GetDefaultOutlet :one
SELECT * FROM outlets
WHERE is_active = true
ORDER BY created_at ASC
LIMIT 1;
