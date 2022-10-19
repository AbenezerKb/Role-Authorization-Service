-- name: GetOpaData :one
select
    json_object_agg(id, result) as services
from
    (
        select
            s.id::text,
            s.status,
            json_object_agg(tenant.tenant_name, tenant) as tenants
        from
            services s
                join
            (
                select
                    tenants.id,
                    tenants.status,
                    tenants.tenant_name,
                    tenants.service_id,
                    json_object_agg(tu.user_id::text, tu) as users
                from
                    tenants
                        join(
                        select
                            u.status,
                            u.user_id,
                            tur.tenant_id,
                            tur.status as user_role_status,
                            json_agg(role) as role
                        from
                            users as u
                                join tenant_users_roles tur on
                                    tur.user_id = u.id
                                join(
                                select
                                    r.id,
                                    r.status,
                                    r.name,
                                    json_agg(rp)as permissions
                                from
                                    roles r
                                        join (
                                        select
                                            p.id,
                                            role_permissions.role_id ,
                                            p.name,
                                            p.description,
                                            p.statment,
                                            p.status,
                                            coalesce(json_agg(ph) filter (where ph.id is not null), '[]') as child
                                        from
                                            role_permissions
                                                join permissions p on
                                                    role_permissions.permission_id = p.id
                                                left join (
                                                select
                                                    p2.name,
                                                    p2.id,
                                                    ph.parent,
                                                    p2.statment,
                                                    p2.status
                                                from
                                                    permissions_hierarchy ph
                                                        join permissions p2 on
                                                            ph.child = p2.id
                                            )as ph on
                                                    ph.parent = p.id
                                        group by
                                            p.id,
                                            role_permissions.role_id ,
                                            p.name,
                                            p.description,
                                            p.statment,
                                            p.status
                                    )as rp on
                                            r.id = rp.role_id
                                group by
                                    r.id,
                                    r.status,
                                    r.name
                            ) as role on
                                    role.id = tur.role_id group by    u.status,
                                                                      u.user_id,
                                                                      tur.tenant_id,
                                                                      tur.status
                    ) as tu on
                            tu.tenant_id = tenants.id
                group by
                    tenants.id,
                    tenants.status,
                    tenants.tenant_name,
                    tenants.service_id
            ) as tenant
            on
                    s.id = tenant.service_id
        group by
            s.id,
            s.status
    ) result;