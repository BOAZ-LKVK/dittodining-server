package fiberfx

import (
	"context"
	"fmt"
	"github.com/BOAZ-LKVK/LKVK-server/pkg/apicontroller"
	"github.com/BOAZ-LKVK/LKVK-server/pkg/errorhandler"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Params struct {
	fx.In

	Config *Config

	APIControllers []apicontroller.APIController `group:"api_controllers"`
	Logger         *zap.Logger
}

type Result struct {
	fx.Out

	Server *fiber.App
}

var Module = fx.Module("server",
	fx.Provide(
		parseConfig,
		New,
	),
	fx.Invoke(startServer),
)

func New(lc fx.Lifecycle, p Params) (Result, error) {
	app := fiber.New(fiber.Config{
		ErrorHandler: errorhandler.NewFiberErrorHandler(p.Logger),
	})

	for _, c := range p.APIControllers {
		p.Logger.Info(fmt.Sprintf("Registering API controller: %s", c.Pattern()))
		group := app.Group(c.Pattern())
		for _, h := range c.Handlers() {
			group.Add(h.Method, h.Pattern, h.Handler)
		}
	}

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() {
				p.Logger.Info(fmt.Sprintf("Starting server on port %d", p.Config.Port))
				if err := app.Listen(fmt.Sprintf(":%d", p.Config.Port)); err != nil {
					p.Logger.Error("Failed to start server", zap.Error(err))
				}
			}()

			return nil
		},
		OnStop: func(context.Context) error {
			if err := app.Shutdown(); err != nil {
				return err
			}
			p.Logger.Info("Server stopped")
			return nil
		},
	})

	return Result{
		Server: app,
	}, nil
}

func startServer(_ *fiber.App) {}
