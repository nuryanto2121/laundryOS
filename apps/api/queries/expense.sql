-- name: CreateExpense :one
INSERT INTO expenses (
    id, outlet_id, expense_date, category, amount, description, created_by, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, NOW(), NOW()
)
RETURNING *;

-- name: ListExpenses :many
SELECT * FROM expenses
WHERE expense_date BETWEEN $1 AND $2
ORDER BY expense_date DESC, created_at DESC;

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
