package initiator

import (
	"2f-authorization/platform/logger"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

func InitDB(url string, idleConnTimeout time.Duration, log logger.Logger) *pgxpool.Pool {
	config, err := pgxpool.ParseConfig(url)
	if err != nil {
		log.Fatal(context.Background(), fmt.Sprintf("Failed to connect to database: %v", err))
	}
	config.ConnConfig.Logger = log.Named("pgx")
	if idleConnTimeout == 0 {
		idleConnTimeout = 4 * time.Minute
	}
	config.MaxConnIdleTime = idleConnTimeout

	conn, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		log.Fatal(context.Background(), fmt.Sprintf("Failed to connect to database: %v", err))
	}
	return conn
}
