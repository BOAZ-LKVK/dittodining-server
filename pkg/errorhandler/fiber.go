package errorhandler

import (
	"fmt"
	"github.com/BOAZ-LKVK/LKVK-server/pkg/customerrors"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
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

		var applicationError *customerrors.ApplicationError
		if errors.As(err, &applicationError) {
			code = applicationError.Code
			message = applicationError.Err.Error()
		}

		// 500 에러인 경우는 message를 보안상 노출하지 않음
		if code == fiber.StatusInternalServerError {
			message = "Internal Server Error"
		}

		return ctx.Status(code).JSON(ErrorResponse{
			Code:    code,
			Message: message,
		})
	}
}
