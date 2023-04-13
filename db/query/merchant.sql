-- name: CreateMerchant :one
INSERT INTO merchants (
  owner,
  profession,
  title,
  about
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: GetMerchant :one
SELECT * FROM merchants
WHERE id = $1 LIMIT 1;

-- name: ListMerchants :many
SELECT * FROM merchants
WHERE owner = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: UpdateMerchant :one
UPDATE merchants
SET balance = COALESCE(sqlc.narg(balance), balance),
    profession = COALESCE(sqlc.narg(profession), profession),
    title = COALESCE(sqlc.narg(title), title),
    about = COALESCE(sqlc.narg(about), about),
    image_url = COALESCE(sqlc.narg(image_url), image_url),
    rating = COALESCE(sqlc.narg(rating), rating)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: AddMerchantBalance :one
UPDATE merchants
SET balance = balance + sqlc.arg(amount)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteMerchant :exec
DELETE FROM merchants
WHERE id = $1;