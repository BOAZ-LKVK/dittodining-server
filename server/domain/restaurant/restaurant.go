package restaurant

import (
	"github.com/shopspring/decimal"
)

type Restaurant struct {
	RestaurantID          int64 `gorm:"primaryKey"`
	Name                  string
	Address               string
	Description           string
	MaximumPricePerPerson decimal.Decimal
	MinimumPricePerPerson decimal.Decimal
	Longitude             decimal.Decimal
	Latitude              decimal.Decimal
	// BusinessHoursJSON is BusinessHour structs in JSON format
	BusinessHoursJSON         string
	RecommendationScore       decimal.Decimal
	AverageScoreFromNaver     decimal.Decimal
	TotalReviewCountFromNaver int64
	TotalReviewCountFromKakao int64
	AverageScoreFromKakao     decimal.Decimal
	TotalReviewCount          int64
}

func (r *Restaurant) TableName() string {
	return "restaurant"
}
