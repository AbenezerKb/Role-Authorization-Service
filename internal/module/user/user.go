package user

import (
	errors "2f-authorization/internal/constants/error"
	"2f-authorization/internal/constants/model/dto"
	"2f-authorization/internal/module"
	"2f-authorization/internal/storage"
	"2f-authorization/platform/logger"
	"2f-authorization/platform/opa"
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type user struct {
	userPersistant storage.User
	log            logger.Logger
	opa            opa.Opa
}

func Init(log logger.Logger, userPersistant storage.User, opa opa.Opa) module.User {
	return &user{
		userPersistant: userPersistant,
		log:            log,
		opa:            opa,
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

	exists, err := u.userPersistant.CheckIfUserExists(ctx, param)
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

func (u *user) UpdateUserStatus(ctx context.Context, param dto.UpdateUserStatus) error {
	var err error
	param.ServiceID, err = uuid.Parse(ctx.Value("x-service-id").(string))
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		u.log.Info(ctx, "invalid input", zap.Error(err), zap.Any("service-id", ctx.Value("x-service-id")))
		return err
	}

	if err = param.Validate(); err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		u.log.Info(ctx, "invalid input", zap.Error(err))
		return err
	}

	if err = u.userPersistant.UpdateUserStatus(ctx, param); err != nil {
		return err
	}

	return u.opa.Refresh(ctx, fmt.Sprintf("Updating user [%v] with status [%v]", param.UserID.String(), param.Status))
}
