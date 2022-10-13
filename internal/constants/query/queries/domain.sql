-- name: CreateDomain :one 
INSERT INTO domains (
    name ,
    service_id
)VALUES (
    $1,$2
) RETURNING *;

-- name: DeleteDomain :exec 
DELETE from domains 
WHERE id = $1;


-- name: GetDomainByServiceId :many
SELECT * FROM domains 
WHERE service_id = $1;

-- name: IsDomainExist :one
SELECT * FROM domains 
WHERE service_id = $1 AND name = $2;