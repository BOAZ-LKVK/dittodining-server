package recommendation

import (
	"github.com/BOAZ-LKVK/LKVK-server/domain/recommendation"
	"gorm.io/gorm"
)

type RestaurantRecommendationRequestRepository interface {
	Create(request *recommendation.RestaurantRecommendationRequest) (*recommendation.RestaurantRecommendationRequest, error)
	FindByID(restaurantRecommendationRequestID int64) (*recommendation.RestaurantRecommendationRequest, error)
}

func NewRestaurantRecommendationRequestRepository(db *gorm.DB) RestaurantRecommendationRequestRepository {
	return &restaurantRecommendationRequestRepository{db: db}
}

type restaurantRecommendationRequestRepository struct {
	db *gorm.DB
}

func (r *restaurantRecommendationRequestRepository) Create(request *recommendation.RestaurantRecommendationRequest) (*recommendation.RestaurantRecommendationRequest, error) {
	result := r.db.Create(request)
	if result.Error != nil {
		return nil, result.Error
	}

	return request, nil
}

func (r *restaurantRecommendationRequestRepository) FindByID(restaurantRecommendationRequestID int64) (*recommendation.RestaurantRecommendationRequest, error) {
	request := &recommendation.RestaurantRecommendationRequest{
		RestaurantRecommendationRequestID: restaurantRecommendationRequestID,
	}
	result := r.db.Find(request)
	if result.Error != nil {
		return nil, result.Error
	}

	return request, nil
}
