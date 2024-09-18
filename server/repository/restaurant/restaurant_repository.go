package restaurant

import (
	"github.com/BOAZ-LKVK/LKVK-server/server/domain/restaurant"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

var ErrRestaurantNotFound = errors.New("restaurant not found")

type RestaurantRepository interface {
	FindByID(restaurantID int64) (*restaurant.Restaurant, error)
	FindByIDs(restaurantIDs []int64) ([]*restaurant.Restaurant, error)
	FindAllOrderByRecommendationScoreDesc(limit int) ([]*restaurant.Restaurant, error)
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
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrRestaurantNotFound
		}

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

func (r *restaurantRepository) FindAllOrderByRecommendationScoreDesc(limit int) ([]*restaurant.Restaurant, error) {
	var existingRestaurants []*restaurant.Restaurant
	result := r.db.
		Order("recommendation_score DESC").
		Limit(limit).
		Find(&existingRestaurants)
	if result.Error != nil {
		return nil, result.Error
	}

	return existingRestaurants, nil
}
