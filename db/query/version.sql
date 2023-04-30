-- name: CreateAppVersion :one
INSERT INTO app_versions (
  tag,
  version 
) VALUES (
  $1, $2
) RETURNING *;

-- name: GetAppVersion :one
SELECT * FROM app_versions
WHERE tag = $1 LIMIT 1;

-- name: ListAppVersions :many
SELECT * FROM app_versions
ORDER BY created_at DESC;

-- name: UpdateAppVersion :one
UPDATE app_versions
SET tag = $1, version = $2
WHERE id = $3
RETURNING *;
