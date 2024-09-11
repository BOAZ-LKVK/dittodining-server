package recommendation

import (
	"errors"
	"github.com/shopspring/decimal"
)

type RequestRestaurantRecommendationRequest struct {
	UserLocation *UserLocationRequest `json:"userLocation"`
}

type UserLocationRequest struct {
	Latitude  decimal.NullDecimal `json:"latitude"`
	Longitude decimal.NullDecimal `json:"longitude"`
}

func (r *RequestRestaurantRecommendationRequest) Validate() error {
	if r.UserLocation == nil {
		return errors.New("userLocation is required")
	}

	if !r.UserLocation.Latitude.Valid {
		return errors.New("userLocation.latitude is required")
	}

	if !r.UserLocation.Longitude.Valid {
		return errors.New("userLocation.longitude is required")
	}

	return nil
}
