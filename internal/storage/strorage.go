package storage

import (
	"2f-authorization/internal/constants/model/dto"
	"context"
)

type Policy interface {
	GetOpaData(ctx context.Context) ([]byte, error)
}

type Service interface {
	CreateService(ctx context.Context, param dto.CreateService) (*dto.CreateServiceResponse, error)
	IsServiceExist(ctx context.Context, param dto.CreateService) (bool, error)
}
