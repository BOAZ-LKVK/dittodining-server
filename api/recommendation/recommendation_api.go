package recommendation

import (
	"github.com/BOAZ-LKVK/LKVK-server/domain/recommendation"
	"github.com/BOAZ-LKVK/LKVK-server/pkg/apicontroller"
	recommendation_repository "github.com/BOAZ-LKVK/LKVK-server/repository/recommendation"
	recommendation_service "github.com/BOAZ-LKVK/LKVK-server/service/recommendation"
	"github.com/gofiber/fiber/v2"
	"time"
)

type RecommendationAPIController struct {
	restaurantRecommendationRequestRepository recommendation_repository.RestaurantRecommendationRequestRepository
	restaurantRecommendationService           recommendation_service.RestaurantRecommendationService
}

func NewRecommendationAPIController(restaurantRecommendationRequestRepository recommendation_repository.RestaurantRecommendationRequestRepository, restaurantRecommendationService recommendation_service.RestaurantRecommendationService) *RecommendationAPIController {
	return &RecommendationAPIController{restaurantRecommendationRequestRepository: restaurantRecommendationRequestRepository, restaurantRecommendationService: restaurantRecommendationService}
}

func (c *RecommendationAPIController) Pattern() string {
	return "/recommendation"
}

func (c *RecommendationAPIController) Handlers() []*apicontroller.APIHandler {
	return []*apicontroller.APIHandler{
		apicontroller.NewAPIHandler("", fiber.MethodPost, c.requestRestaurantRecommendation()),
	}
}

func (c *RecommendationAPIController) requestRestaurantRecommendation() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		request := new(RequestRestaurantRecommendationRequest)
		if err := ctx.BodyParser(request); err != nil {
			return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		if err := request.Validate(); err != nil {
			return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		result, err := c.restaurantRecommendationService.RequestRestaurantRecommendation(
			nil,
			recommendation.UserLocation{
				Latitude:  *request.UserLocation.Latitude,
				Longitude: *request.UserLocation.Longitude,
			},
			time.Now(),
		)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		return ctx.JSON(&RequestRestaurantRecommendationResponse{
			RestaurantRecommendationRequestID: result.RestaurantRecommendationRequestID,
		})
	}
}
