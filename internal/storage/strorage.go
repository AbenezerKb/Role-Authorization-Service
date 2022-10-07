package storage

import "context"

type Policy interface {
	GetOpaData(ctx context.Context) ([]byte, error)
}
