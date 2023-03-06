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
	UserID     uuid.UUID `json:"user_id"`
}

func (db *DBInstance) GetTenantUsersWithRoles(ctx context.Context, filter db_pgnflt.FilterParams, arg GetTenantUsersRoles) ([]dto.TenantUserRoles, *model.MetaData, error) {
	filterParam := db_pgnflt.GetFilterSQLWithCustomWhere(
		fmt.Sprintf("t.tenant_name = '%s' and t.service_id = '%s' and tur.deleted_at is null and u.user_id !='%s'", arg.TenantName, arg.ServiceID, arg.UserID), filter)
	filterParam.GroupBy = "u.user_id, c_x_t_y_b.total_count,u.created_at"
	var v = db_pgnflt.GetSelectColumnsQueryWithJoins([]string{"u.user_id", "json_agg(json_build_object('role_name',rl.name,'status',tur.status,'id',rl.id)) as roles"},
		db_pgnflt.Table{Name: "tenant_users_roles", Alias: "tur"}, []db_pgnflt.JOIN{
			{
				Table: db_pgnflt.Table{
					Alias: "t",
					Name:  "tenants",
				},
				JoinType: "inner join",
				On:       "t.id=tur.tenant_id",
			},
			{
				Table: db_pgnflt.Table{
					Alias: "rl",
					Name:  "roles",
				},
				JoinType: "inner join",
				On:       "rl.id=tur.role_id",
			},
			{
				Table: db_pgnflt.Table{
					Alias: "u",
					Name:  "users",
				},
				JoinType: "inner join",
				On:       "u.id=tur.user_id",
			},
		}, filterParam,
	)

	metadata := &model.MetaData{FilterParams: filter}

	rows, err := db.Pool.Query(ctx, v)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()
	items := []dto.TenantUserRoles{}
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
