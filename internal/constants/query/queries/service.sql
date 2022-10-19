-- name: GetServiceByName :one
SELECT * FROM services WHERE name = $1;

-- name: GetServiceById :one
SELECT * FROM services WHERE id = $1 AND deleted_at IS NULL;

-- name: CreateService :one
with _service as (
    insert
        into
            services (name,
                      password)
            values ($1,
                    $2)
            returning password,
                id as service_id,
                name as service,
                status as service_status
),
    _domain as(
        insert
            into
                domains (name,
                         service_id)
                select
                    'administrator',
                    service_id
                from
                    _service
                returning id as domain
    ),
     _tenant as (
         insert
             into
                 tenants (tenant_name,
                          service_id,domain_id)
                 select
                     'administrator',
                     service_id,
                     domain
                 from
                     _service,_domain
                 returning tenant_name as tenant,
                     id as tenant_id
     ),
     _role as (
         insert
             into
                 roles(name,
                       tenant_id)
                 select
                     'service-admin',
                     tenant_id
                 from
                     _tenant
                 returning id as role_id
     ),
     _permission as(
         insert
             into
                 permissions(name,
                             description,
                             statment,service_id)
                 select 'manage-all',
                        'super admin can perform any action on any domain',
                        json_build_object('action', '*', 'resource', '*', 'effect', 'allow'),service_id from _service
                 returning id as permission_id
     ),
     _user as(
         insert
             into
                 users(user_id,service_id)
                 select $3,service_id from _service
                 returning id as user_id
     ),
     _tenant_user_role as
         (
             insert
                 into
                     tenant_users_roles(tenant_id,
                                        user_id,
                                        role_id)
                     select
                         tenant_id ,
                         user_id,
                         role_id
                     from
                         _tenant,
                         _user,
                         _role
                     returning role_id),
     _role_permission as
         (
             insert
                 into
                     role_permissions (role_id,
                                       permission_id)
                     select
                         role_id,
                         permission_id
                     from
                         _tenant_user_role,
                         _permission
                     returning role_id
         )
select
    service_id,
    password,
    service,
    service_status,
    tenant
from
    _service,
    _tenant;

-- name: DeleteService :exec
DELETE FROM services WHERE id = $1;

-- name: SoftDeleteService :one
UPDATE services set deleted_at = now() WHERE id = $1 AND deleted_at IS NULL returning *;

