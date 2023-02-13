-- name: CreateTenent :exec 
with new_tenant as (
    insert into tenants (
                         domain_id,
                         tenant_name,
                         service_id
        ) values ($1, $2, $3) returning id as tenant_id ),
new_role as (
insert into roles (name, tenant_id)
select 'admin', new_tenant.tenant_id from new_tenant returning id as role_id)
insert into
    role_permissions( role_id, permission_id)
select role_id, id as permission_id
       from  new_role,permissions
       where name = 'manage-all'
returning id;

-- name: GetTenentWithNameAndServiceId :one 
SELECT * FROM tenants WHERE 
tenant_name = $1 AND service_id = $2 AND deleted_at IS NULL;

-- name: TenantRegisterPermission :one
with _permission as(
    INSERT INTO permissions (name,description,statement,service_id,tenant_id,delete_or_update)
    SELECT $1,$2,$3,$4,t.id,true from tenants t where t.tenant_name=$5 and t.deleted_at is null
    RETURNING permissions.id,permissions.statement,permissions.description,permissions.name,permissions.created_at,permissions.service_id, $5::string as tenant
), _ph as(
    insert into permissions_hierarchy(parent, child) select _permission.id,p.id from _permission,permissions p where p.name=ANY($6::string[])  and p.service_id=$4 and p.deleted_at IS NULl ON conflict  do nothing returning id
)
select * from _permission;


-- name: CheckIfPermissionExistsInTenant :one
with _tenant as (
    select tenants.domain_id,tenants.id from tenants where tenant_name =$1 and tenants.service_id=$2 and deleted_at IS NULL
),_sum as(
SELECT count_rows() as count FROM _tenant,permissions p WHERE p.name =$3 and p.service_id=$2 and p.tenant_id=_tenant.id and p.deleted_at IS NULL
union all
SELECT  count_rows() as count from _tenant,permission_domains pd join permissions p2 on pd.permission_id = p2.id where p2.name=$3 and pd.domain_id=_tenant.domain_id
)
select sum(count) from _sum;

-- name: UpdateTenantStatus :one
update tenants set status=$1 where tenant_name=$2 and service_id=$3 returning id;