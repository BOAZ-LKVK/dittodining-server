package recommendation

import (
	"time"
)

type RestaurantRecommendationRequest struct {
	RestaurantRecommendationRequestID int64 `gorm:"primaryKey"`
	UserID                            *int64
	UserLocation                      UserLocation `gorm:"embedded"`
	RequestedAt                       time.Time
}

func (r *RestaurantRecommendationRequest) TableName() string {
	return "restaurant_recommendation_request"
}

func NewRestaurantRecommendationRequest(userID *int64, userLocation UserLocation, requestedAt time.Time) *RestaurantRecommendationRequest {
	return &RestaurantRecommendationRequest{UserID: userID, UserLocation: userLocation, RequestedAt: requestedAt}
}
