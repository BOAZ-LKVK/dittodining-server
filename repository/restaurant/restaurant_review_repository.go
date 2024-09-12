package restaurant

import (
	"github.com/BOAZ-LKVK/LKVK-server/domain/restaurant"
	"gorm.io/gorm"
)

type RestaurantReviewRepository interface {
	FindAllByRestaurantID(restaurantID int64) ([]*restaurant.RestaurantReview, error)
	FindAllByRestaurantIDs(restaurantIDs []int64) ([]*restaurant.RestaurantReview, error)
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
