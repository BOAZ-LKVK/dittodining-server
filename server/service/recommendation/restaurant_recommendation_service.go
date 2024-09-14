package recommendation

import (
	"encoding/json"
	"fmt"
	recommendation_domain "github.com/BOAZ-LKVK/LKVK-server/server/domain/recommendation"
	restaurant_domain "github.com/BOAZ-LKVK/LKVK-server/server/domain/restaurant"
	recommendation_repository "github.com/BOAZ-LKVK/LKVK-server/server/repository/recommendation"
	restaurant_repository "github.com/BOAZ-LKVK/LKVK-server/server/repository/restaurant"
	recommendation_model "github.com/BOAZ-LKVK/LKVK-server/server/service/recommendation/model"
	restaurant_model "github.com/BOAZ-LKVK/LKVK-server/server/service/restaurant/model"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"strconv"
	"time"
)

type RestaurantRecommendationService interface {
	RequestRestaurantRecommendation(userID *int64, userLocation recommendation_domain.UserLocation, now time.Time) (*recommendation_model.RequestRestaurantRecommendationResult, error)
	GetRestaurantRecommendationRequest(restaurantRecommendationRequestID int64) (*recommendation_domain.RestaurantRecommendationRequest, error)
	ListRecommendedRestaurants(restaurantRecommendationRequestID int64, cursorRestaurantRecommendationID *int64, limit *int64) (*recommendation_model.ListRecommendedRestaurantsResult, error)
	SelectRestaurantRecommendation(restaurantRecommendationRequestID int64, restaurantRecommendationIDs []int64) (*recommendation_model.SelectRestaurantRecommendationResult, error)
	GetRestaurantRecommendationResult(restaurantRecommendationRequestID int64) (*recommendation_model.GetRestaurantRecommendationResultResult, error)
}

func NewRestaurantRecommendationService(
	restaurantRecommendationRequestRepository recommendation_repository.RestaurantRecommendationRequestRepository,
	restaurantRecommendationRepository recommendation_repository.RestaurantRecommendationRepository,
	restaurantRepository restaurant_repository.RestaurantRepository,
	restaurantMenuRepository restaurant_repository.RestaurantMenuRepository,
	restaurantReviewRepository restaurant_repository.RestaurantReviewRepository,
	selectedRestaurantRecommendationRepository recommendation_repository.SelectedRestaurantRecommendationRepository,
) RestaurantRecommendationService {
	return &restaurantRecommendationService{
		restaurantRecommendationRequestRepository:  restaurantRecommendationRequestRepository,
		restaurantRecommendationRepository:         restaurantRecommendationRepository,
		restaurantRepository:                       restaurantRepository,
		restaurantMenuRepository:                   restaurantMenuRepository,
		restaurantReviewRepository:                 restaurantReviewRepository,
		selectedRestaurantRecommendationRepository: selectedRestaurantRecommendationRepository,
	}
}

type restaurantRecommendationService struct {
	restaurantRecommendationRequestRepository  recommendation_repository.RestaurantRecommendationRequestRepository
	restaurantRecommendationRepository         recommendation_repository.RestaurantRecommendationRepository
	restaurantRepository                       restaurant_repository.RestaurantRepository
	restaurantMenuRepository                   restaurant_repository.RestaurantMenuRepository
	restaurantReviewRepository                 restaurant_repository.RestaurantReviewRepository
	selectedRestaurantRecommendationRepository recommendation_repository.SelectedRestaurantRecommendationRepository
}

func (s *restaurantRecommendationService) RequestRestaurantRecommendation(userID *int64, userLocation recommendation_domain.UserLocation, now time.Time) (*recommendation_model.RequestRestaurantRecommendationResult, error) {
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

	return &recommendation_model.RequestRestaurantRecommendationResult{
		RestaurantRecommendationRequestID: created.RestaurantRecommendationRequestID,
	}, nil
}

func (s *restaurantRecommendationService) GetRestaurantRecommendationRequest(restaurantRecommendationRequestID int64) (*recommendation_domain.RestaurantRecommendationRequest, error) {
	return s.restaurantRecommendationRequestRepository.FindByID(restaurantRecommendationRequestID)
}

func (s *restaurantRecommendationService) ListRecommendedRestaurants(restaurantRecommendationRequestID int64, cursorRestaurantRecommendationID *int64, limit *int64) (*recommendation_model.ListRecommendedRestaurantsResult, error) {
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

	recommendedRestaurants := make([]*recommendation_model.RecommendedRestaurant, 0, len(recommendations))
	for _, recommendation := range recommendations {
		restaurant := restaurantByID[recommendation.RestaurantID]
		menuItems := menusByRestaurantID[recommendation.RestaurantID]
		reviewItems := reviewsByRestaurantID[recommendation.RestaurantID]

		// TODO: refactor review count는 조회 성능을 위해 restaurant에 저장하도록
		totalCount, err := s.restaurantReviewRepository.CountByRestaurantID(recommendation.RestaurantID)
		if err != nil {
			return nil, err
		}

		recommendedRestaurantModel, err := makeRecommendedRestaurantModel(recommendation, restaurant, menuItems, reviewItems, totalCount)
		if err != nil {
			return nil, err
		}

		recommendedRestaurants = append(recommendedRestaurants, recommendedRestaurantModel)
	}

	var nextCursor *string
	if len(recommendations) > 0 {
		last := recommendations[len(recommendations)-1]
		nextCursor = lo.ToPtr(strconv.FormatInt(last.RestaurantRecommendationID, 10))
	}

	return &recommendation_model.ListRecommendedRestaurantsResult{
		RecommendedRestaurants: recommendedRestaurants,
		NextCursor:             nextCursor,
	}, nil
}

func (s *restaurantRecommendationService) SelectRestaurantRecommendation(restaurantRecommendationRequestID int64, restaurantRecommendationIDs []int64) (*recommendation_model.SelectRestaurantRecommendationResult, error) {
	request, err := s.GetRestaurantRecommendationRequest(restaurantRecommendationRequestID)
	if err != nil {
		return nil, err
	}

	recommendations, err := s.restaurantRecommendationRepository.FindAllByIDs(restaurantRecommendationIDs)
	if err != nil {
		return nil, err
	}
	if len(recommendations) != len(restaurantRecommendationIDs) {
		return nil, errors.New("not exist restaurantRecommendationID")
	}

	selectedRestaurantRecommendations := make([]*recommendation_domain.SelectedRestaurantRecommendation, 0, len(recommendations))
	for _, r := range recommendations {
		selectedRestaurantRecommendations = append(selectedRestaurantRecommendations, recommendation_domain.NewSelectedRestaurantRecommendation(
			request.RestaurantRecommendationRequestID,
			r.RestaurantRecommendationID,
			r.RestaurantID,
		))
	}

	if err := s.selectedRestaurantRecommendationRepository.SaveAll(selectedRestaurantRecommendations); err != nil {
		return nil, err
	}

	return &recommendation_model.SelectRestaurantRecommendationResult{}, nil
}

func (s *restaurantRecommendationService) GetRestaurantRecommendationResult(restaurantRecommendationRequestID int64) (*recommendation_model.GetRestaurantRecommendationResultResult, error) {
	selectedRestaurantRecommendations, err := s.selectedRestaurantRecommendationRepository.FindAllByRestaurantRecommendationRequestID(restaurantRecommendationRequestID)
	if err != nil {
		return nil, err
	}

	restaurantIDs := lo.Map(selectedRestaurantRecommendations, func(item *recommendation_domain.SelectedRestaurantRecommendation, index int) int64 {
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

	restaurantRecommendationIDs := lo.Map(selectedRestaurantRecommendations, func(item *recommendation_domain.SelectedRestaurantRecommendation, index int) int64 {
		return item.RestaurantRecommendationID
	})

	restaurantRecommendations, err := s.restaurantRecommendationRepository.FindAllByIDs(restaurantRecommendationIDs)
	if err != nil {
		return nil, err
	}

	restaurantRecommendationByID := lo.SliceToMap(restaurantRecommendations, func(item *recommendation_domain.RestaurantRecommendation) (int64, *recommendation_domain.RestaurantRecommendation) {
		return item.RestaurantRecommendationID, item
	})

	results := make([]*recommendation_model.RestaurantRecommendationResult, 0, len(selectedRestaurantRecommendations))
	for _, r := range selectedRestaurantRecommendations {
		restaurant := restaurantByID[r.RestaurantID]
		menuItems := menusByRestaurantID[r.RestaurantID]
		reviewItems := reviewsByRestaurantID[r.RestaurantID]
		recommendation := restaurantRecommendationByID[r.RestaurantRecommendationID]

		// TODO: refactor review count는 조회 성능을 위해 restaurant에 저장하도록
		totalCount, err := s.restaurantReviewRepository.CountByRestaurantID(r.RestaurantID)
		if err != nil {
			return nil, err
		}

		recommendedRestaurantModel, err := makeRecommendedRestaurantModel(
			recommendation,
			restaurant,
			menuItems,
			reviewItems,
			totalCount,
		)
		if err != nil {
			return nil, err
		}

		results = append(results, &recommendation_model.RestaurantRecommendationResult{
			RestaurantRecommendationID: r.RestaurantRecommendationID,
			Restaurant:                 recommendedRestaurantModel,
		})
	}

	return &recommendation_model.GetRestaurantRecommendationResultResult{
		Results: results,
	}, nil
}

// TODO: refactor domain이나 model 쪽으로 코드 이전
func makeRecommendedRestaurantModel(
	recommendation *recommendation_domain.RestaurantRecommendation,
	restaurant *restaurant_domain.Restaurant,
	menuItems []*restaurant_domain.RestaurantMenu,
	reviewItems []*restaurant_domain.RestaurantReview,
	totalCount int64,
) (*recommendation_model.RecommendedRestaurant, error) {
	var businessHours []*restaurant_domain.BusinessHour
	if err := json.Unmarshal([]byte(restaurant.BusinessHoursJSON), &businessHours); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal business hours")
	}

	restaurantImageURLs := make([]string, 0)
	for _, image := range restaurant.RestaurantImages {
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

	return &recommendation_model.RecommendedRestaurant{
		Restaurant: recommendation_model.RestaurantRecommendation{
			RestaurantID:        recommendation.RestaurantID,
			Name:                restaurant.Name,
			Description:         restaurant.Description,
			PriceRangePerPerson: fmt.Sprintf("%s ~ %s", restaurant.MinimumPricePerPerson.String(), restaurant.MaximumPricePerPerson.String()),
			Distance:            recommendation.DistanceInMeters.String(),
			BusinessHours:       businessHours,
			RestaurantImageURLs: restaurantImageURLs,
		},
		MenuItems: menuItemModels,
		Review: restaurant_model.RestaurantReview{
			Statistics: &restaurant_model.RestaurantReviewStatistics{
				Kakao: &restaurant_model.RestaurantReviewKakaoStatistics{
					AverageScore: restaurant.AverageScoreFromKakao,
					Count:        int64(len(reviewItems)),
				},
				Naver: &restaurant_model.RestaurantReviewNaverStatistics{
					AverageScore: restaurant.AverageScoreFromNaver,
					Count:        int64(len(reviewItems)),
				},
			},
			Reviews:    reviewModels,
			TotalCount: totalCount,
		},
	}, nil
}
