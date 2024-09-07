package recommendation

import (
	"time"
)

type RestaurantRecommendation struct {
	RestaurantRecommendationID        int64 `gorm:"primaryKey"`
	RestaurantRecommendationRequestID int64
	RestaurantID                      int64

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
