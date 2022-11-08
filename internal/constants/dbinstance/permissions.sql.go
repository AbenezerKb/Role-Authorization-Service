package dbinstance

import (
	"2f-authorization/internal/constants/model/dto"
	"context"

	"github.com/google/uuid"
)

const listPermissions = `-- name: ListPermissions :many
with _tenant as (
    select tenants.domain_id,tenants.id,tenants.inherit from tenants where tenant_name =$1 and tenants.service_id=$2 and deleted_at IS NULL
)
select p.name,p.status,p.description,p.statement,p.created_at,p.service_id,p.id,coalesce(json_agg(json_build_object('name',p2.name,'status',p2.status,'description',p2.description,'statement',p2.statement,'created_at',p2.created_at,'service_id',p2.service_id,'id',p2.id)) filter ( where p2.name is not null ),'[]') AS inherited_permissions from _tenant,permissions p  left join permission_domains pd on p.id = pd.permission_id left join permissions_hierarchy ph on p.id = ph.parent join permissions p2 on p2.id = ph.child where p2.deleted_at is null and pd.domain_id = _tenant.domain_id and _tenant.inherit = true group by p.name,p.status,p.description,p.statement,p.created_at,p.service_id,p.id
UNION
select p.name,p.status,p.description,p.statement,p.created_at,p.service_id,p.id,coalesce(json_agg(json_build_object('name',p2.name,'status',p2.status,'description',p2.description,'statement',p2.statement,'created_at',p2.created_at,'service_id',p2.service_id,'id',p2.id))filter ( where p2.name is not null ),'[]') AS inherited_permissions  from permissions p left join permissions_hierarchy ph on p.id = ph.parent left join permissions p2 on p2.id = ph.child,_tenant where  p2.deleted_at is null and p.tenant_id =_tenant.id and p.deleted_at IS NULL group by p.name,p.status,p.description,p.statement,p.created_at,p.service_id,p.id
`

type ListPermissionsParams struct {
	TenantName string    `json:"tenant_name"`
	ServiceID  uuid.UUID `json:"service_id"`
}

func (q *DBInstance) ListPermissions(ctx context.Context, arg ListPermissionsParams) ([]dto.Permission, error) {
	rows, err := q.Pool.Query(ctx, listPermissions, arg.TenantName, arg.ServiceID)
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
			&i.Description,
			&i.Statement,
			&i.CreatedAt,
			&i.ServiceID,
			&i.ID,
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
