-- name: ListActiveServices :many
SELECT id, service_code, service_name, duration_day, price_per_kg, active
FROM services
WHERE active = true
ORDER BY price_per_kg ASC;

-- name: ListAllServices :many
SELECT id, service_code, service_name, duration_day, price_per_kg, active
FROM services
ORDER BY service_name ASC;

-- name: CreateService :one
INSERT INTO services (id, service_code, service_name, duration_day, price_per_kg, active)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: UpdateService :one
UPDATE services
SET service_name = $2, duration_day = $3, price_per_kg = $4, active = $5
WHERE id = $1
RETURNING *;

-- name: DeleteService :exec
UPDATE services
SET active = false
WHERE id = $1;
