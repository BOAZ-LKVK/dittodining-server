package restaurant

import (
	"github.com/shopspring/decimal"
	"time"
)

type RestaurantReview struct {
	RestaurantReviewID int64 `gorm:"primaryKey"`
	RestaurantID       int64
	WriterName         string
	Score              decimal.NullDecimal
	Content            *string
	WroteAt            time.Time
}

func (r *RestaurantReview) TableName() string {
	return "restaurant_review"
}
