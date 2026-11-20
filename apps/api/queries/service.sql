-- name: ListActiveServices :many
SELECT id, service_code, service_name, duration_day, price_per_kg, active, category_id, category_name, unit
FROM services
WHERE active = true
ORDER BY price_per_kg ASC;

-- name: ListAllServices :many
SELECT id, service_code, service_name, duration_day, price_per_kg, active, category_id, category_name, unit
FROM services
ORDER BY
  CASE WHEN sqlc.narg('sort_by')::VARCHAR = 'service_code' AND sqlc.narg('sort_dir')::VARCHAR = 'asc' THEN service_code END ASC,
  CASE WHEN sqlc.narg('sort_by')::VARCHAR = 'service_code' AND sqlc.narg('sort_dir')::VARCHAR = 'desc' THEN service_code END DESC,
  CASE WHEN sqlc.narg('sort_by')::VARCHAR = 'service_name' AND sqlc.narg('sort_dir')::VARCHAR = 'asc' THEN service_name END ASC,
  CASE WHEN sqlc.narg('sort_by')::VARCHAR = 'service_name' AND sqlc.narg('sort_dir')::VARCHAR = 'desc' THEN service_name END DESC,
  CASE WHEN sqlc.narg('sort_by')::VARCHAR = 'price' AND sqlc.narg('sort_dir')::VARCHAR = 'asc' THEN price_per_kg END ASC,
  CASE WHEN sqlc.narg('sort_by')::VARCHAR = 'price' AND sqlc.narg('sort_dir')::VARCHAR = 'desc' THEN price_per_kg END DESC,
  CASE WHEN sqlc.narg('sort_by')::VARCHAR = 'category_name' AND sqlc.narg('sort_dir')::VARCHAR = 'asc' THEN category_name END ASC,
  CASE WHEN sqlc.narg('sort_by')::VARCHAR = 'category_name' AND sqlc.narg('sort_dir')::VARCHAR = 'desc' THEN category_name END DESC,
  category_name ASC, service_name ASC;

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
