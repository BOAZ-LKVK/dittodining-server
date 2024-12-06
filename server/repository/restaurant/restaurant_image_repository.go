package restaurant

import (
	"context"
	"github.com/BOAZ-LKVK/LKVK-server/server/domain/restaurant"
	"gorm.io/gorm"
)

type RestaurantImageRepository interface {
	FindAllByRestaurantID(ctx context.Context, db *gorm.DB, restaurantID int64) ([]*restaurant.RestaurantImage, error)
	FindAllByRestaurantIDs(ctx context.Context, db *gorm.DB, restaurantIDs []int64) ([]*restaurant.RestaurantImage, error)
}

func NewRestaurantImageRepository() RestaurantImageRepository {
	return &restaurantImageRepository{}
}

type restaurantImageRepository struct{}

func (r *restaurantImageRepository) FindAllByRestaurantID(ctx context.Context, db *gorm.DB, restaurantID int64) ([]*restaurant.RestaurantImage, error) {
	var images []*restaurant.RestaurantImage
	result := db.
		Where(restaurant.RestaurantImage{
			RestaurantID: restaurantID,
		}).
		Find(&images)
	if result.Error != nil {
		return nil, result.Error
	}

	return images, nil
}

func (r *restaurantImageRepository) FindAllByRestaurantIDs(ctx context.Context, db *gorm.DB, restaurantIDs []int64) ([]*restaurant.RestaurantImage, error) {
	var images []*restaurant.RestaurantImage
	result := db.
		Where("restaurant_id IN ?", restaurantIDs).
		Find(&images)
	if result.Error != nil {
		return nil, result.Error
	}

	return images, nil
}
