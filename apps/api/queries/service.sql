-- name: ListActiveServices :many
SELECT id, service_code, service_name, duration_day, price_per_kg, active
FROM services
WHERE active = true
ORDER BY price_per_kg ASC;
