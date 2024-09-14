package recommendation

import "github.com/shopspring/decimal"

type RestaurantRecommendation struct {
	RestaurantRecommendationID        int64 `gorm:"primaryKey"`
	RestaurantRecommendationRequestID int64
	RestaurantID                      int64
	DistanceInMeters                  decimal.Decimal
}

func NewRestaurantRecommendation(
	restaurantRecommendationRequestID int64,
	restaurantID int64,
	distanceInMeters decimal.Decimal,
) *RestaurantRecommendation {
	return &RestaurantRecommendation{
		RestaurantRecommendationRequestID: restaurantRecommendationRequestID,
		RestaurantID:                      restaurantID,
		DistanceInMeters:                  distanceInMeters,
	}
}

func (r *RestaurantRecommendation) TableName() string {
	return "restaurant_recommendation"
}
