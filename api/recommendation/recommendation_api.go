package recommendation

import (
	recommendation_domain "github.com/BOAZ-LKVK/LKVK-server/domain/recommendation"
	"github.com/BOAZ-LKVK/LKVK-server/pkg/apicontroller"
	"github.com/BOAZ-LKVK/LKVK-server/repository/recommendation"
	"github.com/gofiber/fiber/v2"
	"time"
)

type RecommendationAPIController struct {
	restaurantRecommendationRequestRepository recommendation.RestaurantRecommendationRequestRepository
}

func NewRecommendationAPIController(restaurantRecommendationRequestRepository recommendation.RestaurantRecommendationRequestRepository) *RecommendationAPIController {
	return &RecommendationAPIController{restaurantRecommendationRequestRepository: restaurantRecommendationRequestRepository}
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

		recommendationRequest := recommendation_domain.NewRestaurantRecommendationRequest(
			nil,
			recommendation_domain.NewUserLocation(
				*request.UserLocation.Latitude, *request.UserLocation.Longitude,
			),
			// TODO: testablity를 위해 clock interface 개발 후 대체
			time.Now(),
		)
		created, err := c.restaurantRecommendationRequestRepository.Create(recommendationRequest)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		return ctx.JSON(&RequestRestaurantRecommendationResponse{
			RestaurantRecommendationRequestID: created.RestaurantRecommendationRequestID,
		})
	}
}
