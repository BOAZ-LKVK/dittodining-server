package recommendation

import "github.com/shopspring/decimal"

type UserLocation struct {
	Latitude  decimal.Decimal `gorm:"type:decimal(11,8)"`
	Longitude decimal.Decimal `gorm:"type:decimal(11,8)"`
}

func NewUserLocation(latitude decimal.Decimal, longitude decimal.Decimal) UserLocation {
	return UserLocation{Latitude: latitude, Longitude: longitude}
}
