package initiator

import (
	"2f-authorization/platform/logger"
	opa_platform "2f-authorization/platform/opa"
	"context"
	"os"

	"go.uber.org/zap"
)

func InitOpa(ctx context.Context, path string, persistence Persistence, log logger.Logger) opa_platform.Opa {
	policy, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(ctx, "error reading the policy file")
	}

	op := opa_platform.Init(string(policy), persistence.opa, log)
	if err := op.Refresh(ctx, "initiate"); err != nil {
		log.Fatal(ctx, "error getting opa data", zap.Error(err))
	}
	return op
}
