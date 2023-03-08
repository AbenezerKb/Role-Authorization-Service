-- name: CreateRole :one 
with _tenant as(
    select tenants.id as tenant_id from tenants where tenants.tenant_name=$1 AND tenants.service_id=$2 AND tenants.deleted_at IS NULL
), _role as(
    insert into roles(name, tenant_id) select $3,tenant_id from _tenant
        returning  id as role_id,name,created_at,status
),_rp as(
    insert into role_permissions (role_id,permission_id)
        select role_id,permissions.id from _role, permissions where permissions.id =ANY($4::uuid[]) and permissions.deleted_at IS NULL ON CONFLICT DO NOTHING returning id
)select _role.* from _role;

-- name: GetRoleByNameAndTenantName :one 
SELECT roles.id FROM roles join tenants on roles.tenant_id =tenants.id WHERE 
roles.name = $1 AND tenants.tenant_name = $2 and roles.deleted_at IS NULL and tenants.deleted_at IS NULL;



-- name: AssignRole :exec
WITH new_row AS (
    INSERT INTO users (user_id,  service_id) select $1,$2 where not exists(select id from users where user_id=$1 and service_id=$2 and deleted_at is null)returning id
),_user as(
    SELECT id FROM new_row
    UNION
    SELECT id FROM users WHERE user_id=$1 and service_id=$2 and deleted_at is null
)
insert into tenant_users_roles(tenant_id, user_id, role_id)
select tenants.id,_user.id,roles.id
from roles,tenants,_user  where tenants.tenant_name=$5 and tenants.deleted_at IS NULL and roles.id=$3 or roles.name=$4 and roles.tenant_id=tenants.id and roles.deleted_at is null;

-- name: IsRoleAssigned :one 
SELECT count_rows() FROM tenant_users_roles 
WHERE tenant_users_roles.tenant_id in (
    SELECT tenants.id FROM 
    tenants where tenants.tenant_name = $1 and tenants.deleted_at IS NULL
)
and tenant_users_roles.user_id in (
    SELECT users.id from users 
    where users.user_id = $2 and users.deleted_at IS NULL
) and tenant_users_roles.role_id = $3;

-- name: RemovePermissionsFromRole :exec
DELETE FROM role_permissions WHERE role_id=$1 AND NOT permission_id=any($2::uuid[]);

-- name: UpdateRole :exec
INSERT INTO role_permissions (role_id,permission_id)
SELECT $1,permissions.id FROM permissions WHERE permissions.id =ANY($2::uuid[]) AND permissions.deleted_at IS NULL ON conflict do nothing;

-- name: RevokeUserRole :exec 
UPDATE tenant_users_roles 
SET deleted_at= now() WHERE tenant_users_roles.tenant_id = (
    SELECT tenants.id FROM 
    tenants where tenants.tenant_name = $1 and tenants.deleted_at IS NULL
)
and tenant_users_roles.user_id = (
    SELECT users.id from users 
    where users.user_id = $2 and users.deleted_at IS NULL
) and tenant_users_roles.role_id = $3;

-- name: DeleteRole :one
Update roles set deleted_at=now() where roles.id=$1 AND deleted_at IS NULL returning name,id,created_at,updated_at;

-- name: ListRoles :many
select r.name,r.created_at,r.id,r.status from roles r join tenants t on r.tenant_id=t.id where t.tenant_name=$1 AND t.service_id=$2 AND t.deleted_at IS NULL AND r.deleted_at IS NULL;

-- name: UpdateRoleStatus :one
with _tenants as(
    select id from tenants t where t.tenant_name=$1 and t.service_id=$2 and t.deleted_at IS NULL
)
update roles r set status =$3 from _tenants where r.id=$4 and r.deleted_at IS NULL and r.tenant_id=_tenants.id returning r.id;

-- name: GetRoleById :one
select r.name,r.id,r.status,r.created_at,r.updated_at, (select string_to_array(string_agg(p.name,','),',')::string[] from role_permissions join permissions p on role_permissions.permission_id = p.id where role_id=r.id and p.deleted_at is null) as permission  from roles r join tenants t on t.id = r.tenant_id where t.service_id=$1 and t.deleted_at is null and r.id=$2 and r.deleted_at is null;

-- name: RevokeAdminRole :exec
UPDATE tenant_users_roles
SET status = 'INACTIVE'
WHERE tenant_id = $1 AND role_id = 'admin';