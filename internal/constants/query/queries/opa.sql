-- name: GetOpaData :one
select
	coalesce(json_object_agg(id, result) filter ( where id is not null), '{}') as services
from
	(
	select
		s.id::text,
		s.status,
		coalesce(json_object_agg(tenant.tenant_name, tenant)filter ( where tenant.tenant_name is not null), '{}') as tenants
	from
		services s
	left join
            (
		select
			tenants.id,
			tenants.status,
			tenants.tenant_name,
			tenants.service_id,
			coalesce(json_object_agg(tu.user_id::text, tu)filter ( where tu.user_id is not null), '{}') as users
		from
			tenants
		left join(
			select
				u.status,
				u.user_id,
				tur.tenant_id,
				tur.status as user_role_status,
				coalesce(json_agg(role)filter ( where role.id is not null), '[]') as role
			from
				users as u
			left join tenant_users_roles tur on
				tur.user_id = u.id
			left join(
				select
					r.id,
					r.status,
					r.name,
					coalesce(json_agg(rp)filter ( where rp.id is not null), '[]') as permissions
				from
					roles r
				left join (
					select
						p.id,
						role_permissions.role_id ,
						p.name,
						p.description,
						p.statement,
						p.status,
						coalesce(json_agg(ph) filter (where ph.id is not null), '[]') as child
					from
						role_permissions
					left join permissions p on
						role_permissions.permission_id = p.id
					left join (
						select
							p2.name,
							p2.id,
							ph.parent,
							p2.statement,
							p2.status
						from
							permissions_hierarchy ph
						left join permissions p2 on
							ph.child = p2.id
						where
							p2.deleted_at is null
                                            )as ph on
						ph.parent = p.id
					where
						p.deleted_at is null
					group by
						p.id,
						role_permissions.role_id ,
						p.name,
						p.description,
						p.statement,
						p.status
                                    )as rp on
					r.id = rp.role_id
				where
					r.deleted_at is null
				group by
					r.id,
					r.status,
					r.name
                            ) as role on
				role.id = tur.role_id
			where
				u.deleted_at is null
			group by
				u.status,
				u.user_id,
				tur.tenant_id,
				tur.status
                    ) as tu on
			tu.tenant_id = tenants.id
		where
			tenants.deleted_at is null
		group by
			tenants.id,
			tenants.status,
			tenants.tenant_name,
			tenants.service_id
            ) as tenant
            on
		s.id = tenant.service_id
	where
		s.deleted_at is null
	group by
		s.id,
		s.status
    ) result;
