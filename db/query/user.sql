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
SET full_name = COALESCE($2, full_name), phone_number = COALESCE($3, phone_number),
gender = COALESCE($4, gender), birth_date = COALESCE($5, birth_date)
WHERE username = $1
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