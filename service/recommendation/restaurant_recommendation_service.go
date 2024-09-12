package recommendation

import (
	"encoding/json"
	"fmt"
	recommendation_domain "github.com/BOAZ-LKVK/LKVK-server/domain/recommendation"
	restaurant_domain "github.com/BOAZ-LKVK/LKVK-server/domain/restaurant"
	recommendation_repository "github.com/BOAZ-LKVK/LKVK-server/repository/recommendation"
	restaurant_repository "github.com/BOAZ-LKVK/LKVK-server/repository/restaurant"
	"github.com/BOAZ-LKVK/LKVK-server/service/recommendation/model"
	restaurant_model "github.com/BOAZ-LKVK/LKVK-server/service/restaurant/model"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"strconv"
	"time"
)

type RestaurantRecommendationService interface {
	RequestRestaurantRecommendation(userID *int64, userLocation recommendation_domain.UserLocation, now time.Time) (*model.RequestRestaurantRecommendationResult, error)
	GetRestaurantRecommendationRequest(restaurantRecommendationRequestID int64) (*recommendation_domain.RestaurantRecommendationRequest, error)
	ListRecommendedRestaurants(restaurantRecommendationRequestID int64, cursorRestaurantRecommendationID *int64, limit *int64) (*model.ListRecommendedRestaurantsResult, error)
}

func NewRestaurantRecommendationService(
	restaurantRecommendationRequestRepository recommendation_repository.RestaurantRecommendationRequestRepository,
	restaurantRecommendationRepository recommendation_repository.RestaurantRecommendationRepository,
	restaurantRepository restaurant_repository.RestaurantRepository,
	restaurantMenuRepository restaurant_repository.RestaurantMenuRepository,
	restaurantReviewRepository restaurant_repository.RestaurantReviewRepository,
) RestaurantRecommendationService {
	return &restaurantRecommendationService{
		restaurantRecommendationRequestRepository: restaurantRecommendationRequestRepository,
		restaurantRecommendationRepository:        restaurantRecommendationRepository,
		restaurantRepository:                      restaurantRepository,
		restaurantMenuRepository:                  restaurantMenuRepository,
		restaurantReviewRepository:                restaurantReviewRepository,
	}
}

type restaurantRecommendationService struct {
	restaurantRecommendationRequestRepository recommendation_repository.RestaurantRecommendationRequestRepository
	restaurantRecommendationRepository        recommendation_repository.RestaurantRecommendationRepository
	restaurantRepository                      restaurant_repository.RestaurantRepository
	restaurantMenuRepository                  restaurant_repository.RestaurantMenuRepository
	restaurantReviewRepository                restaurant_repository.RestaurantReviewRepository
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

	// TODO: restaurantRecommendations 생성 로직 추가 - 미리 추천 데이터를 만들고 노출하는 구조

	return &model.RequestRestaurantRecommendationResult{
		RestaurantRecommendationRequestID: created.RestaurantRecommendationRequestID,
	}, nil
}

func (s *restaurantRecommendationService) GetRestaurantRecommendationRequest(restaurantRecommendationRequestID int64) (*recommendation_domain.RestaurantRecommendationRequest, error) {
	return s.restaurantRecommendationRequestRepository.FindByID(restaurantRecommendationRequestID)
}

func (s *restaurantRecommendationService) ListRecommendedRestaurants(restaurantRecommendationRequestID int64, cursorRestaurantRecommendationID *int64, limit *int64) (*model.ListRecommendedRestaurantsResult, error) {
	_, err := s.GetRestaurantRecommendationRequest(restaurantRecommendationRequestID)
	if err != nil {
		return nil, err
	}

	recommendations, err := s.restaurantRecommendationRepository.FindAllByRestaurantRecommendationRequestID(
		restaurantRecommendationRequestID,
		cursorRestaurantRecommendationID,
		limit,
	)
	if err != nil {
		return nil, err
	}

	restaurantIDs := lo.Map(recommendations, func(item *recommendation_domain.RestaurantRecommendation, index int) int64 {
		return item.RestaurantID
	})

	restaurants, err := s.restaurantRepository.FindByIDs(restaurantIDs)
	if err != nil {
		return nil, err
	}
	restaurantByID := lo.SliceToMap(restaurants, func(item *restaurant_domain.Restaurant) (int64, *restaurant_domain.Restaurant) {
		return item.RestaurantID, item
	})

	menus, err := s.restaurantMenuRepository.FindAllByRestaurantIDs(restaurantIDs)
	if err != nil {
		return nil, err
	}
	menusByRestaurantID := lo.GroupBy(menus, func(item *restaurant_domain.RestaurantMenu) int64 {
		return item.RestaurantID
	})
	reviews, err := s.restaurantReviewRepository.FindAllByRestaurantIDs(restaurantIDs)
	if err != nil {
		return nil, err
	}
	reviewsByRestaurantID := lo.GroupBy(reviews, func(item *restaurant_domain.RestaurantReview) int64 {
		return item.RestaurantID
	})

	recommendedRestaurants := make([]*model.RecommendedRestaurant, 0, len(recommendations))
	for _, recommendation := range recommendations {
		r := restaurantByID[recommendation.RestaurantID]
		menuItems := menusByRestaurantID[recommendation.RestaurantID]
		reviewItems := reviewsByRestaurantID[recommendation.RestaurantID]

		// TODO: refactor domain 쪽으로 코드 이전
		var businessHours []*restaurant_domain.BusinessHour
		if err := json.Unmarshal([]byte(r.BusinessHoursJSON), &businessHours); err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal business hours")
		}

		restaurantImageURLs := make([]string, 0)
		for _, image := range r.RestaurantImages {
			restaurantImageURLs = append(restaurantImageURLs, image.ImageURL)
		}

		menuItemModels := make([]*restaurant_model.RestaurantMenu, 0, len(menuItems))
		for _, item := range menuItems {
			menuItemModels = append(menuItemModels, &restaurant_model.RestaurantMenu{
				RestaurantMenuID: item.RestaurantMenuID,
				Name:             item.Name,
				Description:      item.Description,
				Price:            item.Price.IntPart(),
			})
		}

		reviewModels := make([]*restaurant_model.RestaurantReviewItem, 0, len(reviewItems))
		for _, item := range reviewItems {
			reviewModels = append(reviewModels, &restaurant_model.RestaurantReviewItem{
				RestaurantReviewID: item.RestaurantReviewID,
				WriterName:         item.WriterName,
				Score:              item.Score,
				Content:            item.Content,
				WroteAt:            item.WroteAt,
			})
		}

		recommendedRestaurants = append(recommendedRestaurants, &model.RecommendedRestaurant{
			Restaurant: model.RestaurantRecommendation{
				RestaurantID:        recommendation.RestaurantID,
				Name:                r.Name,
				Description:         r.Description,
				PriceRangePerPerson: fmt.Sprintf("%s ~ %s", r.MinimumPricePerPerson.String(), r.MaximumPricePerPerson.String()),
				Distance:            recommendation.DistanceInMeters.String(),
				BusinessHours:       businessHours,
				RestaurantImageURLs: restaurantImageURLs,
			},
			MenuItems: menuItemModels,
			Review: restaurant_model.RestaurantReview{
				Statistics: &restaurant_model.RestaurantReviewStatistics{
					// TODO: refactor restaurant에 review_total_count 추가
					Kakao: &restaurant_model.RestaurantReviewKakaoStatistics{
						AverageScore: r.AverageScoreFromKakao,
						Count:        int64(len(reviewItems)),
					},
					Naver: &restaurant_model.RestaurantReviewNaverStatistics{
						AverageScore: r.AverageScoreFromNaver,
						Count:        int64(len(reviewItems)),
					},
				},
				Reviews: reviewModels,
			},
		})
	}

	var nextCursor *string
	if len(recommendations) > 0 {
		last := recommendations[len(recommendations)-1]
		nextCursor = lo.ToPtr(strconv.FormatInt(last.RestaurantRecommendationID, 10))
	}

	return &model.ListRecommendedRestaurantsResult{
		RecommendedRestaurants: recommendedRestaurants,
		NextCursor:             nextCursor,
	}, nil
}
