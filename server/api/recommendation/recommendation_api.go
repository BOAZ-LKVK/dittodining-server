package recommendation

import (
	"github.com/BOAZ-LKVK/LKVK-server/pkg/apicontroller"
	"github.com/BOAZ-LKVK/LKVK-server/server/domain/recommendation"
	recommendation_service "github.com/BOAZ-LKVK/LKVK-server/server/service/recommendation"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"strconv"
	"time"
)

type RecommendationAPIController struct {
	restaurantRecommendationService recommendation_service.RestaurantRecommendationService
	logger                          *zap.Logger
}

func NewRecommendationAPIController(
	restaurantRecommendationService recommendation_service.RestaurantRecommendationService,
	logger *zap.Logger,
) *RecommendationAPIController {
	return &RecommendationAPIController{
		restaurantRecommendationService: restaurantRecommendationService,
		logger:                          logger,
	}
}

func (c *RecommendationAPIController) Pattern() string {
	return "/api/recommendation"
}

func (c *RecommendationAPIController) Handlers() []*apicontroller.APIHandler {
	return []*apicontroller.APIHandler{
		apicontroller.NewAPIHandler("/request", fiber.MethodPost, c.requestRestaurantRecommendation()),
		apicontroller.NewAPIHandler("/request/:restaurantRecommendationRequestId/restaurants", fiber.MethodGet, c.listRecommendedRestaurants()),
		apicontroller.NewAPIHandler("/request/:restaurantRecommendationRequestId/restaurants/select", fiber.MethodPost, c.selectRestaurantRecommendations()),
		apicontroller.NewAPIHandler("/request/:restaurantRecommendationRequestId/result", fiber.MethodGet, c.getRestaurantRecommendationResult()),
		apicontroller.NewAPIHandler("/recommendations/:restaurantRecommendationId", fiber.MethodPost, c.getRestaurantRecommendation()),
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

		return ctx.JSON(&RequestRestaurantRecommendationResponse{
			RestaurantRecommendationRequestID: result.RestaurantRecommendationRequestID,
		})
	}
}

func (c *RecommendationAPIController) listRecommendedRestaurants() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		restaurantRecommendationRequestID, err := ctx.ParamsInt("restaurantRecommendationRequestId")
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
		}
		limit := ctx.QueryInt("limit", 10)
		if limit == 0 {
			return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
		}
		cursorRestaurantRecommendationIDQuery := ctx.Query("cursorRestaurantRecommendationId", "")
		var cursorRestaurantRecommendationID *int64
		if cursorRestaurantRecommendationIDQuery != "" {
			parse, err := strconv.ParseInt(cursorRestaurantRecommendationIDQuery, 10, 64)
			if err != nil {
				return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
			}

			cursorRestaurantRecommendationID = &parse
		}

		listRecommendedRestaurantsResult, err := c.restaurantRecommendationService.ListRecommendedRestaurants(
			int64(restaurantRecommendationRequestID),
			cursorRestaurantRecommendationID,
			lo.ToPtr(int64(limit)),
		)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		return ctx.JSON(&ListRecommendedRestaurantsResponse{
			RecommendedRestaurants: listRecommendedRestaurantsResult.RecommendedRestaurants,
			NextCursor:             listRecommendedRestaurantsResult.NextCursor,
		})
	}
}

func (c *RecommendationAPIController) selectRestaurantRecommendations() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		restaurantRecommendationRequestID, err := ctx.ParamsInt("restaurantRecommendationRequestId")
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
		}
		request, err := parseRequestBody[SelectRestaurantRecommendationsRequest](ctx)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		if _, err := c.restaurantRecommendationService.SelectRestaurantRecommendation(int64(restaurantRecommendationRequestID), request.RestaurantRecommendationIDs); err != nil {
			return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		return ctx.JSON(&SelectRestaurantRecommendationsResponse{})
	}
}

func (c *RecommendationAPIController) getRestaurantRecommendationResult() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		restaurantRecommendationRequestID, err := ctx.ParamsInt("restaurantRecommendationRequestId")
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		result, err := c.restaurantRecommendationService.GetRestaurantRecommendationResult(int64(restaurantRecommendationRequestID))
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		return ctx.JSON(&GetRestaurantRecommendationResultResponse{
			Results: result.Results,
		})
	}
}

func (c *RecommendationAPIController) getRestaurantRecommendation() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		restaurantRecommendationID, err := ctx.ParamsInt("restaurantRecommendationId")
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		result, err := c.restaurantRecommendationService.GetRestaurantRecommendation(int64(restaurantRecommendationID))
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		return ctx.JSON(&GetRestaurantRecommendationResponse{
			Recommendation: result.RecommendedRestaurant,
		})
	}
}

// TODO: refactor 적절한 곳으로 옮기기
func parseRequestBody[T any](ctx *fiber.Ctx) (*T, error) {
	request := new(T)
	if err := ctx.BodyParser(request); err != nil {
		return nil, err
	}

	if v, ok := any(request).(Validator); ok {
		if err := v.Validate(); err != nil {
			return nil, err
		}
	}

	return request, nil
}
