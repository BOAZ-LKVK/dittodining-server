package apihandler

import "github.com/gofiber/fiber/v2"

type APIHandler struct {
	Pattern string
	Method  string
	Handler fiber.Handler
}

func NewAPIHandler(pattern string, method string, handler fiber.Handler) *APIHandler {
	return &APIHandler{
		Pattern: pattern,
		Method:  method,
		Handler: handler,
	}
}
