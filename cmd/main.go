package main

import (
	"github.com/BOAZ-LKVK/LKVK-server/pkg/apihandler"
	"github.com/BOAZ-LKVK/LKVK-server/pkg/apihandler/sample"
	"github.com/BOAZ-LKVK/LKVK-server/pkg/errorhandler"
	sample_repository "github.com/BOAZ-LKVK/LKVK-server/pkg/repository/sample"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"log"
)

func main() {
	// Zap 로거 설정
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Can't initialize zap logger: %v", err)
	}
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			log.Fatalf("Can't sync zap logger: %v", err)
		}
	}(logger)

	app := fiber.New(fiber.Config{
		ErrorHandler: errorhandler.NewFiberErrorHandler(logger),
	})

	apiHandler := app.Group("/api")
	handlers := []apihandler.APIController{
		sample.NewSampleAPIHandler(
			sample_repository.NewSampleRepository(),
		),
	}

	for _, c := range handlers {
		group := apiHandler.Group(c.Pattern())
		for _, h := range c.Handlers() {
			group.Add(h.Method, h.Pattern, h.Handler)
		}
	}

	if err := app.Listen(":3000"); err != nil {
		panic(err)
	}
}
