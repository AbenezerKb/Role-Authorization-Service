package dbinstance

import (
	"2f-authorization/internal/constants/model/dto"
	"context"

	"github.com/google/uuid"
)

const getUserPermissionWithInTenant = `-- name: GetUserPermissionWithInTenant :many
with _user as(
    select id from users u where u.user_id=$1 and u.deleted_at is null
),_tenant as(
    select id from tenants t where t.tenant_name=$2 and t.service_id=$3
)
select p.name,p.status,p.created_at,p.statement,p.id,p.description,coalesce(json_agg(json_build_object('name',p2.name,'status',p2.status,'description',p2.description,'statement',p2.statement,'created_at',p2.created_at,'id',p2.id))filter ( where p2.deleted_at is null and p2.status='ACTIVE' ),'[]') AS inherited_permissions  from   _user,_tenant,tenant_users_roles tur left join roles r on tur.role_id = r.id left join role_permissions rp on r.id = rp.role_id left join permissions p on p.id = rp.permission_id left join permissions_hierarchy ph on p.id = ph.parent left join permissions p2 on p2.id = ph.child where tur.deleted_at is null and tur.status='ACTIVE' and p.status='ACTIVE'  and
tur.tenant_id=_tenant.id and tur.status='ACTIVE' and r.status='ACTIVE' and tur.user_id=_user.id group by  p.name,p.status,p.created_at,p.statement,p.id
`

type GetUserPermissionWithInTenantParams struct {
	UserID     uuid.UUID `json:"user_id"`
	TenantName string    `json:"tenant_name"`
	ServiceID  uuid.UUID `json:"service_id"`
}

func (q *DBInstance) GetUserPermissionWithInTenant(ctx context.Context, arg GetUserPermissionWithInTenantParams) ([]dto.Permission, error) {
	rows, err := q.Pool.Query(ctx, getUserPermissionWithInTenant, arg.UserID, arg.TenantName, arg.ServiceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []dto.Permission
	for rows.Next() {
		var i dto.Permission
		if err := rows.Scan(
			&i.Name,
			&i.Status,
			&i.CreatedAt,
			&i.Statement,
			&i.ID,
			&i.Description,
			&i.InheritedPermissions,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
