package restaurant

import (
	"github.com/BOAZ-LKVK/LKVK-server/server/domain/restaurant"
	"gorm.io/gorm"
)

type RestaurantImageRepository interface {
	FindAllByRestaurantID(restaurantID int64) ([]*restaurant.RestaurantImage, error)
	FindAllByRestaurantIDs(restaurantIDs []int64) ([]*restaurant.RestaurantImage, error)
}

func NewRestaurantImageRepository(db *gorm.DB) RestaurantImageRepository {
	return &restaurantImageRepository{
		db: db,
	}
}

type restaurantImageRepository struct {
	db *gorm.DB
}

func (r *restaurantImageRepository) FindAllByRestaurantID(restaurantID int64) ([]*restaurant.RestaurantImage, error) {
	var images []*restaurant.RestaurantImage
	result := r.db.
		Where(restaurant.RestaurantImage{
			RestaurantID: restaurantID,
		}).
		Find(&images)
	if result.Error != nil {
		return nil, result.Error
	}

	return images, nil
}

func (r *restaurantImageRepository) FindAllByRestaurantIDs(restaurantIDs []int64) ([]*restaurant.RestaurantImage, error) {
	var images []*restaurant.RestaurantImage
	result := r.db.
		Where("restaurant_id IN ?", restaurantIDs).
		Find(&images)
	if result.Error != nil {
		return nil, result.Error
	}

	return images, nil
}
