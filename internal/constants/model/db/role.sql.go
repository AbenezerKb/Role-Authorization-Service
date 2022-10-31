// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: role.sql

package db

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const assignRole = `-- name: AssignRole :exec
insert into tenant_users_roles(tenant_id, user_id, role_id)
select tenants.id,users.id,$1
from tenants,users  where tenants.tenant_name=$2
and users.user_id=$3 and users.deleted_at IS NULL and tenants.deleted_at IS NULL
`

type AssignRoleParams struct {
	RoleID     uuid.UUID `json:"role_id"`
	TenantName string    `json:"tenant_name"`
	UserID     uuid.UUID `json:"user_id"`
}

func (q *Queries) AssignRole(ctx context.Context, arg AssignRoleParams) error {
	_, err := q.db.Exec(ctx, assignRole, arg.RoleID, arg.TenantName, arg.UserID)
	return err
}

const createRole = `-- name: CreateRole :one
with _tenant as(
    select tenants.id as tenant_id from tenants where tenants.tenant_name=$1 AND tenants.service_id=$2 AND tenants.deleted_at IS NULL
), _role as(
    insert into roles(name, tenant_id) select $3,tenant_id from _tenant
        returning  id as role_id,name,created_at,status
),_rp as(
    insert into role_permissions (role_id,permission_id)
        select role_id,permissions.id from _role, permissions where permissions.id =ANY($4::uuid[]) and permissions.deleted_at IS NULL ON CONFLICT DO NOTHING returning id
)select _role.role_id, _role.name, _role.created_at, _role.status from _role
`

type CreateRoleParams struct {
	TenantName string      `json:"tenant_name"`
	ServiceID  uuid.UUID   `json:"service_id"`
	Name       string      `json:"name"`
	Column4    []uuid.UUID `json:"column_4"`
}

type CreateRoleRow struct {
	RoleID    uuid.UUID `json:"role_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	Status    Status    `json:"status"`
}

func (q *Queries) CreateRole(ctx context.Context, arg CreateRoleParams) (CreateRoleRow, error) {
	row := q.db.QueryRow(ctx, createRole,
		arg.TenantName,
		arg.ServiceID,
		arg.Name,
		arg.Column4,
	)
	var i CreateRoleRow
	err := row.Scan(
		&i.RoleID,
		&i.Name,
		&i.CreatedAt,
		&i.Status,
	)
	return i, err
}

const deleteRole = `-- name: DeleteRole :one
Update roles set deleted_at=now() where roles.id=$1 AND deleted_at IS NULL returning name,id,created_at,updated_at
`

type DeleteRoleRow struct {
	Name      string    `json:"name"`
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (q *Queries) DeleteRole(ctx context.Context, id uuid.UUID) (DeleteRoleRow, error) {
	row := q.db.QueryRow(ctx, deleteRole, id)
	var i DeleteRoleRow
	err := row.Scan(
		&i.Name,
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getRoleByNameAndTenantName = `-- name: GetRoleByNameAndTenantName :one
SELECT roles.id FROM roles join tenants on roles.tenant_id =tenants.id WHERE 
roles.name = $1 AND tenants.tenant_name = $2 and roles.deleted_at IS NULL and tenants.deleted_at IS NULL
`

type GetRoleByNameAndTenantNameParams struct {
	Name       string `json:"name"`
	TenantName string `json:"tenant_name"`
}

func (q *Queries) GetRoleByNameAndTenantName(ctx context.Context, arg GetRoleByNameAndTenantNameParams) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, getRoleByNameAndTenantName, arg.Name, arg.TenantName)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const isRoleAssigned = `-- name: IsRoleAssigned :one
SELECT count_rows() FROM tenant_users_roles 
WHERE tenant_users_roles.tenant_id in (
    SELECT tenants.id FROM 
    tenants where tenants.tenant_name = $1 and tenants.deleted_at IS NULL
)
and tenant_users_roles.user_id in (
    SELECT users.id from users 
    where users.user_id = $2 and users.deleted_at IS NULL
) and tenant_users_roles.role_id = $3
`

type IsRoleAssignedParams struct {
	TenantName string    `json:"tenant_name"`
	UserID     uuid.UUID `json:"user_id"`
	RoleID     uuid.UUID `json:"role_id"`
}

func (q *Queries) IsRoleAssigned(ctx context.Context, arg IsRoleAssignedParams) (interface{}, error) {
	row := q.db.QueryRow(ctx, isRoleAssigned, arg.TenantName, arg.UserID, arg.RoleID)
	var count_rows interface{}
	err := row.Scan(&count_rows)
	return count_rows, err
}

const removePermissionsFromRole = `-- name: RemovePermissionsFromRole :exec
DELETE FROM role_permissions WHERE role_id=$1 AND NOT permission_id=any($2::uuid[])
`

type RemovePermissionsFromRoleParams struct {
	RoleID  uuid.UUID   `json:"role_id"`
	Column2 []uuid.UUID `json:"column_2"`
}

func (q *Queries) RemovePermissionsFromRole(ctx context.Context, arg RemovePermissionsFromRoleParams) error {
	_, err := q.db.Exec(ctx, removePermissionsFromRole, arg.RoleID, arg.Column2)
	return err
}

const revokeUserRole = `-- name: RevokeUserRole :exec
UPDATE tenant_users_roles 
SET deleted_at= now() WHERE tenant_users_roles.tenant_id = (
    SELECT tenants.id FROM 
    tenants where tenants.tenant_name = $1 and tenants.deleted_at IS NULL
)
and tenant_users_roles.user_id = (
    SELECT users.id from users 
    where users.user_id = $2 and users.deleted_at IS NULL
) and tenant_users_roles.role_id = $3
`

type RevokeUserRoleParams struct {
	TenantName string    `json:"tenant_name"`
	UserID     uuid.UUID `json:"user_id"`
	RoleID     uuid.UUID `json:"role_id"`
}

func (q *Queries) RevokeUserRole(ctx context.Context, arg RevokeUserRoleParams) error {
	_, err := q.db.Exec(ctx, revokeUserRole, arg.TenantName, arg.UserID, arg.RoleID)
	return err
}

const updateRole = `-- name: UpdateRole :exec
INSERT INTO role_permissions (role_id,permission_id)
SELECT $1,permissions.id FROM permissions WHERE permissions.id =ANY($2::uuid[]) AND permissions.deleted_at IS NULL ON conflict do nothing
`

type UpdateRoleParams struct {
	RoleID  uuid.UUID   `json:"role_id"`
	Column2 []uuid.UUID `json:"column_2"`
}

func (q *Queries) UpdateRole(ctx context.Context, arg UpdateRoleParams) error {
	_, err := q.db.Exec(ctx, updateRole, arg.RoleID, arg.Column2)
	return err
}
