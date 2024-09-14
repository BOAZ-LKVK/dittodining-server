package model

import (
	"encoding/json"
	restaurant2 "github.com/BOAZ-LKVK/LKVK-server/server/domain/restaurant"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type Restaurant struct {
	RestaurantID          int64                       `json:"restaurantId"`
	Name                  string                      `json:"name"`
	Address               string                      `json:"address"`
	Description           string                      `json:"description"`
	MaximumPricePerPerson decimal.Decimal             `json:"maximumPricePerPerson"`
	MinimumPricePerPerson decimal.Decimal             `json:"minimumPricePerPerson"`
	Longitude             decimal.Decimal             `json:"longitude"`
	Latitude              decimal.Decimal             `json:"latitude"`
	BusinessHours         []*restaurant2.BusinessHour `json:"businessHours"`
	RestaurantImageURLs   []string                    `json:"restaurantImageUrls"`
}

func ConvertRestaurantEntityToModel(r *restaurant2.Restaurant) (*Restaurant, error) {
	var businessHours []*restaurant2.BusinessHour
	if err := json.Unmarshal([]byte(r.BusinessHoursJSON), &businessHours); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal business hours")
	}

	restaurantImageURLs := make([]string, 0)
	for _, image := range r.RestaurantImages {
		restaurantImageURLs = append(restaurantImageURLs, image.ImageURL)
	}

	return &Restaurant{
		RestaurantID:          r.RestaurantID,
		Name:                  r.Name,
		Address:               r.Address,
		Description:           r.Description,
		MaximumPricePerPerson: r.MaximumPricePerPerson,
		MinimumPricePerPerson: r.MinimumPricePerPerson,
		Longitude:             r.Longitude,
		Latitude:              r.Latitude,
		BusinessHours:         businessHours,
		RestaurantImageURLs:   restaurantImageURLs,
	}, nil
}
