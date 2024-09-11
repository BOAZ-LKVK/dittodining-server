package recommendation

import (
	"errors"
	"github.com/BOAZ-LKVK/LKVK-server/domain/recommendation"
	"github.com/BOAZ-LKVK/LKVK-server/pkg/apicontroller"
	recommendation_repository "github.com/BOAZ-LKVK/LKVK-server/repository/recommendation"
	recommendation_service "github.com/BOAZ-LKVK/LKVK-server/service/recommendation"
	"github.com/BOAZ-LKVK/LKVK-server/service/restaurant/model"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"time"
)

type RecommendationAPIController struct {
	restaurantRecommendationRequestRepository recommendation_repository.RestaurantRecommendationRequestRepository
	restaurantRecommendationService           recommendation_service.RestaurantRecommendationService
	logger                                    *zap.Logger
}

func NewRecommendationAPIController(
	restaurantRecommendationRequestRepository recommendation_repository.RestaurantRecommendationRequestRepository,
	restaurantRecommendationService recommendation_service.RestaurantRecommendationService,
	logger *zap.Logger,
) *RecommendationAPIController {
	return &RecommendationAPIController{
		restaurantRecommendationRequestRepository: restaurantRecommendationRequestRepository,
		restaurantRecommendationService:           restaurantRecommendationService,
		logger:                                    logger,
	}
}

func (c *RecommendationAPIController) Pattern() string {
	return "/recommendation"
}

func (c *RecommendationAPIController) Handlers() []*apicontroller.APIHandler {
	return []*apicontroller.APIHandler{
		apicontroller.NewAPIHandler("", fiber.MethodPost, c.requestRestaurantRecommendation()),
		apicontroller.NewAPIHandler("/:restaurantRecommendationRequestID/restaurants", fiber.MethodGet, c.listRecommendedRestaurants()),
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
				Latitude:  (*request).UserLocation.Latitude.Decimal,
				Longitude: (*request).UserLocation.Longitude.Decimal,
			},
			time.Now(),
		)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		// TODO: restaurantRecommendations 생성 로직 추가 - 미리 추천 데이터를 만들고 노출하는 구조

		return ctx.JSON(&RequestRestaurantRecommendationResponse{
			RestaurantRecommendationRequestID: result.RestaurantRecommendationRequestID,
		})
	}
}

func (c *RecommendationAPIController) listRecommendedRestaurants() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		restaurantRecommendationRequestID, err := ctx.ParamsInt("restaurantRecommendationRequestID")
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		if _, err := c.restaurantRecommendationRequestRepository.FindByID(int64(restaurantRecommendationRequestID)); err != nil {
			if errors.Is(err, recommendation_repository.ErrRestaurantRecommendationRequestNotFound) {
				return ctx.Status(fiber.StatusNotFound).SendString(err.Error())
			}

			return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		// TODO: restaurants 조회 로직 추가
		// TODO: pagination 추가

		return ctx.JSON(&ListRecommendedRestaurantsResponse{
			RecommendedRestaurants: []RecommendedRestaurant{
				{
					Restaurant: model.Restaurant{},
					MenuItems:  nil,
					Review:     model.RestaurantReview{},
				},
			},
		})
	}
}
