package restaurant

import (
	"context"
	"github.com/BOAZ-LKVK/LKVK-server/server/domain/restaurant"
	"gorm.io/gorm"
)

type RestaurantMenuRepository interface {
	FindAllByRestaurantID(ctx context.Context, db *gorm.DB, restaurantID int64) ([]*restaurant.RestaurantMenu, error)
	FindAllByRestaurantIDs(ctx context.Context, db *gorm.DB, restaurantIDs []int64) ([]*restaurant.RestaurantMenu, error)
}

func NewRestaurantMenuRepository() RestaurantMenuRepository {
	return &restaurantMenuRepository{}
}

type restaurantMenuRepository struct{}

func (r *restaurantMenuRepository) FindAllByRestaurantID(ctx context.Context, db *gorm.DB, restaurantID int64) ([]*restaurant.RestaurantMenu, error) {
	var menus []*restaurant.RestaurantMenu
	result := db.
		Where(restaurant.RestaurantMenu{
			RestaurantID: restaurantID,
		}).
		Find(&menus)
	if result.Error != nil {
		return nil, result.Error
	}

	return menus, nil
}

func (r *restaurantMenuRepository) FindAllByRestaurantIDs(ctx context.Context, db *gorm.DB, restaurantIDs []int64) ([]*restaurant.RestaurantMenu, error) {
	var menus []*restaurant.RestaurantMenu

	result := db.
		Where("restaurant_id IN ?", restaurantIDs).
		Find(&menus)
	if result.Error != nil {
		return nil, result.Error
	}

	return menus, nil
}
