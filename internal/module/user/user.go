package user

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

type user struct {
	userPersistant storage.User
	log            logger.Logger
}

func Init(log logger.Logger, userPersistant storage.User) module.User {
	return &user{
		userPersistant: userPersistant,
		log:            log,
	}
}

func (u *user) RegisterUser(ctx context.Context, param dto.RegisterUser) error {
	var err error
	param.ServiceID, err = uuid.Parse(ctx.Value("x-service-id").(string))
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		u.log.Info(ctx, "invalid input", zap.Error(err), zap.Any("service-id", ctx.Value("x-service-id")))
		return err
	}

	if err = param.Validate(); err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		u.log.Info(ctx, "invalid input", zap.Error(err), zap.Any("input", param))
		return err
	}

	exists, err := u.userPersistant.IsUserExist(ctx, param)
	if err != nil {
		return err
	}

	if exists {
		err := errors.ErrDataExists.Wrap(err, "user with this id already exists")
		u.log.Info(ctx, "user already exists", zap.Error(err), zap.String("user-id", param.UserId.String()), zap.String("service-id", param.ServiceID.String()))
		return err
	}

	if err := u.userPersistant.RegiseterUser(ctx, param); err != nil {
		return err
	}
	return nil
}
