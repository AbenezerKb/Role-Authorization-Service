// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: service.sql

package db

import (
	"context"

	"github.com/google/uuid"
)

const createService = `-- name: CreateService :one
with _service as (
    insert
        into
            services (name,
                      password)
            values ($1,
                    $2)
            returning password,
                id as service_id,
                name as service,
                status as service_status
),
    _domain as(
        insert
            into
                domains (name,
                         service_id)
                select
                    'administrator',
                    service_id
                from
                    _service
                returning id as domain
    ),
     _tenant as (
         insert
             into
                 tenants (tenant_name,
                          service_id,domain_id)
                 select
                     'administrator',
                     service_id,
                     domain
                 from
                     _service,_domain
                 returning tenant_name as tenant,
                     id as tenant_id
     ),
     _role as (
         insert
             into
                 roles(name,
                       tenant_id)
                 select
                     'service-admin',
                     tenant_id
                 from
                     _tenant
                 returning id as role_id
     ),
     _permission as(
         insert
             into
                 permissions(name,
                             description,
                             statement,service_id)
                 select 'manage-all',
                        'super admin can perform any action on any domain',
                        '{"action":"*","resource":"*","effect":"allow","fields":["*"]}',service_id from _service
                 returning id as permission_id
     ),
     _user as(
         insert
             into
                 users(user_id,service_id)
                 select $3,service_id from _service
                 returning id as user_id
     ),
     _tenant_user_role as
         (
             insert
                 into
                     tenant_users_roles(tenant_id,
                                        user_id,
                                        role_id)
                     select
                         tenant_id ,
                         user_id,
                         role_id
                     from
                         _tenant,
                         _user,
                         _role
                     returning role_id),
     _role_permission as
         (
             insert
                 into
                     role_permissions (role_id,
                                       permission_id)
                     select
                         role_id,
                         permission_id
                     from
                         _tenant_user_role,
                         _permission
                     returning role_id
         )
select
    service_id,
    password,
    service,
    service_status,
    tenant
from
    _service,
    _tenant
`

type CreateServiceParams struct {
	Name     string    `json:"name"`
	Password string    `json:"password"`
	UserID   uuid.UUID `json:"user_id"`
}

type CreateServiceRow struct {
	ServiceID     uuid.UUID `json:"service_id"`
	Password      string    `json:"password"`
	Service       string    `json:"service"`
	ServiceStatus Status    `json:"service_status"`
	Tenant        string    `json:"tenant"`
}

func (q *Queries) CreateService(ctx context.Context, arg CreateServiceParams) (CreateServiceRow, error) {
	row := q.db.QueryRow(ctx, createService, arg.Name, arg.Password, arg.UserID)
	var i CreateServiceRow
	err := row.Scan(
		&i.ServiceID,
		&i.Password,
		&i.Service,
		&i.ServiceStatus,
		&i.Tenant,
	)
	return i, err
}

const deleteService = `-- name: DeleteService :exec
DELETE FROM services WHERE id = $1
`

func (q *Queries) DeleteService(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteService, id)
	return err
}

const getServiceById = `-- name: GetServiceById :one
SELECT id, status, name, password, deleted_at, created_at, updated_at FROM services WHERE id = $1 AND deleted_at IS NULL
`

func (q *Queries) GetServiceById(ctx context.Context, id uuid.UUID) (Service, error) {
	row := q.db.QueryRow(ctx, getServiceById, id)
	var i Service
	err := row.Scan(
		&i.ID,
		&i.Status,
		&i.Name,
		&i.Password,
		&i.DeletedAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getServiceByName = `-- name: GetServiceByName :one
SELECT id, status, name, password, deleted_at, created_at, updated_at FROM services WHERE name = $1 and deleted_at IS NULL
`

func (q *Queries) GetServiceByName(ctx context.Context, name string) (Service, error) {
	row := q.db.QueryRow(ctx, getServiceByName, name)
	var i Service
	err := row.Scan(
		&i.ID,
		&i.Status,
		&i.Name,
		&i.Password,
		&i.DeletedAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const softDeleteService = `-- name: SoftDeleteService :one
UPDATE services set deleted_at = now() WHERE id = $1 AND deleted_at IS NULL returning id, status, name, password, deleted_at, created_at, updated_at
`

func (q *Queries) SoftDeleteService(ctx context.Context, id uuid.UUID) (Service, error) {
	row := q.db.QueryRow(ctx, softDeleteService, id)
	var i Service
	err := row.Scan(
		&i.ID,
		&i.Status,
		&i.Name,
		&i.Password,
		&i.DeletedAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateServiceStatus = `-- name: UpdateServiceStatus :one
UPDATE services SET status = $1 WHERE id = $2 AND deleted_at IS NULL RETURNING name,status
`

type UpdateServiceStatusParams struct {
	Status Status    `json:"status"`
	ID     uuid.UUID `json:"id"`
}

type UpdateServiceStatusRow struct {
	Name   string `json:"name"`
	Status Status `json:"status"`
}

func (q *Queries) UpdateServiceStatus(ctx context.Context, arg UpdateServiceStatusParams) (UpdateServiceStatusRow, error) {
	row := q.db.QueryRow(ctx, updateServiceStatus, arg.Status, arg.ID)
	var i UpdateServiceStatusRow
	err := row.Scan(&i.Name, &i.Status)
	return i, err
}
