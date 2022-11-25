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
    select tenants.domain_id,tenants.id,tenants.inherit from tenants where tenant_name =$1 and tenants.service_id=$2 and tenants.deleted_at IS NULL
)
select p.name,p.status,p.description,p.statement,p.created_at,p.service_id,p.id,coalesce(json_agg(json_build_object('name',p2.name,'status',p2.status,'description',p2.description,'statement',p2.statement,'created_at',p2.created_at,'service_id',p2.service_id,'id',p2.id)) filter ( where  p2.deleted_at is null and p2.status='ACTIVE'  ),'[]') AS inherited_permissions from _tenant,permissions p  left join permission_domains pd on p.id = pd.permission_id left join permissions_hierarchy ph on p.id = ph.parent left join permissions p2 on p2.id = ph.child where pd.domain_id = _tenant.domain_id and _tenant.inherit = true group by p.name,p.status,p.description,p.statement,p.created_at,p.service_id,p.id
UNION
select p.name,p.status,p.description,p.statement,p.created_at,p.service_id,p.id,coalesce(json_agg(json_build_object('name',p2.name,'status',p2.status,'description',p2.description,'statement',p2.statement,'created_at',p2.created_at,'service_id',p2.service_id,'id',p2.id))filter ( where  p2.deleted_at is null and p2.status='ACTIVE'  ),'[]') AS inherited_permissions  from permissions p left join permissions_hierarchy ph on p.id = ph.parent left join permissions p2 on p2.id = ph.child,_tenant where p.tenant_id =_tenant.id and p.deleted_at IS NULL group by p.name,p.status,p.description,p.statement,p.created_at,p.service_id,p.id;

-- name: CreatePermissionDependency :exec
with _parent as(
    select id as parant_id from permissions where permissions.name=$1 and permissions.service_id=$2 and permissions.deleted_at IS NULL
)
insert into permissions_hierarchy(parent, child) select _parent.parant_id,p.id from _parent,permissions p where p.name=ANY($3::string[]) and p.service_id=$2  and p.deleted_at IS NULl ON conflict  do nothing;

-- name: DeletePermissions :one 
UPDATE permissions p set deleted_at = now() from tenants t WHERE  t.tenant_name=$1
and p.id = $2 and p.tenant_id=t.id AND p.service_id = $3 AND t.service_id = $3
AND p.deleted_at IS NULL AND p.delete_or_update RETURNING p.id;

-- name: CanBeDeleted :one
select p.delete_or_update from permissions p where p.id=$1 and p.service_id=$2;


-- name: GetPermissionDetails :one
WITH _tenant AS(
    SELECT
        domain_id,
        id AS tenant_id
    FROM
        tenants
    WHERE
            tenant_name = $1
)
SELECT
    p.name,
    p.status,
    p.description,
    p.statement,
    p.id,
    COALESCE(json_agg(json_build_object('name',
                                        p2.name,
                                        'description',
                                        p2.description,
                                        'statement',
                                        p2.statement,
                                        'id',
                                        p2.id)) FILTER (
                 WHERE
                     p2.deleted_at IS NULL
                     AND p2.status = 'ACTIVE' ),
             '[]') AS inherited_permissions
FROM
    _tenant,
    permissions p
        LEFT JOIN permissions_hierarchy ph ON
            p.id = ph.parent
        LEFT JOIN permissions p2 ON
            p2.id = ph.child
        LEFT JOIN permission_domains pd ON
            p.id = pd.permission_id
WHERE
    (p.tenant_id = _tenant.tenant_id
        OR pd.domain_id = _tenant.domain_id)
  AND p.id = $2
  AND p.service_id = $3
  AND p.deleted_at IS NULL
GROUP BY
    p.name,
    p.status,
    p.description,
    p.statement,
    p.id;
