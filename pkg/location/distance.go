package location

import (
	"github.com/golang/geo/s2"
	"github.com/shopspring/decimal"
)

func CalculateDistanceInMeters(lat1, lon1, lat2, lon2 decimal.Decimal) decimal.Decimal {
	// 지구 반지름 (미터 단위)
	radius := 6371000.0

	latlng1 := s2.LatLngFromDegrees(lat1.InexactFloat64(), lon1.InexactFloat64())
	latlng2 := s2.LatLngFromDegrees(lat2.InexactFloat64(), lon2.InexactFloat64())

	return decimal.NewFromFloat(latlng1.Distance(latlng2).Radians() * radius)
}
