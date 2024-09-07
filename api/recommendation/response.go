package recommendation

import (
	"github.com/BOAZ-LKVK/LKVK-server/service/restaurant/model"
)

type RequestRestaurantRecommendationResponse struct {
	RestaurantRecommendationRequestID int64 `json:"restaurantRecommendationRequestID"`
}

type RecommendedRestaurant struct {
	Restaurant model.Restaurant       `json:"restaurant"`
	MenuItems  []model.RestaurantMenu `json:"menuItems"`
	Review     model.RestaurantReview `json:"review"`
}

type ListRecommendedRestaurantsResponse struct {
	RecommendedRestaurants []RecommendedRestaurant `json:"recommendedRestaurants"`
}
