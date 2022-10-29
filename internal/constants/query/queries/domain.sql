-- name: CreateDomain :one 
INSERT INTO domains (
    name ,
    service_id
)VALUES (
    $1,$2
) RETURNING *;

-- name: DeleteDomain :exec 
DELETE from domains 
WHERE id = $1 AND deleted_at IS NULL;

-- name: SoftDeleteDomain :one 
UPDATE domains set deleted_at = now() 
WHERE name = $1 AND service_id = $2 AND deleted_at IS NULL
RETURNING *;


-- name: GetDomainByServiceId :many
SELECT * FROM domains 
WHERE service_id = $1 AND deleted_at IS NULL;

-- name: IsDomainExist :one
SELECT * FROM domains 
WHERE service_id = $1 AND name = $2 AND deleted_at IS NULL;