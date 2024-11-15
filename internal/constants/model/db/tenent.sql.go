// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: tenent.sql

package db

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgtype"
)

const checkIfPermissionExistsInTenant = `-- name: CheckIfPermissionExistsInTenant :one
with _tenant as (
    select tenants.domain_id,tenants.id from tenants where tenant_name =$1 and tenants.service_id=$2 and deleted_at IS NULL
),_sum as(
SELECT count_rows() as count FROM _tenant,permissions p WHERE p.name =$3 and p.service_id=$2 and p.tenant_id=_tenant.id and p.deleted_at IS NULL
union all
SELECT  count_rows() as count from _tenant,permission_domains pd join permissions p2 on pd.permission_id = p2.id where p2.name=$3 and pd.domain_id=_tenant.domain_id
)
select sum(count) from _sum
`

type CheckIfPermissionExistsInTenantParams struct {
	TenantName string    `json:"tenant_name"`
	ServiceID  uuid.UUID `json:"service_id"`
	Name       string    `json:"name"`
}

func (q *Queries) CheckIfPermissionExistsInTenant(ctx context.Context, arg CheckIfPermissionExistsInTenantParams) (int64, error) {
	row := q.db.QueryRow(ctx, checkIfPermissionExistsInTenant, arg.TenantName, arg.ServiceID, arg.Name)
	var sum int64
	err := row.Scan(&sum)
	return sum, err
}

const createTenent = `-- name: CreateTenent :exec
with new_tenant as (
    insert into tenants (
                         domain_id,
                         tenant_name,
                         service_id
        ) values ($1, $2, $3) returning id as tenant_id ),
new_role as (
insert into roles (name, tenant_id)
select 'admin', new_tenant.tenant_id from new_tenant returning id as role_id)
insert into
    role_permissions( role_id, permission_id)
select role_id, id as permission_id
       from  new_role,permissions
       where name = 'manage-all'
returning id
`

type CreateTenentParams struct {
	DomainID   uuid.UUID `json:"domain_id"`
	TenantName string    `json:"tenant_name"`
	ServiceID  uuid.UUID `json:"service_id"`
}

func (q *Queries) CreateTenent(ctx context.Context, arg CreateTenentParams) error {
	_, err := q.db.Exec(ctx, createTenent, arg.DomainID, arg.TenantName, arg.ServiceID)
	return err
}

const getTenentWithNameAndServiceId = `-- name: GetTenentWithNameAndServiceId :one
SELECT id, status, tenant_name, service_id, deleted_at, created_at, updated_at, domain_id, inherit FROM tenants WHERE 
tenant_name = $1 AND service_id = $2 AND deleted_at IS NULL
`

type GetTenentWithNameAndServiceIdParams struct {
	TenantName string    `json:"tenant_name"`
	ServiceID  uuid.UUID `json:"service_id"`
}

func (q *Queries) GetTenentWithNameAndServiceId(ctx context.Context, arg GetTenentWithNameAndServiceIdParams) (Tenant, error) {
	row := q.db.QueryRow(ctx, getTenentWithNameAndServiceId, arg.TenantName, arg.ServiceID)
	var i Tenant
	err := row.Scan(
		&i.ID,
		&i.Status,
		&i.TenantName,
		&i.ServiceID,
		&i.DeletedAt,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DomainID,
		&i.Inherit,
	)
	return i, err
}

const tenantRegisterPermission = `-- name: TenantRegisterPermission :one
with _permission as(
    INSERT INTO permissions (name,description,statement,service_id,tenant_id,delete_or_update)
    SELECT $1,$2,$3,$4,t.id,true from tenants t where t.tenant_name=$5 and t.deleted_at is null
    RETURNING permissions.id,permissions.statement,permissions.description,permissions.name,permissions.created_at,permissions.service_id, $5::string as tenant
), _ph as(
    insert into permissions_hierarchy(parent, child) select _permission.id,p.id from _permission,permissions p where p.name=ANY($6::string[])  and p.service_id=$4 and p.deleted_at IS NULl ON conflict  do nothing returning id
)
select id, statement, description, name, created_at, service_id, tenant from _permission
`

type TenantRegisterPermissionParams struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Statement   pgtype.JSON `json:"statement"`
	ServiceID   uuid.UUID   `json:"service_id"`
	TenantName  string      `json:"tenant_name"`
	Column6     []string    `json:"column_6"`
}

type TenantRegisterPermissionRow struct {
	ID          uuid.UUID   `json:"id"`
	Statement   pgtype.JSON `json:"statement"`
	Description string      `json:"description"`
	Name        string      `json:"name"`
	CreatedAt   time.Time   `json:"created_at"`
	ServiceID   uuid.UUID   `json:"service_id"`
	Tenant      string      `json:"tenant"`
}

func (q *Queries) TenantRegisterPermission(ctx context.Context, arg TenantRegisterPermissionParams) (TenantRegisterPermissionRow, error) {
	row := q.db.QueryRow(ctx, tenantRegisterPermission,
		arg.Name,
		arg.Description,
		arg.Statement,
		arg.ServiceID,
		arg.TenantName,
		arg.Column6,
	)
	var i TenantRegisterPermissionRow
	err := row.Scan(
		&i.ID,
		&i.Statement,
		&i.Description,
		&i.Name,
		&i.CreatedAt,
		&i.ServiceID,
		&i.Tenant,
	)
	return i, err
}

const updateTenantStatus = `-- name: UpdateTenantStatus :one
update tenants set status=$1 where tenant_name=$2 and service_id=$3 returning id
`

type UpdateTenantStatusParams struct {
	Status     Status    `json:"status"`
	TenantName string    `json:"tenant_name"`
	ServiceID  uuid.UUID `json:"service_id"`
}

func (q *Queries) UpdateTenantStatus(ctx context.Context, arg UpdateTenantStatusParams) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, updateTenantStatus, arg.Status, arg.TenantName, arg.ServiceID)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}
