package errorhandler

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func NewFiberErrorHandler(logger *zap.Logger) func(ctx *fiber.Ctx, err error) error {
	return func(ctx *fiber.Ctx, err error) error {
		// Zap 로거를 사용하여 에러 로그 출력
		logger.Error("Unhandled error occurred",
			zap.String("method", ctx.Method()),
			zap.String("url", ctx.OriginalURL()),
			zap.Error(err),
			zap.String("stacktrace", fmt.Sprintf("%+v", err)),
		)

		return ctx.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}
}
