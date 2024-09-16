package model

type GetRestaurantRecommendationResultResult struct {
	Results []*RestaurantRecommendationResult `json:"results"`
}

type RestaurantRecommendationResult struct {
	RestaurantRecommendationID int64                  `json:"restaurantRecommendationId"`
	RecommendedRestaurant      *RecommendedRestaurant `json:"recommendedRestaurant"`
}
