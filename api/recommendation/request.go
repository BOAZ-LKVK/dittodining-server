package recommendation

import (
	"errors"
)

type RequestRestaurantRecommendationRequest struct {
	UserLocation *UserLocationRequest `json:"userLocation"`
}

type UserLocationRequest struct {
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
}

func (r *RequestRestaurantRecommendationRequest) Validate() error {
	if r.UserLocation == nil {
		return errors.New("userLocation is required")
	}

	if r.UserLocation.Latitude == nil {
		return errors.New("userLocation.latitude is required")
	}

	if r.UserLocation.Longitude == nil {
		return errors.New("userLocation.longitude is required")
	}

	return nil
}
