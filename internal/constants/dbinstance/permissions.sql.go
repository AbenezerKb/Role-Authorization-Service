package dbinstance

import (
	"2f-authorization/internal/constants/model/dto"
	"context"

	"github.com/google/uuid"
)

const listPermissions = `-- name: ListPermissions :many
with _tenant as (
    select tenants.domain_id,tenants.id from tenants where tenant_name =$1 and tenants.service_id=$2
)
select p.name,p.status,p.description,p.statement,p.created_at,p.service_id,p.id  from _tenant,permissions p  join permission_domains pd on p.id = pd.permission_id where pd.domain_id = _tenant.domain_id
UNION
select p.name,p.status,p.description,p.statement,p.created_at,p.service_id,p.id  from permissions p,_tenant where p.tenant_id =_tenant.id
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
			&i.ID,
			&i.ServiceID,
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
