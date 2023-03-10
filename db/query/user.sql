-- name: CreateUser :one
INSERT INTO users (
  username,
  hashed_password,
  full_name,
  email,
  phone_number,
  gender,
  birth_date
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE username = $1 LIMIT 1;

-- name: UpdateUser :one
UPDATE users
SET full_name = COALESCE(sqlc.narg(full_name), full_name),
    phone_number = COALESCE(sqlc.narg(phone_number), phone_number),
    gender = COALESCE(sqlc.narg(gender), gender), 
    birth_date = COALESCE(sqlc.narg(birth_date), birth_date),
    image_url = COALESCE(sqlc.narg(image_url), image_url)
WHERE username = sqlc.arg(username)
RETURNING *;

-- name: UpdatePassword :one
UPDATE users
SET hashed_password = $2
WHERE username = $1
RETURNING *;

-- name: UpdateEmail :one
UPDATE users
SET email = $2
WHERE username = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE username = $1;