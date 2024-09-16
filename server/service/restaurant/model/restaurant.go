package model

import (
	restaurant_domain "github.com/BOAZ-LKVK/LKVK-server/server/domain/restaurant"
	"github.com/shopspring/decimal"
)

type Restaurant struct {
	RestaurantID          int64                             `json:"restaurantId"`
	Name                  string                            `json:"name"`
	Address               string                            `json:"address"`
	Description           string                            `json:"description"`
	MaximumPricePerPerson decimal.Decimal                   `json:"maximumPricePerPerson"`
	MinimumPricePerPerson decimal.Decimal                   `json:"minimumPricePerPerson"`
	Longitude             decimal.Decimal                   `json:"longitude"`
	Latitude              decimal.Decimal                   `json:"latitude"`
	BusinessHours         []*restaurant_domain.BusinessHour `json:"businessHours"`
	RestaurantImageURLs   []string                          `json:"restaurantImageUrls"`
}
