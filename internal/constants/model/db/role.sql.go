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

const createRole = `-- name: CreateRole :one
with _tenant as(
    select tenants.id as tenant_id from tenants where tenants.tenant_name=$1 AND tenants.service_id=$2
), _role as(
    insert into roles(name, tenant_id) select $3,tenant_id from _tenant
    returning  id as role_id,name,created_at,status
),_rp as(
insert into role_permissions (role_id,permission_id)
    select role_id,permissions.id from _role, permissions where permissions.id =ANY($4::uuid[])returning id
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
	Status    bool      `json:"status"`
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

const getRoleByNameAndTenantName = `-- name: GetRoleByNameAndTenantName :one
SELECT roles.id FROM roles join tenants on roles.tenant_id =tenants.id WHERE 
roles.name = $1 AND tenants.tenant_name = $2
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
