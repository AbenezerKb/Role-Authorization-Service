package dbinstance

import (
	"2f-authorization/internal/constants/model"
	"2f-authorization/internal/constants/model/dto"
	"context"
	"fmt"

	"github.com/google/uuid"
	db_pgnflt "gitlab.com/2ftimeplc/2fbackend/repo/db-pgnflt"
)

type GetTenantUsersRoles struct {
	TenantName string    `json:"tenant_name"`
	ServiceID  uuid.UUID `json:"service_id"`
}

//	func ComposeSelectColumnsQuery(columns []string, tableName, filterSQL string) string {
//		return fmt.Sprintf("SELECT %v,count(*) over() FROM %s %s", strings.Join(columns, ","), tableName, filterSQL)
//	}
func (db *DBInstance) GetTenantUsersWithRoles(ctx context.Context, filter db_pgnflt.FilterParams, arg GetTenantUsersRoles) ([]dto.TenantUserRoles, *model.MetaData, error) {
	query := fmt.Sprintf(`select
array_agg(rl.name)::string[] as roles,
us.user_id

from
(
	select
		id,
		service_id,
	from
		tenants
	where
		tenant_name = '%s'
) tn
INNER JOIN (

	select
		role_id,
		tenant_id,
		user_id
	from
		tenant_users_roles

) tur ON tn.id = tur.tenant_id
INNER JOIN (
	select
		id,
		name
	from
		roles
) rl ON tur.role_id = rl.id
INNER JOIN (
	select
		id,
		user_id
	from
		users
) us ON us.id = tur.user_id

  GROUP  BY us.user_id
  `, arg.TenantName)
	// var count int64
	metadata := &model.MetaData{FilterParams: filter}
	var filtParam = db_pgnflt.FilterWithCustomWhere(fmt.Sprintf("service_id= '%s'", arg.ServiceID), filter)
	rows, err := db.Pool.Query(ctx, db_pgnflt.ComposeSelectColumnsQuery([]string{"user_id", "roles"}, query, filtParam))
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()
	var items []dto.TenantUserRoles
	for rows.Next() {
		var i dto.TenantUserRoles
		if err := rows.Scan(&i.UserId, &i.Roles, &metadata.Total); err != nil {
			return nil, nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, nil, err
	}

	return items, metadata, nil
}
