-- name: CreateMerchant :one
INSERT INTO merchants (
  owner,
  balance,
  profession,
  title,
  about
) VALUES (
  $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetMerchant :one
SELECT * FROM merchants
WHERE id = $1 LIMIT 1;

-- name: GetMerchantForUpdate :one
SELECT * FROM merchants
WHERE id = $1 LIMIT 1
FOR NO KEY UPDATE;

-- name: ListMerchants :many
SELECT * FROM merchants
WHERE owner = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: UpdateMerchant :one
UPDATE merchants
SET balance = COALESCE($2, balance), 
    profession = COALESCE($3, profession),
    title = COALESCE($4, title),
    about = COALESCE($5, about),
    image_url = COALESCE($6, image_url),
    rating = COALESCE($7, rating)
WHERE id = $1
RETURNING *;

-- name: AddMerchantBalance :one
UPDATE merchants
SET balance = balance + sqlc.arg(amount)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteMerchant :exec
DELETE FROM merchants
WHERE id = $1;