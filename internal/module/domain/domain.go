package domain

import (
	errors "2f-authorization/internal/constants/error"
	"2f-authorization/internal/constants/model/dto"
	"2f-authorization/internal/module"
	"2f-authorization/internal/storage"
	"2f-authorization/platform/logger"
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type domain struct {
	domainPersistant storage.Domain
	log              logger.Logger
}

func Init(log logger.Logger, domainPersistant storage.Domain) module.Domain {
	return &domain{

		domainPersistant: domainPersistant,
		log:              log,
	}
}

func (d *domain) CreateDomain(ctx context.Context, param dto.CreateDomain) (*dto.Domain, error) {

	var err error
	if err = param.Validated(); err != nil {

		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		d.log.Info(ctx, "invalid input", zap.Error(err))
		return nil, err
	}

	isExist, err := d.domainPersistant.IsDomainExist(ctx, param)
	if err != nil {
		d.log.Info(ctx, "domain already exists", zap.String("name", param.Name))
		return nil, errors.ErrDataExists.Wrap(err, "domain  with this name and service already exists")
	}

	if isExist {
		d.log.Info(ctx, "domain already exists", zap.String("name", param.Name))
		return nil, errors.ErrDataExists.Wrap(err, "domain with this name and service already exists")
	}

	domain, err := d.domainPersistant.CreateDomain(ctx, param)
	if err != nil {
		return nil, err
	}

	return domain, nil

}
func (d *domain) DeleteDomain(ctx context.Context, param dto.DeleteDomain) error {
	var err error
	param.ServiceID, err = uuid.Parse(ctx.Value("x-service-id").(string))
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		d.log.Info(ctx, "invalid input", zap.Error(err))
		return err
	}
	err = param.Validate()
	if err != nil {
		err = errors.ErrInvalidUserInput.Wrap(err, "error on validating input")
		d.log.Error(ctx, "error validating input of domain", zap.Error(err), zap.String("service-id", param.ServiceID.String()))
		return err

	}
	if err := d.domainPersistant.SoftDeleteDomain(ctx, param); err != nil {

		return err
	}

	return nil

}
