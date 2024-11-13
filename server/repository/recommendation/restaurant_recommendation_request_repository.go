package recommendation

import (
	"context"
	"errors"
	"github.com/BOAZ-LKVK/LKVK-server/server/domain/recommendation"
	"gorm.io/gorm"
)

var ErrRestaurantRecommendationRequestNotFound = errors.New("restaurant recommendation request not found")

type RestaurantRecommendationRequestRepository interface {
	Save(ctx context.Context, db *gorm.DB, request *recommendation.RestaurantRecommendationRequest) (*recommendation.RestaurantRecommendationRequest, error)
	FindByID(ctx context.Context, db *gorm.DB, restaurantRecommendationRequestID int64) (*recommendation.RestaurantRecommendationRequest, error)
}

func NewRestaurantRecommendationRequestRepository() RestaurantRecommendationRequestRepository {
	return &restaurantRecommendationRequestRepository{}
}

type restaurantRecommendationRequestRepository struct {
	db *gorm.DB
}

func (r *restaurantRecommendationRequestRepository) Save(ctx context.Context, db *gorm.DB, request *recommendation.RestaurantRecommendationRequest) (*recommendation.RestaurantRecommendationRequest, error) {
	result := db.Create(request)
	if result.Error != nil {
		return nil, result.Error
	}

	return request, nil
}

func (r *restaurantRecommendationRequestRepository) FindByID(ctx context.Context, db *gorm.DB, restaurantRecommendationRequestID int64) (*recommendation.RestaurantRecommendationRequest, error) {
	var request *recommendation.RestaurantRecommendationRequest

	result := r.db.Where(recommendation.RestaurantRecommendationRequest{
		RestaurantRecommendationRequestID: restaurantRecommendationRequestID,
	}).Find(&request)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrRestaurantRecommendationRequestNotFound
		}

		return nil, result.Error
	}

	return request, nil
}
