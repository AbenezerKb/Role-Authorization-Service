package permission

import (
	"2f-authorization/internal/constants/dbinstance"
	errors "2f-authorization/internal/constants/error"
	"2f-authorization/internal/constants/model/db"
	"2f-authorization/internal/constants/model/dto"
	"2f-authorization/internal/storage"
	"2f-authorization/platform/logger"
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	"go.uber.org/zap"
)

type permission struct {
	db  dbinstance.DBInstance
	log logger.Logger
}

func Init(db dbinstance.DBInstance, log logger.Logger) storage.Permission {
	return &permission{
		db:  db,
		log: log,
	}
}

func (p *permission) CreatePermission(ctx context.Context, param dto.CreatePermission) (uuid.UUID, error) {

	statement, _ := param.Statement.Value()
	permissionId, err := p.db.CreateOrGetPermission(ctx, db.CreateOrGetPermissionParams{
		Name:        param.Name,
		ServiceID:   param.ServiceID,
		Description: param.Description,
		Statment: pgtype.JSON{
			Bytes:  statement,
			Status: pgtype.Present,
		},
	})

	if err != nil {
		err := errors.ErrWriteError.Wrap(err, "could not create or get permission")
		p.log.Error(ctx, "unable to create or get permission", zap.Error(err), zap.Any("permission", param))
		return uuid.UUID{}, err
	}

	return permissionId, nil
}

func (p *permission) AddToDomain(ctx context.Context, permissionId, domain uuid.UUID) error {
	if err := p.db.AssignDomain(ctx, db.AssignDomainParams{
		ID:           domain,
		PermissionID: permissionId,
	}); err != nil {
		err := errors.ErrWriteError.Wrap(err, "could not assign domain to permission")
		p.log.Error(ctx, "unable to  assign domain to permission", zap.Error(err), zap.String("permission", permissionId.String()), zap.String("domain", domain.String()))
		return err
	}
	return nil
}
