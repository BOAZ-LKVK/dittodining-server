package model

import "github.com/shopspring/decimal"

type RestaurantMenu struct {
	RestaurantMenuID int64           `json:"restaurantMenuId"`
	ImageURL         *string         `json:"imageUrl"`
	Name             string          `json:"name"`
	Price            decimal.Decimal `json:"price"`
	Description      *string         `json:"description"`
}
