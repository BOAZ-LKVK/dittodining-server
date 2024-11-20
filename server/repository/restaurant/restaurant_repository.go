package restaurant

import (
	"context"
	"fmt"
	recommendation_domain "github.com/BOAZ-LKVK/LKVK-server/server/domain/recommendation"
	"github.com/BOAZ-LKVK/LKVK-server/server/domain/restaurant"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

var ErrRestaurantNotFound = errors.New("restaurant not found")

type RestaurantRepository interface {
	FindByID(ctx context.Context, db *gorm.DB, restaurantID int64) (*restaurant.Restaurant, error)
	FindByIDs(ctx context.Context, db *gorm.DB, restaurantIDs []int64) ([]*restaurant.Restaurant, error)
	FindNearbyAllOrderByRecommendationScoreDesc(
		ctx context.Context,
		db *gorm.DB,
		userLocation recommendation_domain.UserLocation,
		minutes int64,
		cursorRecommendationScore *int64,
		limit *int64,
	) ([]*restaurant.Restaurant, error)
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

func (r *restaurantRepository) FindNearbyAllOrderByRecommendationScoreDesc(
	ctx context.Context,
	db *gorm.DB,
	userLocation recommendation_domain.UserLocation,
	minutes int64,
	cursorRecommendationScore *int64,
	limit *int64,
) ([]*restaurant.Restaurant, error) {
	// 이동 반경 계산 (80m/min × 시간)
	radius := 80 * minutes

	// 거리 계산 SQL 표현식
	distanceExpr := fmt.Sprintf(
		"ST_Distance_Sphere(POINT(longitude, latitude), POINT(%f, %f))",
		userLocation.Longitude.InexactFloat64(), userLocation.Latitude.InexactFloat64(),
	)

	// 쿼리 실행
	queryBuilder := db.
		Model(&restaurant.Restaurant{}).
		Select("*, "+distanceExpr+" AS distance").
		Where(distanceExpr+" <= ?", radius) // 반경 조건

	if cursorRecommendationScore != nil {
		queryBuilder.Where("recommendation_score < ?", cursorRecommendationScore)
	}

	if limit != nil {
		queryBuilder.Limit(int(*limit))
	} else {
		// default limit
		queryBuilder.Limit(10)
	}

	var existingRestaurants []*restaurant.Restaurant
	result := queryBuilder.
		Order("recommendation_score DESC, distance ASC").
		Find(&existingRestaurants)
	if result.Error != nil {
		return nil, result.Error
	}

	return existingRestaurants, nil
}
