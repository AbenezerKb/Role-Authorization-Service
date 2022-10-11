package module

import (
	"2f-authorization/internal/constants/model/dto"
	"context"
)

type Service interface {
	CreateService(ctx context.Context, param dto.CreateService) (*dto.CreateServiceResponse, error)
}
