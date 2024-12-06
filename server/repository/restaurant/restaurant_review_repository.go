package restaurant

import (
	"context"
	"github.com/BOAZ-LKVK/LKVK-server/server/domain/restaurant"
	"gorm.io/gorm"
)

type RestaurantReviewRepository interface {
	// TODO: limit 추가
	FindAllByRestaurantID(ctx context.Context, db *gorm.DB, restaurantID int64) ([]*restaurant.RestaurantReview, error)
	// TODO: limit 추가
	FindAllByRestaurantIDs(ctx context.Context, db *gorm.DB, restaurantIDs []int64) ([]*restaurant.RestaurantReview, error)
	CountByRestaurantID(ctx context.Context, db *gorm.DB, restaurantID int64) (int64, error)
}

func NewRestaurantReviewRepository() RestaurantReviewRepository {
	return &restaurantReviewRepository{}
}

type restaurantReviewRepository struct{}

func (r *restaurantReviewRepository) FindAllByRestaurantID(ctx context.Context, db *gorm.DB, restaurantID int64) ([]*restaurant.RestaurantReview, error) {
	var reviews []*restaurant.RestaurantReview
	result := db.
		Where(restaurant.RestaurantReview{
			RestaurantID: restaurantID,
		}).
		Find(&reviews)
	if result.Error != nil {
		return nil, result.Error
	}

	return reviews, nil
}

func (r *restaurantReviewRepository) FindAllByRestaurantIDs(ctx context.Context, db *gorm.DB, restaurantIDs []int64) ([]*restaurant.RestaurantReview, error) {
	var reviews []*restaurant.RestaurantReview

	result := db.
		Where("restaurant_id IN ?", restaurantIDs).
		Find(&reviews)
	if result.Error != nil {
		return nil, result.Error
	}

	return reviews, nil
}

func (r *restaurantReviewRepository) CountByRestaurantID(ctx context.Context, db *gorm.DB, restaurantID int64) (int64, error) {
	var count int64
	result := db.
		Model(&restaurant.RestaurantReview{}).
		Where(restaurant.RestaurantReview{
			RestaurantID: restaurantID,
		}).
		Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}
