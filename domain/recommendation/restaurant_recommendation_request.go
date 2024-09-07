package recommendation

import (
	"time"
)

type RestaurantRecommendationRequest struct {
	RestaurantRecommendationRequestID int64 `gorm:"primaryKey"`
	UserID                            *int64
	RequestedAt                       time.Time

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
