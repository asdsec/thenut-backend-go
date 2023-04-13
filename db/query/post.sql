-- name: CreatePost :one
INSERT INTO posts (
  merchant_id,
  title,
  image_url
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetPost :one
SELECT * FROM posts
WHERE id = $1 LIMIT 1;

-- name: ListMerchantPosts :many
SELECT * FROM posts
WHERE merchant_id = $1
ORDER BY created_at DESC
LIMIT $2
OFFSET $3;

-- name: ListPosts :many
SELECT * FROM posts
ORDER BY created_at DESC
LIMIT $1
OFFSET $2;

-- name: DeletePost :exec
DELETE FROM posts
WHERE id = $1;