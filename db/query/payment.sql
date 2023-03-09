-- name: CreatePayment :one
INSERT INTO payments (
  merchant_id,
  customer_id,
  amount
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetPayment :one
SELECT * FROM payments
WHERE id = $1 LIMIT 1;

-- name: ListPayments :many
SELECT * FROM payments
WHERE 
    merchant_id = $1 OR
    customer_id = $2
ORDER BY id
LIMIT $3
OFFSET $4;