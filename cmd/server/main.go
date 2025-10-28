package main

import (
	"context"
	"fmt"

	"github.com/WithSoull/ChatServer/internal/app"
	"github.com/WithSoull/platform_common/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	appCtx := context.Background()

	a, err := app.NewApp(appCtx)
	if err != nil {
		panic(fmt.Sprintf("failed to init app with error: %v", err))
	}

	if err := a.Run(appCtx); err != nil {
		logger.Fatal(appCtx, "failed to run app", zap.Error(err))
	}
}
