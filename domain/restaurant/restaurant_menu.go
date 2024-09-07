package restaurant

import "time"

type RestaurantMenu struct {
	RestaurantMenuID int64 `gorm:"primaryKey"`
	RestaurantID     int64
	Name             string
	Price            int64
	Description      string
	ImageURL         *string

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
