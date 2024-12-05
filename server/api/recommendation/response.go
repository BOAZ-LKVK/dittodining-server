package recommendation

import "github.com/BOAZ-LKVK/LKVK-server/server/service/recommendation/model"

type RequestRestaurantRecommendationResponse struct {
	RestaurantRecommendationRequestID int64 `json:"restaurantRecommendationRequestId"`
	IsAvailableLocation               bool  `json:"isAvailableLocation"`
}

type ListRecommendedRestaurantsResponse struct {
	RecommendedRestaurants []*model.RecommendedRestaurant `json:"recommendedRestaurants"`
	NextCursor             *string                        `json:"nextCursor"`
}

type SelectRestaurantRecommendationsResponse struct{}

type GetRestaurantRecommendationResultResponse struct {
	Results []*model.RestaurantRecommendationResult `json:"results"`
}

type GetRestaurantRecommendationResponse struct {
	Recommendation *model.RecommendedRestaurant `json:"recommendation"`
}
