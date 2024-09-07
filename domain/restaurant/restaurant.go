package restaurant

import (
	"time"
)

type Restaurant struct {
	RestaurantID        int64 `gorm:"primaryKey"`
	Name                string
	Catchphrase         string
	PriceRangePerPerson string
	Distance            string
	BusinessHoursJSON   string
	RestaurantImages    []RestaurantImage `gorm:"references:RestaurantID"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
