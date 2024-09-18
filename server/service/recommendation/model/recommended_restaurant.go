package model

import (
	restaurant_domain "github.com/BOAZ-LKVK/LKVK-server/server/domain/restaurant"
	restaurant_model "github.com/BOAZ-LKVK/LKVK-server/server/service/restaurant/model"
	"github.com/shopspring/decimal"
)

type RecommendedRestaurant struct {
	Restaurant RestaurantRecommendation           `json:"restaurant"`
	MenuItems  []*restaurant_model.RestaurantMenu `json:"menuItems"`
	Review     restaurant_model.RestaurantReview  `json:"review"`
}

type RestaurantRecommendation struct {
	RestaurantRecommendationID int64                             `json:"restaurantRecommendationId"`
	RestaurantID               int64                             `json:"restaurantId"`
	Name                       string                            `json:"name"`
	Description                string                            `json:"description"`
	Location                   Location                          `json:"location"`
	MaximumPricePerPerson      decimal.Decimal                   `json:"maximumPricePerPerson"`
	MinimumPricePerPerson      decimal.Decimal                   `json:"minimumPricePerPerson"`
	PriceRangePerPerson        string                            `json:"priceRangePerPerson"`
	DistanceInMeters           decimal.Decimal                   `json:"distanceInMeters"`
	BusinessHours              []*restaurant_domain.BusinessHour `json:"businessHours"`
	RestaurantImageURLs        []string                          `json:"restaurantImageUrls"`
}
