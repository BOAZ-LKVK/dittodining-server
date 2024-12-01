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
	"golang.org/x/net/context"
	"gorm.io/gorm"
	"strconv"
	"time"
)

const minutesToFindNearbyRestaurants = 15

type RestaurantRecommendationService interface {
	RequestRestaurantRecommendation(ctx context.Context, userID *int64, userLocation recommendation_domain.UserLocation, now time.Time) (*recommendation_model.RequestRestaurantRecommendationResult, error)
	GetRestaurantRecommendationRequest(ctx context.Context, restaurantRecommendationRequestID int64) (*recommendation_domain.RestaurantRecommendationRequest, error)
	ListRecommendedRestaurants(ctx context.Context, restaurantRecommendationRequestID int64, cursorRestaurantRecommendationID *int64, limit *int64) (*recommendation_model.ListRecommendedRestaurantsResult, error)
	SelectRestaurantRecommendation(ctx context.Context, restaurantRecommendationRequestID int64, restaurantRecommendationIDs []int64) (*recommendation_model.SelectRestaurantRecommendationResult, error)
	GetRestaurantRecommendationResult(ctx context.Context, restaurantRecommendationRequestID int64) (*recommendation_model.GetRestaurantRecommendationResultResult, error)
	GetRestaurantRecommendation(ctx context.Context, restaurantRecommendationID int64) (*recommendation_model.GetRestaurantRecommendationResult, error)
}

func NewRestaurantRecommendationService(
	restaurantRecommendationRequestRepository recommendation_repository.RestaurantRecommendationRequestRepository,
	restaurantRecommendationRepository recommendation_repository.RestaurantRecommendationRepository,
	selectedRestaurantRecommendationRepository recommendation_repository.SelectedRestaurantRecommendationRepository,
	restaurantRepository restaurant_repository.RestaurantRepository,
	restaurantMenuRepository restaurant_repository.RestaurantMenuRepository,
	restaurantReviewRepository restaurant_repository.RestaurantReviewRepository,
	restaurantImageRepository restaurant_repository.RestaurantImageRepository,
	db *gorm.DB,
) RestaurantRecommendationService {
	return &restaurantRecommendationService{
		restaurantRecommendationRequestRepository:  restaurantRecommendationRequestRepository,
		restaurantRecommendationRepository:         restaurantRecommendationRepository,
		selectedRestaurantRecommendationRepository: selectedRestaurantRecommendationRepository,
		restaurantRepository:                       restaurantRepository,
		restaurantMenuRepository:                   restaurantMenuRepository,
		restaurantReviewRepository:                 restaurantReviewRepository,
		restaurantImageRepository:                  restaurantImageRepository,
		db:                                         db,
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
	db                                         *gorm.DB
}

func (s *restaurantRecommendationService) RequestRestaurantRecommendation(ctx context.Context, userID *int64, userLocation recommendation_domain.UserLocation, now time.Time) (*recommendation_model.RequestRestaurantRecommendationResult, error) {
	var createdRecommendationRequest *recommendation_domain.RestaurantRecommendationRequest
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		recommendationRequest := recommendation_domain.NewRestaurantRecommendationRequest(
			userID,
			recommendation_domain.NewUserLocation(
				userLocation.Latitude, userLocation.Longitude,
			),
			// TODO: testablity를 위해 clock interface 개발 후 대체
			now,
		)
		created, err := s.restaurantRecommendationRequestRepository.Save(ctx, s.db, recommendationRequest)
		if err != nil {
			return err
		}
		createdRecommendationRequest = created

		if _, err := s.createRecommendations(
			ctx,
			tx,
			userLocation,
			recommendationRequest.RestaurantRecommendationRequestID,
			minutesToFindNearbyRestaurants,
			nil,
			nil,
		); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &recommendation_model.RequestRestaurantRecommendationResult{
		RestaurantRecommendationRequestID: createdRecommendationRequest.RestaurantRecommendationRequestID,
	}, nil
}

func (s *restaurantRecommendationService) createRecommendations(
	ctx context.Context,
	tx *gorm.DB,
	userLocation recommendation_domain.UserLocation,
	restaurantRecommendationRequestID int64,
	minutes int64,
	cursorRecommendationScore *int64,
	limit *int64,
) ([]*recommendation_domain.RestaurantRecommendation, error) {
	restaurantsOrderByRecommendationScoreDesc, err := s.restaurantRepository.FindNearbyAllOrderByRecommendationScoreDesc(
		ctx,
		tx,
		userLocation,
		minutes,
		cursorRecommendationScore,
		limit,
	)
	if err != nil {
		return nil, err
	}

	recommendations := make([]*recommendation_domain.RestaurantRecommendation, 0, len(restaurantsOrderByRecommendationScoreDesc))
	for _, r := range restaurantsOrderByRecommendationScoreDesc {
		distanceInMeters := location.CalculateDistanceInMeters(userLocation.Latitude, userLocation.Longitude, r.Latitude, r.Longitude)

		recommendations = append(recommendations,
			recommendation_domain.NewRestaurantRecommendation(
				restaurantRecommendationRequestID,
				r.RestaurantID,
				distanceInMeters,
			),
		)
	}

	if err := s.restaurantRecommendationRepository.SaveAll(ctx, s.db, recommendations); err != nil {
		return nil, err
	}

	return recommendations, nil
}

func (s *restaurantRecommendationService) GetRestaurantRecommendationRequest(ctx context.Context, restaurantRecommendationRequestID int64) (*recommendation_domain.RestaurantRecommendationRequest, error) {
	return s.restaurantRecommendationRequestRepository.FindByID(ctx, s.db, restaurantRecommendationRequestID)
}

func (s *restaurantRecommendationService) ListRecommendedRestaurants(
	ctx context.Context,
	restaurantRecommendationRequestID int64,
	cursorRestaurantRecommendationID *int64,
	limit *int64,
) (*recommendation_model.ListRecommendedRestaurantsResult, error) {
	recommendationRequest, err := s.GetRestaurantRecommendationRequest(ctx, restaurantRecommendationRequestID)
	if err != nil {
		return nil, err
	}

	recommendations, err := s.restaurantRecommendationRepository.FindAllByRestaurantRecommendationRequestID(
		ctx,
		s.db,
		restaurantRecommendationRequestID,
		cursorRestaurantRecommendationID,
		limit,
	)
	if err != nil {
		return nil, err
	}

	if len(recommendations) == 0 {
		lastRecommendation, err := s.restaurantRecommendationRepository.FindLastOneByRestaurantRecommendationRequestID(
			ctx,
			s.db,
			restaurantRecommendationRequestID,
		)
		if err != nil {
			return nil, err
		}

		lastRecommendedRestaurant, err := s.restaurantRepository.FindByID(ctx, s.db, lastRecommendation.RestaurantID)
		if err != nil {
			return nil, err
		}

		recommendations, err = s.createRecommendations(
			ctx,
			s.db,
			recommendationRequest.UserLocation,
			recommendationRequest.RestaurantRecommendationRequestID,
			minutesToFindNearbyRestaurants,
			lo.ToPtr(lastRecommendedRestaurant.RecommendationScore.IntPart()),
			limit,
		)
		if err != nil {
			return nil, err
		}
	}

	restaurantIDs := lo.Map(recommendations, func(item *recommendation_domain.RestaurantRecommendation, index int) int64 {
		return item.RestaurantID
	})

	restaurants, err := s.restaurantRepository.FindByIDs(ctx, s.db, restaurantIDs)
	if err != nil {
		return nil, err
	}
	restaurantByID := lo.SliceToMap(restaurants, func(item *restaurant_domain.Restaurant) (int64, *restaurant_domain.Restaurant) {
		return item.RestaurantID, item
	})

	restaurantImages, err := s.restaurantImageRepository.FindAllByRestaurantIDs(ctx, s.db, restaurantIDs)
	if err != nil {
		return nil, err
	}
	restaurantImagesByRestaurantID := lo.GroupBy(restaurantImages, func(item *restaurant_domain.RestaurantImage) int64 {
		return item.RestaurantID
	})

	menus, err := s.restaurantMenuRepository.FindAllByRestaurantIDs(ctx, s.db, restaurantIDs)
	if err != nil {
		return nil, err
	}
	menusByRestaurantID := lo.GroupBy(menus, func(item *restaurant_domain.RestaurantMenu) int64 {
		return item.RestaurantID
	})

	reviews, err := s.restaurantReviewRepository.FindAllByRestaurantIDs(ctx, s.db, restaurantIDs)
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
		restaurantImageItems := restaurantImagesByRestaurantID[recommendation.RestaurantID]

		recommendedRestaurantModel, err := makeRecommendedRestaurantModel(recommendation, restaurant, menuItems, reviewItems, restaurantImageItems)
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

func (s *restaurantRecommendationService) SelectRestaurantRecommendation(ctx context.Context, restaurantRecommendationRequestID int64, restaurantRecommendationIDs []int64) (*recommendation_model.SelectRestaurantRecommendationResult, error) {
	request, err := s.GetRestaurantRecommendationRequest(ctx, restaurantRecommendationRequestID)
	if err != nil {
		return nil, err
	}

	recommendations, err := s.restaurantRecommendationRepository.FindAllByIDs(ctx, s.db, restaurantRecommendationIDs)
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

	if err := s.selectedRestaurantRecommendationRepository.SaveAll(ctx, s.db, selectedRestaurantRecommendations); err != nil {
		return nil, err
	}

	return &recommendation_model.SelectRestaurantRecommendationResult{}, nil
}

func (s *restaurantRecommendationService) GetRestaurantRecommendationResult(ctx context.Context, restaurantRecommendationRequestID int64) (*recommendation_model.GetRestaurantRecommendationResultResult, error) {
	selectedRestaurantRecommendations, err := s.selectedRestaurantRecommendationRepository.FindAllByRestaurantRecommendationRequestID(ctx, s.db, restaurantRecommendationRequestID)
	if err != nil {
		return nil, err
	}

	restaurantIDs := lo.Map(selectedRestaurantRecommendations, func(item *recommendation_domain.SelectedRestaurantRecommendation, index int) int64 {
		return item.RestaurantID
	})

	restaurants, err := s.restaurantRepository.FindByIDs(ctx, s.db, restaurantIDs)
	if err != nil {
		return nil, err
	}
	restaurantByID := lo.SliceToMap(restaurants, func(item *restaurant_domain.Restaurant) (int64, *restaurant_domain.Restaurant) {
		return item.RestaurantID, item
	})

	restaurantImages, err := s.restaurantImageRepository.FindAllByRestaurantIDs(ctx, s.db, restaurantIDs)
	if err != nil {
		return nil, err
	}
	restaurantImagesByRestaurantID := lo.GroupBy(restaurantImages, func(item *restaurant_domain.RestaurantImage) int64 {
		return item.RestaurantID
	})

	menus, err := s.restaurantMenuRepository.FindAllByRestaurantIDs(ctx, s.db, restaurantIDs)
	if err != nil {
		return nil, err
	}
	menusByRestaurantID := lo.GroupBy(menus, func(item *restaurant_domain.RestaurantMenu) int64 {
		return item.RestaurantID
	})

	reviews, err := s.restaurantReviewRepository.FindAllByRestaurantIDs(ctx, s.db, restaurantIDs)
	if err != nil {
		return nil, err
	}
	reviewsByRestaurantID := lo.GroupBy(reviews, func(item *restaurant_domain.RestaurantReview) int64 {
		return item.RestaurantID
	})

	restaurantRecommendationIDs := lo.Map(selectedRestaurantRecommendations, func(item *recommendation_domain.SelectedRestaurantRecommendation, index int) int64 {
		return item.RestaurantRecommendationID
	})

	restaurantRecommendations, err := s.restaurantRecommendationRepository.FindAllByIDs(ctx, s.db, restaurantRecommendationIDs)
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

func (s *restaurantRecommendationService) GetRestaurantRecommendation(ctx context.Context, restaurantRecommendationID int64) (*recommendation_model.GetRestaurantRecommendationResult, error) {
	existingRecommendation, err := s.restaurantRecommendationRepository.FindByID(ctx, s.db, restaurantRecommendationID)
	if err != nil {
		return nil, err
	}

	restaurant, err := s.restaurantRepository.FindByID(ctx, s.db, existingRecommendation.RestaurantID)
	if err != nil {
		return nil, err
	}

	menuItems, err := s.restaurantMenuRepository.FindAllByRestaurantID(ctx, s.db, existingRecommendation.RestaurantID)
	if err != nil {
		return nil, err
	}

	reviewItems, err := s.restaurantReviewRepository.FindAllByRestaurantID(ctx, s.db, existingRecommendation.RestaurantID)
	if err != nil {
		return nil, err
	}

	restaurantImageItems, err := s.restaurantImageRepository.FindAllByRestaurantID(ctx, s.db, existingRecommendation.RestaurantID)
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
