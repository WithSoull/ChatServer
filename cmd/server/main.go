package main

import (
	"context"

	"github.com/WithSoull/ChatServer/internal/app"
	"github.com/WithSoull/platform_common/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	appCtx := context.Background()

	a, err := app.NewApp(appCtx)
	if err != nil {
		logger.Fatal(appCtx, "failed to init app", zap.Error(err))
	}

	if err := a.Run(); err != nil {
		logger.Fatal(appCtx, "failed to run app", zap.Error(err))
	}
}
