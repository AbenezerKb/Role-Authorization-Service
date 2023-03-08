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
	_, sql := db_pgnflt.GetFilterSQL(filter)

	where := fmt.Sprintf("WHERE t.tenant_name ='%s' AND t.service_id ='%s' ", arg.TenantName, arg.ServiceID.String())

	if len(sql.Where) != 0 {
		where += fmt.Sprintf(" AND (%s)", sql.Where)
	}
	if len(sql.Search) != 0 {
		if len(sql.Where) != 0 {
			where += fmt.Sprintf(" AND (%s) ", sql.Search)
		} else {
			where += fmt.Sprintf(" WHERE (%s) ", sql.Search)
		}
	}

	orderBy := ""
	limitOffset := ""
	if len(sql.Sort) != 0 {
		orderBy += fmt.Sprintf(" ORDER BY %s ", sql.Sort)
	}

	if sql.Limit > 0 {
		limitOffset += fmt.Sprintf(" LIMIT %d  OFFSET %d ", sql.Limit, sql.Offset)
	}

	v := fmt.Sprintf(`
	WITH ur AS (
	    SELECT u.user_id,
	           json_agg(json_build_object('role_name', rl.name, 'status', tur.status, 'id', rl.id)) AS roles,
	           u.created_at
	    FROM tenant_users_roles AS tur
	             INNER JOIN tenants AS t ON t.id = tur.tenant_id
	             INNER JOIN roles AS rl ON rl.id = tur.role_id
	             INNER JOIN users AS u ON u.id = tur.user_id
		 %s
	  GROUP BY u.user_id, u.created_at
	)
	SELECT ur.user_id,
	       ur.roles,
	(select count(*) from ur)total_account FROM ur %s %s
`, where, orderBy, limitOffset)
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
