package recommendation

import (
	"github.com/BOAZ-LKVK/LKVK-server/pkg/apicontroller"
	"github.com/BOAZ-LKVK/LKVK-server/pkg/customerrors"
	"github.com/BOAZ-LKVK/LKVK-server/server/domain/recommendation"
	recommendation_repository "github.com/BOAZ-LKVK/LKVK-server/server/repository/recommendation"
	restaurant_repository "github.com/BOAZ-LKVK/LKVK-server/server/repository/restaurant"
	recommendation_service "github.com/BOAZ-LKVK/LKVK-server/server/service/recommendation"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
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
		apicontroller.NewAPIHandler("/recommendations/:restaurantRecommendationId", fiber.MethodGet, c.getRestaurantRecommendation()),
	}
}

func (c *RecommendationAPIController) requestRestaurantRecommendation() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		request := new(RequestRestaurantRecommendationRequest)
		if err := ctx.BodyParser(request); err != nil {
			return &customerrors.ApplicationError{
				Code: fiber.StatusBadRequest,
				Err:  errors.New("Invalid RestaurantRecommendation Request body"),
			}
		}

		// assign new error with return error string of validate method
		if err := request.Validate(); err != nil {
			return &customerrors.ApplicationError{
				Code: fiber.StatusBadRequest,
				// err.Error() function gets error message
				Err: errors.Errorf("Invalid RestaurantRecommendation Request body: %s", err.Error()),
			}
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
			return &customerrors.ApplicationError{
				Code: fiber.StatusInternalServerError,
				Err:  err,
			}
		}

		return ctx.JSON(&RequestRestaurantRecommendationResponse{
			RestaurantRecommendationRequestID: result.RestaurantRecommendationRequestID,
		})
	}
}

func (c *RecommendationAPIController) listRecommendedRestaurants() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// get parameter value from restaurantRecommendationRequestId
		restaurantRecommendationRequestID, err := ctx.ParamsInt("restaurantRecommendationRequestId")
		if err != nil {
			return &customerrors.ApplicationError{
				Code: fiber.StatusBadRequest,
				Err:  errors.New("Cannot Convert restaurantRecommendationRequestId route parameter into integer"),
			}
		}
		// get querystring value from limit, if key is existed but cannot assign value then default value set
		// If They don't have key parameters then error value in 0
		limit := ctx.QueryInt("limit", 10)
		if limit == 0 {
			return &customerrors.ApplicationError{
				Code: fiber.StatusBadRequest,
				Err:  errors.New("Cannot Convert limit route parameter into integer"),
			}
		}

		cursorRestaurantRecommendationIDQuery := ctx.Query("cursorRestaurantRecommendationId", "")
		var cursorRestaurantRecommendationID *int64
		if cursorRestaurantRecommendationIDQuery != "" {
			parse, err := strconv.ParseInt(cursorRestaurantRecommendationIDQuery, 10, 64)
			if err != nil {
				return &customerrors.ApplicationError{
					Code: fiber.StatusBadRequest,
					Err:  errors.New("Cannot Convert cursorRestaurantRecommendationId query parameter into integer"),
				}
			}
			cursorRestaurantRecommendationID = &parse
		}

		listRecommendedRestaurantsResult, err := c.restaurantRecommendationService.ListRecommendedRestaurants(
			int64(restaurantRecommendationRequestID),
			cursorRestaurantRecommendationID,
			lo.ToPtr(int64(limit)),
		)
		if err != nil {
			return &customerrors.ApplicationError{
				Code: fiber.StatusInternalServerError,
				Err:  err,
			}
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
			return &customerrors.ApplicationError{
				Code: fiber.StatusBadRequest,
				Err:  errors.New("Cannot Convert restaurantRecommendationId route parameter into integer"),
			}
		}
		request, err := parseRequestBody[SelectRestaurantRecommendationsRequest](ctx)
		if err != nil {
			return &customerrors.ApplicationError{
				Code: fiber.StatusBadRequest,
				Err:  err,
			}
		}

		if _, err := c.restaurantRecommendationService.SelectRestaurantRecommendation(int64(restaurantRecommendationRequestID), request.RestaurantRecommendationIDs); err != nil {
			return &customerrors.ApplicationError{
				Code: fiber.StatusInternalServerError,
				Err:  err,
			}
		}

		return ctx.JSON(&SelectRestaurantRecommendationsResponse{})
	}
}

func (c *RecommendationAPIController) getRestaurantRecommendationResult() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		restaurantRecommendationRequestID, err := ctx.ParamsInt("restaurantRecommendationRequestId")
		if err != nil {
			return &customerrors.ApplicationError{
				Code: fiber.StatusBadRequest,
				Err:  errors.New("Cannot Convert restaurantRecommendationId route parameter into integer"),
			}
		}

		result, err := c.restaurantRecommendationService.GetRestaurantRecommendationResult(int64(restaurantRecommendationRequestID))
		if err != nil {
			return &customerrors.ApplicationError{
				Code: fiber.StatusInternalServerError,
				Err:  err,
			}
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
			return &customerrors.ApplicationError{
				Code: fiber.StatusBadRequest,
				Err:  errors.New("Cannot Convert restaurantRecommendationId route parameter into integer"),
			}
		}

		result, err := c.restaurantRecommendationService.GetRestaurantRecommendation(int64(restaurantRecommendationID))
		if err != nil {
			if errors.Is(err, restaurant_repository.ErrRestaurantNotFound) {
				return &customerrors.ApplicationError{
					Code: fiber.StatusNotFound,
					Err:  errors.New("Cannot found restaurant provided from restaurantRecommendationId"),
				}
			}

			if errors.Is(err, recommendation_repository.ErrRestaurantRecommendationNotFound) {
				return &customerrors.ApplicationError{
					Code: fiber.StatusNotFound,
					Err:  errors.New("Cannot found recommend restaurant list"),
				}
			}

			return &customerrors.ApplicationError{
				Code: fiber.StatusInternalServerError,
				Err:  errors.New("GetRestaurantRecommendation service function error"),
			}
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
		return nil, errors.New("Cannot get value from request body")
	}

	if v, ok := any(request).(Validator); ok {
		if err := v.Validate(); err != nil {

			return nil, errors.New("Invalidate request body value")
		}
	}

	return request, nil
}
