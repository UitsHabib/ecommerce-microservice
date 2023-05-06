-- name: GetFeature :one
SELECT * FROM features
WHERE id = $1 LIMIT 1;

-- name: ListFeatures :many
SELECT * FROM features 
ORDER BY id
LIMIT $1
OFFSET $2;