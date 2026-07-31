-- name: CreateCustomer :one
INSERT INTO customers (
    id, customer_code, full_name, phone, email, address,
    latitude, longitude, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW()
)
RETURNING *;

-- name: GetCustomer :one
SELECT * FROM customers
WHERE id = $1 AND deleted_at IS NULL
LIMIT 1;

-- name: ListCustomers :many
SELECT * FROM customers
WHERE deleted_at IS NULL
ORDER BY created_at DESC;

-- name: UpdateCustomer :one
UPDATE customers
SET
    full_name   = $2,
    phone       = $3,
    email       = $4,
    address     = $5,
    latitude    = $6,
    longitude   = $7,
    updated_at  = NOW()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteCustomer :exec
UPDATE customers
SET deleted_at = NOW()
WHERE id = $1 AND deleted_at IS NULL;
