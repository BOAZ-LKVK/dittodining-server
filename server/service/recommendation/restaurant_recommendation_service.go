package recommendation

import (
	"encoding/json"
	"github.com/BOAZ-LKVK/LKVK-server/pkg/location"
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
	GetRestaurantRecommendation(restaurantRecommendationID int64) (*recommendation_model.GetRestaurantRecommendationResult, error)
}

func NewRestaurantRecommendationService(
	restaurantRecommendationRequestRepository recommendation_repository.RestaurantRecommendationRequestRepository,
	restaurantRecommendationRepository recommendation_repository.RestaurantRecommendationRepository,
	selectedRestaurantRecommendationRepository recommendation_repository.SelectedRestaurantRecommendationRepository,
	restaurantRepository restaurant_repository.RestaurantRepository,
	restaurantMenuRepository restaurant_repository.RestaurantMenuRepository,
	restaurantReviewRepository restaurant_repository.RestaurantReviewRepository,
	restaurantImageRepository restaurant_repository.RestaurantImageRepository,
) RestaurantRecommendationService {
	return &restaurantRecommendationService{
		restaurantRecommendationRequestRepository:  restaurantRecommendationRequestRepository,
		restaurantRecommendationRepository:         restaurantRecommendationRepository,
		selectedRestaurantRecommendationRepository: selectedRestaurantRecommendationRepository,
		restaurantRepository:                       restaurantRepository,
		restaurantMenuRepository:                   restaurantMenuRepository,
		restaurantReviewRepository:                 restaurantReviewRepository,
		restaurantImageRepository:                  restaurantImageRepository,
	}
}

type restaurantRecommendationService struct {
	restaurantRecommendationRequestRepository  recommendation_repository.RestaurantRecommendationRequestRepository
	restaurantRecommendationRepository         recommendation_repository.RestaurantRecommendationRepository
	selectedRestaurantRecommendationRepository recommendation_repository.SelectedRestaurantRecommendationRepository
	restaurantRepository                       restaurant_repository.RestaurantRepository
	restaurantMenuRepository                   restaurant_repository.RestaurantMenuRepository
	restaurantReviewRepository                 restaurant_repository.RestaurantReviewRepository
	restaurantImageRepository                  restaurant_repository.RestaurantImageRepository
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
	created, err := s.restaurantRecommendationRequestRepository.Save(recommendationRequest)
	if err != nil {
		return nil, errors.New("Cannot Save recommendationRequest to restaurantRecommendationRequestRepository")
	}

	restaurantsOrderByRecommendationScoreDesc, err := s.restaurantRepository.FindAllOrderByRecommendationScoreDesc(10)
	if err != nil {
		return nil, errors.New("Cannot get descend order on RecommendationScore")
	}

	recommendations := make([]*recommendation_domain.RestaurantRecommendation, 0, len(restaurantsOrderByRecommendationScoreDesc))
	for _, r := range restaurantsOrderByRecommendationScoreDesc {
		distanceInMeters := location.CalculateDistanceInMeters(userLocation.Latitude, userLocation.Longitude, r.Latitude, r.Longitude)

		recommendations = append(recommendations,
			recommendation_domain.NewRestaurantRecommendation(
				recommendationRequest.RestaurantRecommendationRequestID,
				r.RestaurantID,
				distanceInMeters,
			),
		)
	}

	if err := s.restaurantRecommendationRepository.SaveAll(recommendations); err != nil {
		return nil, errors.New("Cannot Save recommendations to restaurantRecommendationRepository")
	}

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
		return nil, errors.New("Cannot get RestaurantRecommendationRequest from restaurantRecommendationRequestID")
	}

	recommendations, err := s.restaurantRecommendationRepository.FindAllByRestaurantRecommendationRequestID(
		restaurantRecommendationRequestID,
		cursorRestaurantRecommendationID,
		limit,
	)
	if err != nil {
		return nil, errors.New("Cannot found All of recommendations from restaurantRecommendationRequestID")
	}

	// TODO: if recommendations == nil -> 추천 데이터 더 추가하도록

	restaurantIDs := lo.Map(recommendations, func(item *recommendation_domain.RestaurantRecommendation, index int) int64 {
		return item.RestaurantID
	})

	restaurants, err := s.restaurantRepository.FindByIDs(restaurantIDs)
	if err != nil {
		return nil, errors.New("Cannot found restaurants from restaurantID list")
	}
	restaurantByID := lo.SliceToMap(restaurants, func(item *restaurant_domain.Restaurant) (int64, *restaurant_domain.Restaurant) {
		return item.RestaurantID, item
	})

	restaurantImages, err := s.restaurantImageRepository.FindAllByRestaurantIDs(restaurantIDs)
	if err != nil {
		return nil, err
	}
	restaurantImagesByRestaurantID := lo.GroupBy(restaurantImages, func(item *restaurant_domain.RestaurantImage) int64 {
		return item.RestaurantID
	})

	menus, err := s.restaurantMenuRepository.FindAllByRestaurantIDs(restaurantIDs)
	if err != nil {
		return nil, errors.New("Cannot found restaurant Menu from restaurantID list")
	}
	menusByRestaurantID := lo.GroupBy(menus, func(item *restaurant_domain.RestaurantMenu) int64 {
		return item.RestaurantID
	})

	reviews, err := s.restaurantReviewRepository.FindAllByRestaurantIDs(restaurantIDs)
	if err != nil {
		return nil, errors.New("Cannot found restaurant Review from restaurantID list")
	}
	reviewsByRestaurantID := lo.GroupBy(reviews, func(item *restaurant_domain.RestaurantReview) int64 {
		return item.RestaurantID
	})

	recommendedRestaurants := make([]*recommendation_model.RecommendedRestaurant, 0, len(recommendations))
	for _, recommendation := range recommendations {
		restaurant := restaurantByID[recommendation.RestaurantID]
		menuItems := menusByRestaurantID[recommendation.RestaurantID]
		reviewItems := reviewsByRestaurantID[recommendation.RestaurantID]
		restaurantImageItems := restaurantImagesByRestaurantID[recommendation.RestaurantID]

		recommendedRestaurantModel, err := makeRecommendedRestaurantModel(recommendation, restaurant, menuItems, reviewItems, restaurantImageItems)
		if err != nil {
			return nil, errors.New("Cannot made recommendedRestaurantModel")
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
		return nil, errors.New("not exist restaurantRecommendationId")
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

	restaurantImages, err := s.restaurantImageRepository.FindAllByRestaurantIDs(restaurantIDs)
	if err != nil {
		return nil, err
	}
	restaurantImagesByRestaurantID := lo.GroupBy(restaurantImages, func(item *restaurant_domain.RestaurantImage) int64 {
		return item.RestaurantID
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
		restaurantImageItems := restaurantImagesByRestaurantID[recommendation.RestaurantID]

		recommendedRestaurantModel, err := makeRecommendedRestaurantModel(
			recommendation,
			restaurant,
			menuItems,
			reviewItems,
			restaurantImageItems,
		)
		if err != nil {
			return nil, err
		}

		results = append(results, &recommendation_model.RestaurantRecommendationResult{
			RestaurantRecommendationID: r.RestaurantRecommendationID,
			RecommendedRestaurant:      recommendedRestaurantModel,
		})
	}

	return &recommendation_model.GetRestaurantRecommendationResultResult{
		Results: results,
	}, nil
}

func (s *restaurantRecommendationService) GetRestaurantRecommendation(restaurantRecommendationID int64) (*recommendation_model.GetRestaurantRecommendationResult, error) {
	existingRecommendation, err := s.restaurantRecommendationRepository.FindByID(restaurantRecommendationID)
	if err != nil {
		return nil, err
	}

	restaurant, err := s.restaurantRepository.FindByID(existingRecommendation.RestaurantID)
	if err != nil {
		return nil, err
	}

	menuItems, err := s.restaurantMenuRepository.FindAllByRestaurantID(existingRecommendation.RestaurantID)
	if err != nil {
		return nil, err
	}

	reviewItems, err := s.restaurantReviewRepository.FindAllByRestaurantID(existingRecommendation.RestaurantID)
	if err != nil {
		return nil, err
	}

	restaurantImageItems, err := s.restaurantImageRepository.FindAllByRestaurantID(existingRecommendation.RestaurantID)
	if err != nil {
		return nil, err
	}

	recommendedRestaurantModel, err := makeRecommendedRestaurantModel(existingRecommendation, restaurant, menuItems, reviewItems, restaurantImageItems)
	if err != nil {
		return nil, err
	}

	return &recommendation_model.GetRestaurantRecommendationResult{
		RecommendedRestaurant: recommendedRestaurantModel,
	}, nil
}

// TODO: refactor domain이나 model 쪽으로 코드 이전
func makeRecommendedRestaurantModel(
	recommendation *recommendation_domain.RestaurantRecommendation,
	restaurant *restaurant_domain.Restaurant,
	menuItems []*restaurant_domain.RestaurantMenu,
	reviewItems []*restaurant_domain.RestaurantReview,
	restaurantImages []*restaurant_domain.RestaurantImage,
) (*recommendation_model.RecommendedRestaurant, error) {
	var businessHours []*restaurant_domain.BusinessHour
	if err := json.Unmarshal([]byte(restaurant.BusinessHoursJSON), &businessHours); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal business hours")
	}

	restaurantImageURLs := make([]string, 0, len(restaurantImages))
	for _, image := range restaurantImages {
		restaurantImageURLs = append(restaurantImageURLs, image.ImageURL)
	}

	menuItemModels := make([]*restaurant_model.RestaurantMenu, 0, len(menuItems))
	for _, item := range menuItems {
		menuItemModels = append(menuItemModels, &restaurant_model.RestaurantMenu{
			RestaurantMenuID: item.RestaurantMenuID,
			Name:             item.Name,
			Description:      item.Description,
			Price:            item.Price,
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
			RestaurantRecommendationID: recommendation.RestaurantRecommendationID,
			RestaurantID:               recommendation.RestaurantID,
			Name:                       restaurant.Name,
			Description:                restaurant.Description,
			Location: recommendation_model.Location{
				Latitude:  restaurant.Latitude,
				Longitude: restaurant.Longitude,
			},
			MinimumPricePerPerson: restaurant.MinimumPricePerPerson,
			MaximumPricePerPerson: restaurant.MaximumPricePerPerson,
			DistanceInMeters:      recommendation.DistanceInMeters,
			BusinessHours:         businessHours,
			RestaurantImageURLs:   restaurantImageURLs,
		},
		MenuItems: menuItemModels,
		Review: restaurant_model.RestaurantReview{
			Statistics: &restaurant_model.RestaurantReviewStatistics{
				Kakao: &restaurant_model.RestaurantReviewKakaoStatistics{
					AverageScore: restaurant.AverageScoreFromKakao,
					Count:        restaurant.TotalReviewCountFromKakao,
				},
				Naver: &restaurant_model.RestaurantReviewNaverStatistics{
					AverageScore: restaurant.AverageScoreFromNaver,
					Count:        restaurant.TotalReviewCountFromNaver,
				},
			},
			Reviews:    reviewModels,
			TotalCount: restaurant.TotalReviewCount,
		},
	}, nil
}
