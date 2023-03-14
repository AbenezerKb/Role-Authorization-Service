package initiator

import (
	"2f-authorization/platform/logger"
	opa_platform "2f-authorization/platform/opa"
	"context"
	"fmt"
	"os"

	"go.uber.org/zap"
)

func InitOpa(ctx context.Context, rego, data, server string, persistence Persistence, log logger.Logger) opa_platform.Opa {

	policy, err := os.ReadFile(rego)
	if err != nil {
		log.Fatal(ctx, "error reading the policy file")
	}
	fmt.Println("file path", data)
	op := opa_platform.Init(string(policy), persistence.opa, data, rego, server, log)
	if err := op.Refresh(ctx, "initiate"); err != nil {
		log.Fatal(ctx, "error getting opa data", zap.Error(err))
	}

	return op
}
