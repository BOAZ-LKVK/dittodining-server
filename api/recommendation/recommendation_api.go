package recommendation

import (
	"github.com/BOAZ-LKVK/LKVK-server/pkg/apicontroller"
	"github.com/gofiber/fiber/v2"
)

type RecommendationAPIController struct {
}

func (a *RecommendationAPIController) Pattern() string {
	return "/recommendation"
}

func (a *RecommendationAPIController) Handlers() []*apicontroller.APIHandler {
	return []*apicontroller.APIHandler{
		apicontroller.NewAPIHandler("", fiber.MethodPost, a.requestRestaurantRecommendation()),
	}
}

func (a *RecommendationAPIController) requestRestaurantRecommendation() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		request := new(RequestRestaurantRecommendationRequest)
		if err := ctx.BodyParser(request); err != nil {
			return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		return ctx.JSON(&RequestRestaurantRecommendationResponse{
			recommendationID: 0,
		})
	}
}
