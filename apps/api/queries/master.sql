-- name: ListActivePerfumes :many
SELECT * FROM perfumes
WHERE is_active = true
ORDER BY name ASC;

-- name: ListAllPerfumes :many
SELECT * FROM perfumes
ORDER BY name ASC;

-- name: CreatePerfume :one
INSERT INTO perfumes (id, name, description, is_active, created_at)
VALUES ($1, $2, $3, $4, NOW())
RETURNING *;

-- name: UpdatePerfume :one
UPDATE perfumes
SET name = $2, description = $3, is_active = $4
WHERE id = $1
RETURNING *;

-- name: DeletePerfume :exec
UPDATE perfumes
SET is_active = false
WHERE id = $1;

-- name: ListActiveDiscounts :many
SELECT * FROM discounts
WHERE is_active = true
ORDER BY name ASC;

-- name: ListAllDiscounts :many
SELECT * FROM discounts
ORDER BY name ASC;

-- name: CreateDiscount :one
INSERT INTO discounts (id, name, discount_type, value, min_order_amount, is_active, created_at)
VALUES ($1, $2, $3, $4, $5, $6, NOW())
RETURNING *;

-- name: UpdateDiscount :one
UPDATE discounts
SET name = $2, discount_type = $3, value = $4, min_order_amount = $5, is_active = $6
WHERE id = $1
RETURNING *;

-- name: DeleteDiscount :exec
UPDATE discounts
SET is_active = false
WHERE id = $1;

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

-- name: ListExpenseCategories :many
SELECT * FROM expense_categories
WHERE is_active = true
ORDER BY name ASC;

-- name: ListAllExpenseCategories :many
SELECT * FROM expense_categories
ORDER BY
  CASE WHEN sqlc.narg('sort_by')::VARCHAR = 'code' AND sqlc.narg('sort_dir')::VARCHAR = 'asc' THEN code END ASC,
  CASE WHEN sqlc.narg('sort_by')::VARCHAR = 'code' AND sqlc.narg('sort_dir')::VARCHAR = 'desc' THEN code END DESC,
  CASE WHEN sqlc.narg('sort_by')::VARCHAR = 'name' AND sqlc.narg('sort_dir')::VARCHAR = 'asc' THEN name END ASC,
  CASE WHEN sqlc.narg('sort_by')::VARCHAR = 'name' AND sqlc.narg('sort_dir')::VARCHAR = 'desc' THEN name END DESC,
  name ASC;

-- name: CreateExpenseCategory :one
INSERT INTO expense_categories (id, code, name, is_active, created_at)
VALUES ($1, $2, $3, $4, NOW())
RETURNING *;

-- name: UpdateExpenseCategory :one
UPDATE expense_categories
SET name = $2, is_active = $3
WHERE id = $1
RETURNING *;

-- name: DeleteExpenseCategory :exec
UPDATE expense_categories
SET is_active = false
WHERE id = $1;

-- name: ListEmployees :many
SELECT * FROM employees
WHERE is_active = true
ORDER BY role ASC, full_name ASC;

-- name: ListAllEmployees :many
SELECT * FROM employees
ORDER BY created_at DESC;

-- name: CreateEmployee :one
INSERT INTO employees (id, username, full_name, phone, role, pin, is_active, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
RETURNING *;

-- name: UpdateEmployee :one
UPDATE employees
SET full_name = $2, phone = $3, role = $4, pin = $5, is_active = $6, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteEmployee :exec
UPDATE employees
SET is_active = false, updated_at = NOW()
WHERE id = $1;

-- name: ListActiveServiceCategories :many
SELECT * FROM service_categories
WHERE is_active = true
ORDER BY display_order ASC, name ASC;

-- name: ListAllServiceCategories :many
SELECT * FROM service_categories
ORDER BY display_order ASC, name ASC;

-- name: GetServiceCategory :one
SELECT * FROM service_categories
WHERE id = $1;

-- name: CreateServiceCategory :one
INSERT INTO service_categories (id, code, name, default_unit, display_order, is_active, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
RETURNING *;

-- name: UpdateServiceCategory :one
UPDATE service_categories
SET name = $2, default_unit = $3, display_order = $4, is_active = $5, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteServiceCategory :exec
UPDATE service_categories
SET is_active = false, updated_at = NOW()
WHERE id = $1;

