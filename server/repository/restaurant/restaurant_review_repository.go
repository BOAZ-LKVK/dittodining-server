package restaurant

import (
	"github.com/BOAZ-LKVK/LKVK-server/server/domain/restaurant"
	"gorm.io/gorm"
)

type RestaurantReviewRepository interface {
	FindAllByRestaurantID(restaurantID int64) ([]*restaurant.RestaurantReview, error)
	FindAllByRestaurantIDs(restaurantIDs []int64) ([]*restaurant.RestaurantReview, error)
	CountByRestaurantID(restaurantID int64) (int64, error)
}

func NewRestaurantReviewRepository(db *gorm.DB) RestaurantReviewRepository {
	return &restaurantReviewRepository{db: db}
}

type restaurantReviewRepository struct {
	db *gorm.DB
}

func (r *restaurantReviewRepository) FindAllByRestaurantID(restaurantID int64) ([]*restaurant.RestaurantReview, error) {
	var reviews []*restaurant.RestaurantReview
	result := r.db.
		Where(restaurant.RestaurantReview{
			RestaurantID: restaurantID,
		}).
		Find(&reviews)
	if result.Error != nil {
		return nil, result.Error
	}

	return reviews, nil
}

func (r *restaurantReviewRepository) FindAllByRestaurantIDs(restaurantIDs []int64) ([]*restaurant.RestaurantReview, error) {
	var reviews []*restaurant.RestaurantReview

	result := r.db.
		Where("restaurant_id IN ?", restaurantIDs).
		Find(&reviews)
	if result.Error != nil {
		return nil, result.Error
	}

	return reviews, nil
}

func (r *restaurantReviewRepository) CountByRestaurantID(restaurantID int64) (int64, error) {
	var count int64
	result := r.db.
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
