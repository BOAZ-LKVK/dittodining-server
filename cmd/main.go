package main

import (
	sample_api "github.com/BOAZ-LKVK/LKVK-server/api/sample"
	"github.com/BOAZ-LKVK/LKVK-server/pkg/fiberfx"
	"github.com/BOAZ-LKVK/LKVK-server/pkg/zapfx"
	sample_repository "github.com/BOAZ-LKVK/LKVK-server/repository/sample"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func main() {
	fx.New(
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger { return &fxevent.ZapLogger{Logger: log} }),
		fx.Provide(
			zapfx.NewZapLogger,
		),
		fx.Provide(
			sample_repository.NewSampleRepository,
			fiberfx.AsAPIController(
				sample_api.NewSampleAPIHandler,
			),
		),
		fiberfx.Module,
	).Run()
}
