-- name: CreateOrGetPermission :one
WITH new_row AS (
    INSERT INTO permissions (name,description,statement,service_id)
        SELECT $1,$2,$3,$4
        WHERE NOT EXISTS (SELECT id FROM permissions WHERE name =$1 and service_id=$4 and deleted_at IS NULL)
        RETURNING id
),_permission as(
    SELECT id FROM new_row
    UNION
    SELECT id FROM permissions WHERE name = $1 and service_id=$4 and deleted_at IS NULL
)
,pd as (insert into permission_domains (domain_id,permission_id)
select domains.id,_permission.id from domains,_permission where domains.id =ANY($5::uuid[]) ON CONFLICT DO NOTHING returning permission_domains.id
)
select id from _permission;

-- name: AssignDomain :exec
with _domain as(
    select domains.id as domain_id from domains where domains.id=$1 and deleted_at IS NULL
) INSERT INTO permission_domains(domain_id,permission_id)
SELECT  domain_id, $2 from _domain
WHERE NOT exists(select permission_id from permission_domains where permission_id=$2 and domain_id=_domain.domain_id and deleted_at IS NULL);

-- name: ListPermissions :many
with _tenant as (
    select tenants.domain_id,tenants.id,tenants.inherit from tenants where tenant_name =$1 and tenants.service_id=$2 and deleted_at IS NULL
)
select p.name,p.status,p.description,p.statement,p.created_at,p.service_id,p.id  from _tenant,permissions p  join permission_domains pd on p.id = pd.permission_id where pd.domain_id = _tenant.domain_id and _tenant.inherit = true
UNION
select p.name,p.status,p.description,p.statement,p.created_at,p.service_id,p.id  from permissions p,_tenant where p.tenant_id =_tenant.id and deleted_at IS NULL;