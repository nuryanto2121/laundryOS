-- name: ListActiveServices :many
SELECT id, service_code, service_name, duration_day, price_per_kg, active, category_id, category_name, unit
FROM services
WHERE active = true
ORDER BY price_per_kg ASC;

-- name: ListAllServices :many
SELECT id, service_code, service_name, duration_day, price_per_kg, active, category_id, category_name, unit
FROM services
ORDER BY category_name ASC, service_name ASC;

-- name: CreateService :one
INSERT INTO services (id, service_code, service_name, duration_day, price_per_kg, active, category_id, category_name, unit)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: UpdateService :one
UPDATE services
SET service_name = $2, duration_day = $3, price_per_kg = $4, active = $5, category_id = $6, category_name = $7, unit = $8
WHERE id = $1
RETURNING *;

-- name: DeleteService :exec
UPDATE services
SET active = false
WHERE id = $1;
