package restaurant

import (
	"github.com/BOAZ-LKVK/LKVK-server/server/domain/restaurant"
	"gorm.io/gorm"
)

type RestaurantMenuRepository interface {
	FindAllByRestaurantID(restaurantID int64) ([]*restaurant.RestaurantMenu, error)
	FindAllByRestaurantIDs(restaurantIDs []int64) ([]*restaurant.RestaurantMenu, error)
}

func NewRestaurantMenuRepository(db *gorm.DB) RestaurantMenuRepository {
	return &restaurantMenuRepository{db: db}
}

type restaurantMenuRepository struct {
	db *gorm.DB
}

func (r *restaurantMenuRepository) FindAllByRestaurantID(restaurantID int64) ([]*restaurant.RestaurantMenu, error) {
	var menus []*restaurant.RestaurantMenu
	result := r.db.
		Where(restaurant.RestaurantMenu{
			RestaurantID: restaurantID,
		}).
		Find(&menus)
	if result.Error != nil {
		return nil, result.Error
	}

	return menus, nil
}

func (r *restaurantMenuRepository) FindAllByRestaurantIDs(restaurantIDs []int64) ([]*restaurant.RestaurantMenu, error) {
	var menus []*restaurant.RestaurantMenu

	result := r.db.
		Where("restaurant_id IN ?", restaurantIDs).
		Find(&menus)
	if result.Error != nil {
		return nil, result.Error
	}

	return menus, nil
}
