-- name: CreateOrGetPermission :one
WITH new_row AS (
    INSERT INTO permissions (name,description,statment,service_id)
        SELECT $1,$2,$3,$4
        WHERE NOT EXISTS (SELECT id FROM permissions WHERE name = $1 and service_id=$4)
        RETURNING id
)
SELECT id FROM new_row
UNION
SELECT id FROM permissions WHERE name = $1 and service_id=$4;

-- name: AssignDomain :exec
with _domain as(
    select domains.id as domain_id from domains where domains.id=$1
) INSERT INTO permission_domains(domain_id,permission_id)
SELECT  domain_id, $2 from _domain
WHERE NOT exists(select permission_id from permission_domains where permission_id=$2 and domain_id=_domain.domain_id);

