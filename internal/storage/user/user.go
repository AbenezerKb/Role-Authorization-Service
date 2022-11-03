package user

import (
	"2f-authorization/internal/constants/dbinstance"
	errors "2f-authorization/internal/constants/error"
	"2f-authorization/internal/constants/error/sqlcerr"
	"2f-authorization/internal/constants/model/db"
	"2f-authorization/internal/constants/model/dto"
	"2f-authorization/internal/storage"
	"2f-authorization/platform/logger"
	"context"

	"go.uber.org/zap"
)

type user struct {
	db  dbinstance.DBInstance
	log logger.Logger
}

func Init(db dbinstance.DBInstance, log logger.Logger) storage.User {
	return &user{
		db:  db,
		log: log,
	}
}

func (u *user) RegiseterUser(ctx context.Context, param dto.RegisterUser) error {

	if err := u.db.RegisterUser(ctx, db.RegisterUserParams{
		UserID:    param.UserId,
		ServiceID: param.ServiceID,
	}); err != nil {
		err = errors.ErrWriteError.Wrap(err, "could not register user")
		u.log.Error(ctx, "unable to register user", zap.Error(err), zap.Any("user-id", param.UserId), zap.Any("service-id", param.ServiceID))
		return err
	}
	return nil
}

func (u *user) IsUserExist(ctx context.Context, param dto.RegisterUser) (bool, error) {
	_, err := u.db.GetUserWithUserIdAndServiceId(ctx, db.GetUserWithUserIdAndServiceIdParams{
		UserID:    param.UserId,
		ServiceID: param.ServiceID,
	})
	if err != nil {
		if sqlcerr.Is(err, sqlcerr.ErrNoRows) {
			return false, nil
		} else {
			err = errors.ErrReadError.Wrap(err, "error reading the user data")
			u.log.Error(ctx, "unable to get service data", zap.Error(err), zap.Any("service-id", param.ServiceID), zap.Any("user-id", param.UserId))
			return false, err
		}
	}
	return true, nil
}

func (u *user) UpdateUserStatus(ctx context.Context, param dto.UpdateUserStatus) error {
	_, err := u.db.UpdateUserStatus(ctx, db.UpdateUserStatusParams{
		Status:    db.Status(param.Status),
		UserID:    param.UserID,
		ServiceID: param.ServiceID,
	})

	if err != nil {
		if sqlcerr.Is(err, sqlcerr.ErrNoRows) {
			err := errors.ErrNoRecordFound.Wrap(err, "user not found")
			u.log.Error(ctx, "error changing user's status", zap.Error(err), zap.String("service", param.ServiceID.String()), zap.String("user-status", param.Status), zap.String("user-id", param.UserID.String()))
			return err
		} else {
			err = errors.ErrUpdateError.Wrap(err, "error changing user's status")
			u.log.Error(ctx, "error changing user's status", zap.Error(err), zap.String("service", param.ServiceID.String()), zap.String("user-status", param.Status), zap.String("user-id", param.UserID.String()))
			return err
		}
	}
	return nil
}
