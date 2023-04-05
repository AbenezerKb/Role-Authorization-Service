package dbinstance

import (
	"2f-authorization/internal/constants/model/dto"
	"context"

	"github.com/google/uuid"
)

const getUserPermissionWithInTenant = `-- name: GetUserPermissionWithInTenant :many
with _user as(
    select id from users u where u.user_id=$1 and u.service_id=$2 and u.deleted_at is null and u.status='ACTIVE'
),_tenant as(
    select id from tenants t where t.tenant_name=$3 and t.service_id=$2 and t.deleted_at is null
)
select p.name,p.status,p.created_at,p.statement,p.id,p.description,coalesce(json_agg(json_build_object('name',p2.name,'status',p2.status,'description',p2.description,'statement',p2.statement,'created_at',p2.created_at,'id',p2.id))filter ( where p2.deleted_at is null and p2.status='ACTIVE' ),'[]') AS inherited_permissions
from   _user,_tenant,tenant_users_roles tur left join roles r on tur.role_id = r.id left join role_permissions rp on r.id = rp.role_id left join permissions p on p.id = rp.permission_id left join permissions_hierarchy ph on p.id = ph.parent left join permissions p2 on p2.id = ph.child where tur.deleted_at is null and tur.status='ACTIVE'
and tur.user_id=_user.id and  tur.tenant_id=_tenant.id and p.status='ACTIVE' and p.deleted_at is null and r.status='ACTIVE' and r.deleted_at is null group by  p.name,p.status,p.created_at,p.statement,p.id,p.description
`

type GetUserPermissionWithInTenantParams struct {
	UserID     uuid.UUID `json:"user_id"`
	TenantName string    `json:"tenant_name"`
	ServiceID  uuid.UUID `json:"service_id"`
}

func (q *DBInstance) GetUserPermissionWithInTenant(ctx context.Context, arg GetUserPermissionWithInTenantParams) ([]dto.Permission, error) {
	rows, err := q.Pool.Query(ctx, getUserPermissionWithInTenant, arg.UserID, arg.ServiceID, arg.TenantName)
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
			&i.CreatedAt,
			&i.Statement,
			&i.ID,
			&i.Description,
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

const getUserPermissionWithInTheDomain = `-- name: GetUserPermissionWithInTheDomain :many
with _user as(
    select id from users u where u.user_id=$1 and u.service_id=$2 and u.deleted_at is null and u.status='ACTIVE'
), _permissions as(
    select p.name,p.status,p.created_at,p.statement,p.id,p.description,t.tenant_name ,coalesce(json_agg(json_build_object('name',p2.name,'status',p2.status,'description',p2.description,'statement',p2.statement,'created_at',p2.created_at,'id',p2.id))filter ( where p2.deleted_at is null and p2.status='ACTIVE' ),'[]') AS inherited_permissions  from   _user,tenants t,tenant_users_roles tur left join roles r on tur.role_id = r.id left join role_permissions rp on r.id = rp.role_id left join permissions p on p.id = rp.permission_id left join permissions_hierarchy ph on p.id = ph.parent left join permissions p2 on p2.id = ph.child where p.status='ACTIVE' and p.deleted_at is null  and tur.deleted_at is null and tur.status='ACTIVE' and tur.tenant_id=t.id  and tur.user_id=_user.id  and t.domain_id=$3 and t.deleted_at is null and t.status='ACTIVE' and t.service_id=$2 and r.status='ACTIVE' group by  p.name,p.status,p.created_at,p.statement,p.id,p.description,p.deleted_at,t.tenant_name
)
select tenant_name as tenant,coalesce(json_agg(_permissions)filter ( where _permissions.status='ACTIVE' ),'[]')  as permissions from _permissions group by tenant_name
`

type GetUserPermissionWithInTheDomainParams struct {
	UserID    uuid.UUID `json:"user_id"`
	ServiceID uuid.UUID `json:"service_id"`
	DomainID  uuid.UUID `json:"domain_id"`
}

func (q *DBInstance) GetUserPermissionWithInTheDomain(ctx context.Context, arg GetUserPermissionWithInTheDomainParams) ([]dto.DomainPermissions, error) {
	rows, err := q.Pool.Query(ctx, getUserPermissionWithInTheDomain, arg.UserID, arg.ServiceID, arg.DomainID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []dto.DomainPermissions
	for rows.Next() {
		var i dto.DomainPermissions
		if err := rows.Scan(&i.Tenant, &i.Permissions); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
