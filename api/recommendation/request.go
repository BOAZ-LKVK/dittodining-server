package recommendation

type RequestRestaurantRecommendationRequest struct {
	latitude  float64
	longitude float64
}

type RequestRestaurantRecommendationResponse struct {
	recommendationID int64
}
