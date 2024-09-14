package recommendation

import "github.com/BOAZ-LKVK/LKVK-server/server/service/recommendation/model"

type RequestRestaurantRecommendationResponse struct {
	RestaurantRecommendationRequestID int64 `json:"restaurantRecommendationRequestID"`
}

type ListRecommendedRestaurantsResponse struct {
	RecommendedRestaurants []*model.RecommendedRestaurant `json:"recommendedRestaurants"`
	NextCursor             *string                        `json:"nextCursor"`
}

type SelectRestaurantRecommendationsResponse struct{}
