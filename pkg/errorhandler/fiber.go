package errorhandler

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewFiberErrorHandler(logger *zap.Logger) fiber.ErrorHandler {
	return func(ctx *fiber.Ctx, err error) error {
		logger.Error("Unhandled error occurred",
			zap.String("method", ctx.Method()),
			zap.String("url", ctx.OriginalURL()),
			zap.String("stacktrace", fmt.Sprintf("%+v", err)),
			zap.Error(err),
		)

		code := fiber.StatusInternalServerError
		message := "Internal Server Error"

		var e *fiber.Error
		if errors.As(err, &e) {
			code = e.Code
			message = e.Message
		}

		return ctx.Status(code).JSON(ErrorResponse{
			Code:    code,
			Message: message,
		})
	}
}
