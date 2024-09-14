package restaurant

type RestaurantImage struct {
	RestaurantImageID int64 `gorm:"primaryKey"`
	RestaurantID      int64
	ImageURL          string
}

func (i *RestaurantImage) TableName() string {
	return "restaurant_image"
}
