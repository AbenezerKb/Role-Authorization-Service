package opa

import (
	"2f-authorization/internal/constants/dbinstance"
	errors "2f-authorization/internal/constants/error"
	"2f-authorization/internal/storage"
	"2f-authorization/platform/logger"
	"context"
	"encoding/json"

	"go.uber.org/zap"
)

type opa struct {
	db  dbinstance.DBInstance
	log logger.Logger
}

func Init(db dbinstance.DBInstance, log logger.Logger) storage.Policy {
	return &opa{
		db:  db,
		log: log,
	}
}

func (o *opa) GetOpaData(ctx context.Context) ([]byte, error) {
	data, err := o.db.GetOpaData(ctx)
	if err != nil {
		err := errors.ErrOpaUpdatePolicyError.Wrap(err, "can not update opa policy data")
		o.log.Error(ctx, "error getting opa data", zap.Error(err))
		return nil, err
	}
	
	opaData, err := json.Marshal(data)
	if err != nil {
		err := errors.ErrOpaUpdatePolicyError.Wrap(err, "can not update opa policy data")
		o.log.Error(ctx, "error getting opa data", zap.Error(err))
		return nil, err
	}
	return opaData, nil
}
