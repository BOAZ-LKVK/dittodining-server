package restaurant

import (
	"github.com/shopspring/decimal"
)

type RestaurantMenu struct {
	RestaurantMenuID int64 `gorm:"primaryKey"`
	RestaurantID     int64
	Name             string
	Price            decimal.Decimal
	Description      *string
	ImageURL         *string
}

func (m *RestaurantMenu) TableName() string {
	return "restaurant_menu"
}
