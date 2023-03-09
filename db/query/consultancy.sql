-- name: CreateConsultancy :one
INSERT INTO consultancies (
  merchant_id,
  customer_id,
  cost
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetConsultancy :one
SELECT * FROM consultancies
WHERE id = $1 LIMIT 1;

-- name: ListConsultancies :many
SELECT * FROM consultancies
WHERE 
    merchant_id = $1 OR
    customer_id = $2
ORDER BY id
LIMIT $3
OFFSET $4;