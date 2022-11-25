package dbinstance

import (
	"2f-authorization/internal/constants/model/dto"
	"context"

	"github.com/google/uuid"
)

const listPermissions = `-- name: ListPermissions :many
with _tenant as (
    select tenants.domain_id,tenants.id,tenants.inherit from tenants where tenant_name =$1 and tenants.service_id=$2 and tenants.deleted_at IS NULL
)
select p.name,p.status,p.description,p.statement,p.created_at,p.service_id,p.id,coalesce(json_agg(json_build_object('name',p2.name,'status',p2.status,'description',p2.description,'statement',p2.statement,'created_at',p2.created_at,'service_id',p2.service_id,'id',p2.id)) filter ( where  p2.deleted_at is null and p2.status='ACTIVE'  ),'[]') AS inherited_permissions from _tenant,permissions p  left join permission_domains pd on p.id = pd.permission_id left join permissions_hierarchy ph on p.id = ph.parent left join permissions p2 on p2.id = ph.child where pd.domain_id = _tenant.domain_id and _tenant.inherit = true group by p.name,p.status,p.description,p.statement,p.created_at,p.service_id,p.id
UNION
select p.name,p.status,p.description,p.statement,p.created_at,p.service_id,p.id,coalesce(json_agg(json_build_object('name',p2.name,'status',p2.status,'description',p2.description,'statement',p2.statement,'created_at',p2.created_at,'service_id',p2.service_id,'id',p2.id))filter ( where  p2.deleted_at is null and p2.status='ACTIVE'  ),'[]') AS inherited_permissions  from permissions p left join permissions_hierarchy ph on p.id = ph.parent left join permissions p2 on p2.id = ph.child,_tenant where p.tenant_id =_tenant.id and p.deleted_at IS NULL group by p.name,p.status,p.description,p.statement,p.created_at,p.service_id,p.id
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

const getPermissionDetails = `-- name: GetPermissionDetails :one
WITH _tenant AS(
    SELECT
        domain_id,
        id AS tenant_id
    FROM
        tenants
    WHERE
            tenant_name = $1
)
SELECT
    p.name,
    p.status,
    p.description,
    p.statement,
    p.id,
   COALESCE(json_agg(json_build_object('name',
                                        p2.name,
                                        'description',
                                        p2.description,
                                        'statement',
                                        p2.statement,
                                        'id',
                                        p2.id)) FILTER (
                 WHERE
                     p2.deleted_at IS NULL
                     AND p2.status = 'ACTIVE' ),
             '[]') AS inherited_permissions
FROM
    _tenant,
    permissions p
        LEFT JOIN permissions_hierarchy ph ON
            p.id = ph.parent
        LEFT JOIN permissions p2 ON
            p2.id = ph.child
        LEFT JOIN permission_domains pd ON
            p.id = pd.permission_id
WHERE
    (p.tenant_id = _tenant.tenant_id
        OR pd.domain_id = _tenant.domain_id)
  AND p.id = $2
  AND p.service_id = $3
  AND p.deleted_at IS NULL
GROUP BY
    p.name,
    p.status,
    p.description,
    p.statement,
    p.id
`

type GetPermissionDetailsParams struct {
	TenantName string    `json:"tenant_name"`
	ID         uuid.UUID `json:"id"`
	ServiceID  uuid.UUID `json:"service_id"`
}

func (q *DBInstance) GetPermissionDetails(ctx context.Context, arg GetPermissionDetailsParams) (dto.Permission, error) {
	row := q.Pool.QueryRow(ctx, getPermissionDetails, arg.TenantName, arg.ID, arg.ServiceID)
	var i dto.Permission
	err := row.Scan(
		&i.Name,
		&i.Status,
		&i.Description,
		&i.Statement,
		&i.ID,
		&i.InheritedPermissions,
	)
	return i, err
}
