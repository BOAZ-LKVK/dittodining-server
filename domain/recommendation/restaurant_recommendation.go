package recommendation

import "github.com/shopspring/decimal"

type RestaurantRecommendation struct {
	RestaurantRecommendationID        int64 `gorm:"primaryKey"`
	RestaurantRecommendationRequestID int64
	RestaurantID                      int64
	DistanceInMeters                  decimal.Decimal
}

func (r *RestaurantRecommendation) TableName() string {
	return "restaurant_recommendation"
}
