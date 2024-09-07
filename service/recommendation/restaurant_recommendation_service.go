package recommendation

import (
	recommendation_domain "github.com/BOAZ-LKVK/LKVK-server/domain/recommendation"
	recommendation_repository "github.com/BOAZ-LKVK/LKVK-server/repository/recommendation"
	"github.com/BOAZ-LKVK/LKVK-server/service/recommendation/model"
	"time"
)

type RestaurantRecommendationService interface {
	RequestRestaurantRecommendation(userID *int64, userLocation recommendation_domain.UserLocation, now time.Time) (*model.RequestRestaurantRecommendationResult, error)
}

func NewRestaurantRecommendationService(restaurantRecommendationRequestRepository recommendation_repository.RestaurantRecommendationRequestRepository) RestaurantRecommendationService {
	return &restaurantRecommendationService{restaurantRecommendationRequestRepository: restaurantRecommendationRequestRepository}
}

type restaurantRecommendationService struct {
	restaurantRecommendationRequestRepository recommendation_repository.RestaurantRecommendationRequestRepository
}

func (s *restaurantRecommendationService) RequestRestaurantRecommendation(userID *int64, userLocation recommendation_domain.UserLocation, now time.Time) (*model.RequestRestaurantRecommendationResult, error) {
	recommendationRequest := recommendation_domain.NewRestaurantRecommendationRequest(
		userID,
		recommendation_domain.NewUserLocation(
			userLocation.Latitude, userLocation.Longitude,
		),
		// TODO: testablity를 위해 clock interface 개발 후 대체
		now,
	)
	created, err := s.restaurantRecommendationRequestRepository.Create(recommendationRequest)
	if err != nil {
		return nil, err
	}

	return &model.RequestRestaurantRecommendationResult{
		RestaurantRecommendationRequestID: created.RestaurantRecommendationRequestID,
	}, nil
}
