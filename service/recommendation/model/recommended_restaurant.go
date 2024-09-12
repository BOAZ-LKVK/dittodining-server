package model

import (
	"github.com/BOAZ-LKVK/LKVK-server/domain/restaurant"
	"github.com/BOAZ-LKVK/LKVK-server/service/restaurant/model"
)

type RecommendedRestaurant struct {
	Restaurant RestaurantRecommendation `json:"restaurant"`
	MenuItems  []*model.RestaurantMenu  `json:"menuItems"`
	Review     model.RestaurantReview   `json:"review"`
}

type RestaurantRecommendation struct {
	RestaurantID        int64                      `json:"restaurantId"`
	Name                string                     `json:"name"`
	Description         string                     `json:"description"`
	PriceRangePerPerson string                     `json:"priceRangePerPerson"`
	Distance            string                     `json:"distance"`
	BusinessHours       []*restaurant.BusinessHour `json:"businessHours"`
	RestaurantImageURLs []string                   `json:"restaurantImageUrls"`
}
