-- name: CreateComment :one
INSERT INTO comments (
  comment_type,
  post_id,
  merchant_id,
  owner,
  comment 
) VALUES (
  $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetComment :one
SELECT * FROM comments
WHERE id = $1 LIMIT 1;

-- name: ListPostComments :many
SELECT * FROM comments
WHERE post_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: ListMerchantComments :many
SELECT * FROM comments
WHERE merchant_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: DeleteComment :exec
DELETE FROM comments
WHERE id = $1;