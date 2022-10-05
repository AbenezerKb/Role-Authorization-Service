package initiator

import (
	"go.uber.org/zap"
	"log"
)

func InitLogger() *zap.Logger {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}
	return logger
}
