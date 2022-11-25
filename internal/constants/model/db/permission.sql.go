// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: permission.sql

package db

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgtype"
)

const assignDomain = `-- name: AssignDomain :exec
with _domain as(
    select domains.id as domain_id from domains where domains.id=$1 and deleted_at IS NULL
) INSERT INTO permission_domains(domain_id,permission_id)
SELECT  domain_id, $2 from _domain
WHERE NOT exists(select permission_id from permission_domains where permission_id=$2 and domain_id=_domain.domain_id and deleted_at IS NULL)
`

type AssignDomainParams struct {
	ID           uuid.UUID `json:"id"`
	PermissionID uuid.UUID `json:"permission_id"`
}

func (q *Queries) AssignDomain(ctx context.Context, arg AssignDomainParams) error {
	_, err := q.db.Exec(ctx, assignDomain, arg.ID, arg.PermissionID)
	return err
}

const canBeDeleted = `-- name: CanBeDeleted :one
select p.delete_or_update from permissions p where p.id=$1 and p.service_id=$2
`

type CanBeDeletedParams struct {
	ID        uuid.UUID `json:"id"`
	ServiceID uuid.UUID `json:"service_id"`
}

func (q *Queries) CanBeDeleted(ctx context.Context, arg CanBeDeletedParams) (bool, error) {
	row := q.db.QueryRow(ctx, canBeDeleted, arg.ID, arg.ServiceID)
	var delete_or_update bool
	err := row.Scan(&delete_or_update)
	return delete_or_update, err
}

const createOrGetPermission = `-- name: CreateOrGetPermission :one
WITH new_row AS (
    INSERT INTO permissions (name,description,statement,service_id)
        SELECT $1,$2,$3,$4
        WHERE NOT EXISTS (SELECT id FROM permissions WHERE name =$1 and service_id=$4 and deleted_at IS NULL)
        RETURNING id
),_permission as(
    SELECT id FROM new_row
    UNION
    SELECT id FROM permissions WHERE name = $1 and service_id=$4 and deleted_at IS NULL
)
,pd as (insert into permission_domains (domain_id,permission_id)
select domains.id,_permission.id from domains,_permission where domains.id =ANY($5::uuid[]) ON CONFLICT DO NOTHING returning permission_domains.id
)
select id from _permission
`

type CreateOrGetPermissionParams struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Statement   pgtype.JSON `json:"statement"`
	ServiceID   uuid.UUID   `json:"service_id"`
	Column5     []uuid.UUID `json:"column_5"`
}

func (q *Queries) CreateOrGetPermission(ctx context.Context, arg CreateOrGetPermissionParams) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, createOrGetPermission,
		arg.Name,
		arg.Description,
		arg.Statement,
		arg.ServiceID,
		arg.Column5,
	)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const createPermissionDependency = `-- name: CreatePermissionDependency :exec
with _parent as(
    select id as parant_id from permissions where permissions.name=$1 and permissions.service_id=$2 and permissions.deleted_at IS NULL
)
insert into permissions_hierarchy(parent, child) select _parent.parant_id,p.id from _parent,permissions p where p.name=ANY($3::string[]) and p.service_id=$2  and p.deleted_at IS NULl ON conflict  do nothing
`

type CreatePermissionDependencyParams struct {
	Name      string    `json:"name"`
	ServiceID uuid.UUID `json:"service_id"`
	Column3   []string  `json:"column_3"`
}

func (q *Queries) CreatePermissionDependency(ctx context.Context, arg CreatePermissionDependencyParams) error {
	_, err := q.db.Exec(ctx, createPermissionDependency, arg.Name, arg.ServiceID, arg.Column3)
	return err
}

const deletePermissions = `-- name: DeletePermissions :one
UPDATE permissions p set deleted_at = now() from tenants t WHERE  t.tenant_name=$1
and p.id = $2 and p.tenant_id=t.id AND p.service_id = $3 AND t.service_id = $3
AND p.deleted_at IS NULL AND p.delete_or_update RETURNING p.id
`

type DeletePermissionsParams struct {
	TenantName string    `json:"tenant_name"`
	ID         uuid.UUID `json:"id"`
	ServiceID  uuid.UUID `json:"service_id"`
}

func (q *Queries) DeletePermissions(ctx context.Context, arg DeletePermissionsParams) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, deletePermissions, arg.TenantName, arg.ID, arg.ServiceID)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
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

type GetPermissionDetailsRow struct {
	Name                 string      `json:"name"`
	Status               Status      `json:"status"`
	Description          string      `json:"description"`
	Statement            pgtype.JSON `json:"statement"`
	ID                   uuid.UUID   `json:"id"`
	InheritedPermissions interface{} `json:"inherited_permissions"`
}

func (q *Queries) GetPermissionDetails(ctx context.Context, arg GetPermissionDetailsParams) (GetPermissionDetailsRow, error) {
	row := q.db.QueryRow(ctx, getPermissionDetails, arg.TenantName, arg.ID, arg.ServiceID)
	var i GetPermissionDetailsRow
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

type ListPermissionsRow struct {
	Name                 string      `json:"name"`
	Status               Status      `json:"status"`
	Description          string      `json:"description"`
	Statement            pgtype.JSON `json:"statement"`
	CreatedAt            time.Time   `json:"created_at"`
	ServiceID            uuid.UUID   `json:"service_id"`
	ID                   uuid.UUID   `json:"id"`
	InheritedPermissions interface{} `json:"inherited_permissions"`
}

func (q *Queries) ListPermissions(ctx context.Context, arg ListPermissionsParams) ([]ListPermissionsRow, error) {
	rows, err := q.db.Query(ctx, listPermissions, arg.TenantName, arg.ServiceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListPermissionsRow
	for rows.Next() {
		var i ListPermissionsRow
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
