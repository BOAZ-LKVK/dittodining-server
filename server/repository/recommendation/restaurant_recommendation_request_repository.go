package recommendation

import (
	"errors"
	"github.com/BOAZ-LKVK/LKVK-server/server/domain/recommendation"
	"gorm.io/gorm"
)

var ErrRestaurantRecommendationRequestNotFound = errors.New("restaurant recommendation request not found")

type RestaurantRecommendationRequestRepository interface {
	Save(request *recommendation.RestaurantRecommendationRequest) (*recommendation.RestaurantRecommendationRequest, error)
	FindByID(restaurantRecommendationRequestID int64) (*recommendation.RestaurantRecommendationRequest, error)
}

func NewRestaurantRecommendationRequestRepository(db *gorm.DB) RestaurantRecommendationRequestRepository {
	return &restaurantRecommendationRequestRepository{db: db}
}

type restaurantRecommendationRequestRepository struct {
	db *gorm.DB
}

func (r *restaurantRecommendationRequestRepository) Save(request *recommendation.RestaurantRecommendationRequest) (*recommendation.RestaurantRecommendationRequest, error) {
	result := r.db.Create(request)
	if result.Error != nil {
		return nil, result.Error
	}

	return request, nil
}

func (r *restaurantRecommendationRequestRepository) FindByID(restaurantRecommendationRequestID int64) (*recommendation.RestaurantRecommendationRequest, error) {
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
