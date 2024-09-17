package api

import (
	"github.com/BOAZ-LKVK/LKVK-server/pkg/apicontroller"
	"github.com/gofiber/fiber/v2"
)

type HealthCheckAPIController struct {
}

func NewHomeAPIController() *HealthCheckAPIController {
	return &HealthCheckAPIController{}
}

func (c *HealthCheckAPIController) Pattern() string {
	return ""
}

func (c *HealthCheckAPIController) Handlers() []*apicontroller.APIHandler {
	return []*apicontroller.APIHandler{
		apicontroller.NewAPIHandler("/health", "GET", c.healthCheck()),
	}
}

func (c *HealthCheckAPIController) healthCheck() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusOK).SendString("OK")
	}
}
