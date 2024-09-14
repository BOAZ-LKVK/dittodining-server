package model

type GetRestaurantRecommendationResultResult struct {
	Results []*RestaurantRecommendationResult `json:"results"`
}

type RestaurantRecommendationResult struct {
	RestaurantRecommendationID int64                  `json:"restaurantRecommendationId"`
	Restaurant                 *RecommendedRestaurant `json:"restaurant"`
}
