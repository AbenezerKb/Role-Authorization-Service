package dbinstance

import (
	"2f-authorization/internal/constants/model/dto"
	"context"

	"github.com/google/uuid"
)

const listRoles = `-- name: ListRoles :many
select r.name,r.created_at,r.id,r.status from roles r join tenants t on r.tenant_id=t.id where t.tenant_name=$1 AND t.service_id=$2 AND t.deleted_at IS NULL AND r.deleted_at IS NULL
`

type ListRolesParams struct {
	TenantName string    `json:"tenant_name"`
	ServiceID  uuid.UUID `json:"service_id"`
}

func (db *DBInstance) ListRoles(ctx context.Context, arg ListRolesParams) ([]dto.Role, error) {
	rows, err := db.Pool.Query(ctx, listRoles, arg.TenantName, arg.ServiceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []dto.Role
	for rows.Next() {
		var i dto.Role
		if err := rows.Scan(&i.Name, &i.CreatedAt, &i.ID,&i.Status); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
