-- name: CreateTenent :exec 
INSERT INTO tenants (
domain_id,
tenant_name,
service_id

) VALUES (
 $1,$2,$3
) ;

-- name: GetTenentWithNameAndServiceId :one 
SELECT * FROM tenants WHERE 
tenant_name = $1 AND service_id = $2 AND deleted_at IS NULL;

-- name: TenantRegisterPermission :one
INSERT INTO permissions (name,description,statement,service_id,tenant_id)
SELECT $1,$2,$3,$4,t.id from tenants t where t.tenant_name=$5 and t.deleted_at is null
RETURNING permissions.id,permissions.statement,permissions.description,permissions.name,permissions.created_at,permissions.service_id, $5::string as tenant;

-- name: IsPermissionExistsInTenant :one
SELECT count_rows() FROM permissions p join tenants t on p.tenant_id=t.id WHERE p.name =$1 and p.service_id=$2 and t.tenant_name=$3 and p.deleted_at IS NULL and p.deleted_at IS NULL ;