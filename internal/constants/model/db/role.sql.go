// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: role.sql

package db

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const assignRole = `-- name: AssignRole :exec
WITH new_row AS (
    INSERT INTO users (user_id,  service_id) select $1,$2 where not exists(select id from users where user_id=$1 and service_id=$2 and deleted_at is null)returning id
),_user as(
    SELECT id FROM new_row
    UNION
    SELECT id FROM users WHERE user_id=$1 and service_id=$2 and deleted_at is null
)
insert into tenant_users_roles(tenant_id, user_id, role_id)
select tenants.id,_user.id,roles.id
from roles,tenants,_user  where tenants.tenant_name=$5 and tenants.deleted_at IS NULL and (roles.id=$3 or roles.name=$4) and roles.tenant_id=tenants.id and roles.deleted_at is null
`

type AssignRoleParams struct {
	UserID     uuid.UUID `json:"user_id"`
	ServiceID  uuid.UUID `json:"service_id"`
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	TenantName string    `json:"tenant_name"`
}

func (q *Queries) AssignRole(ctx context.Context, arg AssignRoleParams) error {
	_, err := q.db.Exec(ctx, assignRole,
		arg.UserID,
		arg.ServiceID,
		arg.ID,
		arg.Name,
		arg.TenantName,
	)
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
Update roles set deleted_at=now() where roles.id=$1 and r.name !='admin' AND deleted_at IS NULL returning name,id,created_at,updated_at
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

const getRoleById = `-- name: GetRoleById :one
select r.name,r.id,r.status,r.created_at,r.updated_at, (select string_to_array(string_agg(p.name,','),',')::string[] from role_permissions join permissions p on role_permissions.permission_id = p.id where role_id=r.id and p.deleted_at is null) as permission  from roles r join tenants t on t.id = r.tenant_id where t.service_id=$1 and t.deleted_at is null and r.id=$2 and r.deleted_at is null
`

type GetRoleByIdParams struct {
	ServiceID uuid.UUID `json:"service_id"`
	ID        uuid.UUID `json:"id"`
}

type GetRoleByIdRow struct {
	Name       string    `json:"name"`
	ID         uuid.UUID `json:"id"`
	Status     Status    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Permission []string  `json:"permission"`
}

func (q *Queries) GetRoleById(ctx context.Context, arg GetRoleByIdParams) (GetRoleByIdRow, error) {
	row := q.db.QueryRow(ctx, getRoleById, arg.ServiceID, arg.ID)
	var i GetRoleByIdRow
	err := row.Scan(
		&i.Name,
		&i.ID,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Permission,
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
select count_rows() from tenant_users_roles tur
    join tenants t on t.id = tur.tenant_id
    join users u on u.id = tur.user_id
    join roles r on r.id = tur.role_id 
    where t.tenant_name=$1
and u.user_id=$2 and( r.id=$3 or r.name=$4)
 and tur.deleted_at is null and r.deleted_at is null
`

type IsRoleAssignedParams struct {
	TenantName string    `json:"tenant_name"`
	UserID     uuid.UUID `json:"user_id"`
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
}

func (q *Queries) IsRoleAssigned(ctx context.Context, arg IsRoleAssignedParams) (interface{}, error) {
	row := q.db.QueryRow(ctx, isRoleAssigned,
		arg.TenantName,
		arg.UserID,
		arg.ID,
		arg.Name,
	)
	var count_rows interface{}
	err := row.Scan(&count_rows)
	return count_rows, err
}

const listRoles = `-- name: ListRoles :many
select r.name,r.created_at,r.id,r.status from roles r join tenants t on r.tenant_id=t.id where t.tenant_name=$1 AND t.service_id=$2 AND t.deleted_at IS NULL AND r.deleted_at IS NULL
`

type ListRolesParams struct {
	TenantName string    `json:"tenant_name"`
	ServiceID  uuid.UUID `json:"service_id"`
}

type ListRolesRow struct {
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	ID        uuid.UUID `json:"id"`
	Status    Status    `json:"status"`
}

func (q *Queries) ListRoles(ctx context.Context, arg ListRolesParams) ([]ListRolesRow, error) {
	rows, err := q.db.Query(ctx, listRoles, arg.TenantName, arg.ServiceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListRolesRow
	for rows.Next() {
		var i ListRolesRow
		if err := rows.Scan(
			&i.Name,
			&i.CreatedAt,
			&i.ID,
			&i.Status,
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

const revokeAdminRole = `-- name: RevokeAdminRole :exec
UPDATE tenant_users_roles tur
SET status = 'INACTIVE'
FROM roles r, tenants t
WHERE t.tenant_name = $1
  AND r.name = 'admin'
  AND r.id = tur.role_id
  AND tur.tenant_id = t.id
`

func (q *Queries) RevokeAdminRole(ctx context.Context, tenantName string) error {
	_, err := q.db.Exec(ctx, revokeAdminRole, tenantName)
	return err
}

const revokeUserRole = `-- name: RevokeUserRole :exec
UPDATE tenant_users_roles tur
SET deleted_at= now() FROM roles r, tenants t,users u
WHERE t.tenant_name =  $1
  AND tur.role_id =  $2
  AND u.user_id=$3
  AND r.id = tur.role_id
  AND u.id=tur.user_id
  AND tur.tenant_id = t.id and tur.deleted_at is null and r.deleted_at is null and u.deleted_at is null
`

type RevokeUserRoleParams struct {
	TenantName string    `json:"tenant_name"`
	RoleID     uuid.UUID `json:"role_id"`
	UserID     uuid.UUID `json:"user_id"`
}

func (q *Queries) RevokeUserRole(ctx context.Context, arg RevokeUserRoleParams) error {
	_, err := q.db.Exec(ctx, revokeUserRole, arg.TenantName, arg.RoleID, arg.UserID)
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

const updateRoleStatus = `-- name: UpdateRoleStatus :one
with _tenants as(
    select id from tenants t where t.tenant_name=$1 and t.service_id=$2 and t.deleted_at IS NULL
)
update roles r set status =$3 from _tenants where r.id=$4 and r.name !='admin' and r.deleted_at IS NULL and r.tenant_id=_tenants.id returning r.id
`

type UpdateRoleStatusParams struct {
	TenantName string    `json:"tenant_name"`
	ServiceID  uuid.UUID `json:"service_id"`
	Status     Status    `json:"status"`
	ID         uuid.UUID `json:"id"`
}

func (q *Queries) UpdateRoleStatus(ctx context.Context, arg UpdateRoleStatusParams) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, updateRoleStatus,
		arg.TenantName,
		arg.ServiceID,
		arg.Status,
		arg.ID,
	)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}
