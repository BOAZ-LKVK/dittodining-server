package zapfx

import (
	"context"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"log"
)

func NewZapLogger(lc fx.Lifecycle) *zap.Logger {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Can't initialize zap zapfx: %v", err)
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			err := logger.Sync()
			if err != nil {
				return err
			}
			return nil
		},
	})

	return logger
}
