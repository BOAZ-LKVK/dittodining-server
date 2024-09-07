package model

import "github.com/BOAZ-LKVK/LKVK-server/domain/restaurant"

type Restaurant struct {
	RestaurantID        int64                     `json:"restaurantId"`
	Name                string                    `json:"name"`
	Catchphrase         string                    `json:"catchphrase"`
	PriceRangePerPerson string                    `json:"priceRangePerPerson"`
	Distance            string                    `json:"distance"`
	BusinessHours       []restaurant.BusinessHour `json:"businessHours"`
	RestaurantImageUrls []string                  `json:"restaurantImageUrls"`
}
