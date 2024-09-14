package recommendation

type SelectedRestaurantRecommendation struct {
	SelectedRestaurantRecommendationID int64 `gorm:"primaryKey"`
	RestaurantRecommendationRequestID  int64
	RestaurantRecommendationID         int64
	RestaurantID                       int64
}

func (s *SelectedRestaurantRecommendation) TableName() string {
	return "selected_restaurant_recommendation"
}

func NewSelectedRestaurantRecommendation(
	restaurantRecommendationRequestID int64,
	restaurantRecommendationID int64,
	restaurantID int64,
) *SelectedRestaurantRecommendation {
	return &SelectedRestaurantRecommendation{
		RestaurantRecommendationRequestID: restaurantRecommendationRequestID,
		RestaurantRecommendationID:        restaurantRecommendationID,
		RestaurantID:                      restaurantID,
	}
}
