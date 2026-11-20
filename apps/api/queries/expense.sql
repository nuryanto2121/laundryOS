-- name: CreateExpense :one
INSERT INTO expenses (
    id, outlet_id, expense_date, category, amount, description, created_by, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, NOW(), NOW()
)
RETURNING *;

-- name: UpdateExpense :one
UPDATE expenses
SET
    expense_date = $2,
    category     = $3,
    amount       = $4,
    description  = $5,
    updated_at   = NOW()
WHERE id = $1
RETURNING *;

-- name: ListExpenses :many
SELECT * FROM expenses
WHERE expense_date BETWEEN $1 AND $2
ORDER BY expense_date DESC, created_at DESC;

-- name: ListExpensesPaginated :many
SELECT * FROM expenses
WHERE expense_date BETWEEN $1 AND $2
  AND (sqlc.narg('category')::VARCHAR IS NULL OR category = sqlc.narg('category'))
ORDER BY
  CASE WHEN sqlc.narg('sort_by')::VARCHAR = 'expense_date' AND sqlc.narg('sort_dir')::VARCHAR = 'asc' THEN expense_date END ASC,
  CASE WHEN sqlc.narg('sort_by')::VARCHAR = 'expense_date' AND sqlc.narg('sort_dir')::VARCHAR = 'desc' THEN expense_date END DESC,
  CASE WHEN sqlc.narg('sort_by')::VARCHAR = 'category' AND sqlc.narg('sort_dir')::VARCHAR = 'asc' THEN category END ASC,
  CASE WHEN sqlc.narg('sort_by')::VARCHAR = 'category' AND sqlc.narg('sort_dir')::VARCHAR = 'desc' THEN category END DESC,
  CASE WHEN sqlc.narg('sort_by')::VARCHAR = 'amount' AND sqlc.narg('sort_dir')::VARCHAR = 'asc' THEN amount END ASC,
  CASE WHEN sqlc.narg('sort_by')::VARCHAR = 'amount' AND sqlc.narg('sort_dir')::VARCHAR = 'desc' THEN amount END DESC,
  CASE WHEN sqlc.narg('sort_by')::VARCHAR = 'created_by' AND sqlc.narg('sort_dir')::VARCHAR = 'asc' THEN created_by END ASC,
  CASE WHEN sqlc.narg('sort_by')::VARCHAR = 'created_by' AND sqlc.narg('sort_dir')::VARCHAR = 'desc' THEN created_by END DESC,
  CASE WHEN sqlc.narg('sort_by')::VARCHAR = 'created_at' AND sqlc.narg('sort_dir')::VARCHAR = 'asc' THEN created_at END ASC,
  CASE WHEN sqlc.narg('sort_by')::VARCHAR = 'created_at' AND sqlc.narg('sort_dir')::VARCHAR = 'desc' THEN created_at END DESC,
  expense_date DESC, created_at DESC
LIMIT $3 OFFSET $4;

-- name: CountExpenses :one
SELECT COUNT(*) FROM expenses
WHERE expense_date BETWEEN $1 AND $2
  AND (sqlc.narg('category')::VARCHAR IS NULL OR category = sqlc.narg('category'));

-- name: GetExpense :one
SELECT * FROM expenses
WHERE id = $1
LIMIT 1;

-- name: DeleteExpense :exec
DELETE FROM expenses
WHERE id = $1;

-- name: GetExpenseSummary :many
SELECT category, COALESCE(SUM(amount), 0)::DECIMAL(12,2) AS total_amount, COUNT(*)::BIGINT AS count
FROM expenses
WHERE expense_date BETWEEN $1 AND $2
GROUP BY category
ORDER BY total_amount DESC;

-- name: GetTotalExpenseByDateRange :one
SELECT COALESCE(SUM(amount), 0)::DECIMAL(12,2) AS total_expense
FROM expenses
WHERE expense_date BETWEEN $1 AND $2;
