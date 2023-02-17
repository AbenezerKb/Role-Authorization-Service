package dbinstance

import (
	"2f-authorization/internal/constants/model"
	"2f-authorization/internal/constants/model/dto"
	"context"
	"fmt"

	db_pgnflt "gitlab.com/2ftimeplc/2fbackend/repo/db-pgnflt"

	"github.com/google/uuid"
)

type ListRolesParams struct {
	TenantName string    `json:"tenant_name"`
	ServiceID  uuid.UUID `json:"service_id"`
}

func (db *DBInstance) ListRoles(ctx context.Context, filter db_pgnflt.FilterParams, arg ListRolesParams) ([]dto.Role, *model.MetaData, error) {

	metadata := &model.MetaData{FilterParams: filter}
	_, filterParam := db_pgnflt.GetFilterSQL(filter)
	filterParam.Where = fmt.Sprintf("WHERE (tenant_name = '%s' AND service_id= '%s') AND (%s)", arg.TenantName, arg.ServiceID, filterParam.Where)
	rows, err := db.Pool.Query(ctx, db_pgnflt.GetSelectColumnsQuery([]string{"name", "created_at", "id", "status", "updated_at"}, "role_tenant", filterParam))
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()
	var items []dto.Role
	for rows.Next() {
		var i dto.Role
		if err := rows.Scan(&i.Name, &i.CreatedAt, &i.ID, &i.Status, &i.UpdatedAt, &metadata.Total); err != nil {
			return nil, nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, nil, err
	}
	return items, metadata, nil
}
