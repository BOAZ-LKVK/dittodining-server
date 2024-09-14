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
