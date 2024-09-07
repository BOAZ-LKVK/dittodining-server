package recommendation

type UserLocation struct {
	Latitude  float64
	Longitude float64
}

func NewUserLocation(latitude float64, longitude float64) UserLocation {
	return UserLocation{Latitude: latitude, Longitude: longitude}
}
