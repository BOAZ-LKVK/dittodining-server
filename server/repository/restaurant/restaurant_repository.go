package restaurant

import (
	"context"
	"github.com/BOAZ-LKVK/LKVK-server/server/domain/restaurant"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

var ErrRestaurantNotFound = errors.New("restaurant not found")

type RestaurantRepository interface {
	FindByID(ctx context.Context, db *gorm.DB, restaurantID int64) (*restaurant.Restaurant, error)
	FindByIDs(ctx context.Context, db *gorm.DB, restaurantIDs []int64) ([]*restaurant.Restaurant, error)
	FindAllOrderByRecommendationScoreDesc(ctx context.Context, db *gorm.DB, limit int) ([]*restaurant.Restaurant, error)
}

func NewRestaurantRepository() RestaurantRepository {
	return &restaurantRepository{}
}

type restaurantRepository struct{}

func (r *restaurantRepository) FindByID(ctx context.Context, db *gorm.DB, restaurantID int64) (*restaurant.Restaurant, error) {
	var existingRestaurant *restaurant.Restaurant
	result := db.
		Where(restaurant.Restaurant{
			RestaurantID: restaurantID,
		}).
		Find(&existingRestaurant)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrRestaurantNotFound
		}

		return nil, result.Error
	}

	return existingRestaurant, nil
}

func (r *restaurantRepository) FindByIDs(ctx context.Context, db *gorm.DB, restaurantIDs []int64) ([]*restaurant.Restaurant, error) {
	var existingRestaurants []*restaurant.Restaurant
	result := db.
		Where("restaurant_id IN ?", restaurantIDs).
		Find(&existingRestaurants)
	if result.Error != nil {
		return nil, result.Error
	}

	return existingRestaurants, nil
}

func (r *restaurantRepository) FindAllOrderByRecommendationScoreDesc(ctx context.Context, db *gorm.DB, limit int) ([]*restaurant.Restaurant, error) {
	var existingRestaurants []*restaurant.Restaurant
	result := db.
		Order("recommendation_score DESC").
		Limit(limit).
		Find(&existingRestaurants)
	if result.Error != nil {
		return nil, result.Error
	}

	return existingRestaurants, nil
}
