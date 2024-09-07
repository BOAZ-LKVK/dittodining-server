package recommendation

import (
	"time"
)

type RestaurantRecommendationRequest struct {
	RestaurantRecommendationRequestID int64 `gorm:"primaryKey"`
	UserID                            *int64
	UserLocation                      UserLocation `gorm:"embedded"`

	RequestedAt time.Time

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func NewRestaurantRecommendationRequest(userID *int64, userLocation UserLocation, requestedAt time.Time) *RestaurantRecommendationRequest {
	return &RestaurantRecommendationRequest{UserID: userID, UserLocation: userLocation, RequestedAt: requestedAt}
}
