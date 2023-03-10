-- name: CreateCustomer :one
INSERT INTO customers (
  owner
) VALUES (
  $1
) RETURNING *;

-- name: GetCustomer :one
SELECT * FROM customers
WHERE id = $1 LIMIT 1;

-- name: UpdateCustomer :one
UPDATE customers
SET image_url = $2
WHERE id = $1
RETURNING *;

-- name: DeleteCustomer :exec
DELETE FROM customers
WHERE id = $1;