package apihandler

import "github.com/gofiber/fiber/v2"

type APIHandler struct {
	Pattern string
	Method  string
	Handler fiber.Handler
}

type APIHandlerGroup interface {
	Pattern() string
	Handlers() []*APIHandler
}
