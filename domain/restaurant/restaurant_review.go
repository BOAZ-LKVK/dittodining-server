package restaurant

import "time"

type RestaurantReview struct {
	RestaurantReviewID int64 `gorm:"primaryKey"`
	RestaurantID       int64
	WriterName         string
	Score              float32
	Content            string
	WrittenAt          time.Time

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
