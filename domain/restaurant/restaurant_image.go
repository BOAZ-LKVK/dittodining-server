package restaurant

import (
	"time"
)

type RestaurantImage struct {
	RestaurantImageID int64 `gorm:"primaryKey"`
	RestaurantID      int64
	ImageURL          string

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
