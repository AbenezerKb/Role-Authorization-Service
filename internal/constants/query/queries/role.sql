-- name: CreateRole :one 
with _tenant as(
    select tenants.id as tenant_id from tenants where tenants.tenant_name=$1 AND tenants.service_id=$2
), _role as(
    insert into roles(name, tenant_id) select $3,tenant_id from _tenant
    returning  id as role_id,name,created_at,status
),_rp as(
insert into role_permissions (role_id,permission_id)
    select role_id,permissions.id from _role, permissions where permissions.id =ANY($4::uuid[])returning id
)select _role.* from _role;

-- name: GetRoleByNameAndTenantName :one 
SELECT roles.id FROM roles join tenants on roles.tenant_id =tenants.id WHERE 
roles.name = $1 AND tenants.tenant_name = $2;

