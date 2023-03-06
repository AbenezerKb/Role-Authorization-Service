package dbinstance

import (
	"2f-authorization/internal/constants/model/dto"
	"context"

	"github.com/google/uuid"
)

const listPermissions = `-- name: ListPermissions :many
WITH _tenant AS (
SELECT
	tenants.domain_id,
	tenants.id,
	tenants.inherit
FROM
	tenants
WHERE
	tenant_name = $1
	AND tenants.service_id = $2
	AND tenants.deleted_at IS NULL )
SELECT
	p.name,
	p.status,
	p.description,
	p.statement,
	p.id,
	p.delete_or_update,
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
LEFT JOIN
    permission_domains pd 
        ON
	p.id = pd.permission_id
LEFT JOIN
    permissions_hierarchy ph 
        ON
	p.id = ph.parent
LEFT JOIN
    permissions p2 
        ON
	p2.id = ph.child
WHERE
	p.deleted_at IS NULL
	AND(
        p.tenant_id = _tenant.id
		OR pd.domain_id = _tenant.domain_id 
    )
GROUP BY
	p.name,
	p.status,
	p.description,
	p.statement,
	p.id,
	p.delete_or_update
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
	items := []dto.Permission{}
	for rows.Next() {
		var i dto.Permission
		if err := rows.Scan(
			&i.Name,
			&i.Status,
			&i.Description,
			&i.Statement,
			&i.ID,
			&i.DeleteOrUpdate,
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
    p.delete_or_update,
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
    p.id,
    p.delete_or_update
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
		&i.DeleteOrUpdate,
		&i.InheritedPermissions,
	)
	return i, err
}
