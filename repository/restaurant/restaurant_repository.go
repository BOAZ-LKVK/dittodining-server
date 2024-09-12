package restaurant

import (
	"github.com/BOAZ-LKVK/LKVK-server/domain/restaurant"
	"gorm.io/gorm"
)

type RestaurantRepository interface {
	FindByID(restaurantID int64) (*restaurant.Restaurant, error)
	FindByIDs(restaurantIDs []int64) ([]*restaurant.Restaurant, error)
}

func NewRestaurantRepository(db *gorm.DB) RestaurantRepository {
	return &restaurantRepository{db: db}
}

type restaurantRepository struct {
	db *gorm.DB
}

func (r *restaurantRepository) FindByID(restaurantID int64) (*restaurant.Restaurant, error) {
	var existingRestaurant *restaurant.Restaurant
	result := r.db.
		Where(restaurant.Restaurant{
			RestaurantID: restaurantID,
		}).
		Find(&existingRestaurant)
	if result.Error != nil {
		return nil, result.Error
	}

	return existingRestaurant, nil
}

func (r *restaurantRepository) FindByIDs(restaurantIDs []int64) ([]*restaurant.Restaurant, error) {
	var existingRestaurants []*restaurant.Restaurant
	result := r.db.
		Where("restaurant_id IN ?", restaurantIDs).
		Find(&existingRestaurants)
	if result.Error != nil {
		return nil, result.Error
	}

	return existingRestaurants, nil
}
