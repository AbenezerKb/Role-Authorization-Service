package permission

import (
	"2f-authorization/internal/constants/dbinstance"
	errors "2f-authorization/internal/constants/error"
	"2f-authorization/internal/constants/error/sqlcerr"
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
		Statement: pgtype.JSON{
			Bytes:  statement,
			Status: pgtype.Present,
		},
		Column5: param.Domain,
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

func (p *permission) ListAllPermission(ctx context.Context, param dto.GetAllPermissionsReq) ([]dto.Permission, error) {
	permission, err := p.db.ListPermissions(ctx, dbinstance.ListPermissionsParams{
		TenantName: param.TenantName,
		ServiceID:  param.ServiceID,
	})
	if err != nil {
		if sqlcerr.Is(err, sqlcerr.ErrNoRows) {
			err := errors.ErrNoRecordFound.Wrap(err, "no permisisons found")
			p.log.Info(ctx, "no permissions were found", zap.Error(err), zap.String("tenany-name", param.TenantName), zap.String("service-id", param.ServiceID.String()))
			return []dto.Permission{}, err
		} else {
			err = errors.ErrReadError.Wrap(err, "error reading permissions")
			p.log.Error(ctx, "error reading permissions", zap.Error(err), zap.String("tenany-name", param.TenantName), zap.String("service-id", param.ServiceID.String()))
			return []dto.Permission{}, err
		}
	}
	return permission, nil
}

func (p *permission) CreatePermissionDependency(ctx context.Context, param dto.CreatePermissionDependency, serviceId uuid.UUID) error {
	if err := p.db.CreatePermissionDependency(ctx, db.CreatePermissionDependencyParams{
		Name:      param.PermissionName,
		Column3:   param.InheritedPermissions,
		ServiceID: serviceId,
	}); err != nil {
		err = errors.ErrWriteError.Wrap(err, "could not create inheritance between permissions")
		p.log.Error(ctx, "unable to create inheritance", zap.Error(err), zap.Any("parent-permission", param.PermissionName), zap.Any("child-permissions", param.InheritedPermissions))
		return err
	}

	return nil
}

func (p *permission) DeletePermission(ctx context.Context, serviceId, permissionId uuid.UUID, tenantName string) error {

	if _, err := p.db.DeletePermissions(ctx, db.DeletePermissionsParams{
		ID:         permissionId,
		ServiceID:  serviceId,
		TenantName: tenantName,
	}); err != nil {
		if sqlcerr.Is(err, sqlcerr.ErrNoRows) {
			err := errors.ErrNoRecordFound.Wrap(err, "permission does not exists")
			p.log.Info(ctx, "permission not found", zap.Error(err), zap.String("service-id", serviceId.String()), zap.String("permission-id", permissionId.String()))
			return err
		}
		err = errors.ErrDBDelError.Wrap(err, "error deleting permission")
		p.log.Error(ctx, "error deleting permission", zap.Error(err), zap.String("service-id", serviceId.String()), zap.String("permission-id", permissionId.String()))
		return err
	}
	return nil
}

func (p *permission) CanBeDeletedOrUpdated(ctx context.Context, permissionId, serviceId uuid.UUID) (bool, error) {
	delete_or_update, err := p.db.CanBeDeletedOrUpdated(ctx, db.CanBeDeletedOrUpdatedParams{
		ID:        permissionId,
		ServiceID: serviceId,
	})
	if err != nil {
		if sqlcerr.Is(err, sqlcerr.ErrNoRows) {
			err := errors.ErrNoRecordFound.Wrap(err, "permission does not exists")
			p.log.Info(ctx, "permission not found", zap.Error(err), zap.String("service-id", serviceId.String()), zap.String("permission-id", permissionId.String()))
			return false, err
		}
		err := errors.ErrReadError.Wrap(err, "could not read permission data")
		p.log.Error(ctx, "unable to read the permission data", zap.Error(err), zap.Any("permission id", permissionId), zap.Any("service id", serviceId), zap.String("permission-id", permissionId.String()))
		return false, err
	}

	return delete_or_update, nil
}

func (p *permission) GetPermission(ctx context.Context, permissionId, serviceId uuid.UUID, tenantName string) (*dto.Permission, error) {
	permission, err := p.db.GetPermissionDetails(ctx, dbinstance.GetPermissionDetailsParams{
		TenantName: tenantName,
		ServiceID:  serviceId,
		ID:         permissionId,
	})
	if err != nil {
		if sqlcerr.Is(err, sqlcerr.ErrNoRows) {
			err := errors.ErrNoRecordFound.Wrap(err, "permission does not exists")
			p.log.Info(ctx, "permission not found", zap.Error(err), zap.String("service-id", serviceId.String()), zap.String("permission-id", permissionId.String()), zap.String("tenant", tenantName))
			return nil, err
		}
		err := errors.ErrReadError.Wrap(err, "could not read permission data")
		p.log.Error(ctx, "unable to read the permission data", zap.Error(err), zap.Any("permission id", permissionId), zap.Any("service id", serviceId), zap.String("tenant", tenantName))
		return nil, err
	}

	return &permission, nil
}

func (p *permission) UpdatePermissionStatus(ctx context.Context, param dto.UpdatePermissionStatus, permissionId, serviceId uuid.UUID, tenant string) error {
	_, err := p.db.UpdatePermissionStatus(ctx, db.UpdatePermissionStatusParams{
		ID:         permissionId,
		TenantName: tenant,
		ServiceID:  serviceId,
		Status:     db.Status(param.Status),
	})
	if err != nil {
		if sqlcerr.Is(err, sqlcerr.ErrNoRows) {
			err := errors.ErrNoRecordFound.Wrap(err, "permission does not exists")
			p.log.Info(ctx, "permission not found", zap.Error(err), zap.String("service-id", serviceId.String()), zap.String("permission-id", permissionId.String()))
			return err
		}
		err = errors.ErrUpdateError.Wrap(err, "error changing permission's status")
		p.log.Error(ctx, "error changing permission's status", zap.Error(err), zap.String("service", serviceId.String()), zap.String("permission-status", param.Status), zap.String("permission-id", permissionId.String()), zap.String("tenant", tenant))
		return err
	}
	return nil
}
