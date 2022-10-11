package dbinstance

import (
	"2f-authorization/internal/constants/model/db"

	"github.com/jackc/pgx/v4/pgxpool"
)

type DBInstance struct {
	*db.Queries
	Pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) DBInstance {
	return DBInstance{
		Queries: db.New(pool),
		Pool:    pool,
	}
}

