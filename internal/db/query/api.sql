-- name: CreateAPI :one
INSERT INTO apis (
    key
) VALUES (
    $1
) RETURNING *;

-- name: GetAPI :one
SELECT * FROM apis
ORDER BY RANDOM()
LIMIT 1;

-- name: UpdateAPIUsageCount :exec
UPDATE apis
SET usage_count = $2
WHERE id = $1;