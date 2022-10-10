-- name: GetServiceByName :one
SELECT * FROM services WHERE name = $1;

-- name: CreateService :one
INSERT INTO services (
    name,
    password
) VALUES (
    $1, $2
) RETURNING *;

-- name: DeleteService :exec
DELETE FROM services WHERE id = $1;