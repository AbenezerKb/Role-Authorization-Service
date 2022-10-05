package initiator

import (
	"2f-authorization/platform/logger"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Persistence struct {
}

func InitPersistence(conn *pgxpool.Pool, log logger.Logger) Persistence {
	return Persistence{}
}
