package restaurant

import (
	"encoding/json"
	restaurant_domain "github.com/BOAZ-LKVK/LKVK-server/server/domain/restaurant"
	"github.com/BOAZ-LKVK/LKVK-server/server/repository/restaurant"
	restaurant_model "github.com/BOAZ-LKVK/LKVK-server/server/service/restaurant/model"
	"github.com/pkg/errors"
)

type RestaurantService interface {
	GetRestaurant(restaurantID int64) (*restaurant_model.Restaurant, error)
	ListRestaurantMenus(restaurantID int64) ([]*restaurant_model.RestaurantMenu, error)
	GetRestaurantReview(restaurantID int64) (*restaurant_model.RestaurantReview, error)
}

func NewRestaurantService(
	restaurantRepository restaurant.RestaurantRepository,
	restaurantMenuRepository restaurant.RestaurantMenuRepository,
	restaurantReviewRepository restaurant.RestaurantReviewRepository,
	restaurantImageRepository restaurant.RestaurantImageRepository,
) RestaurantService {
	return &restaurantService{
		restaurantRepository:       restaurantRepository,
		restaurantMenuRepository:   restaurantMenuRepository,
		restaurantReviewRepository: restaurantReviewRepository,
		restaurantImageRepository:  restaurantImageRepository,
	}
}

type restaurantService struct {
	restaurantRepository       restaurant.RestaurantRepository
	restaurantMenuRepository   restaurant.RestaurantMenuRepository
	restaurantReviewRepository restaurant.RestaurantReviewRepository
	restaurantImageRepository  restaurant.RestaurantImageRepository
}

func (s *restaurantService) GetRestaurant(restaurantID int64) (*restaurant_model.Restaurant, error) {
	r, err := s.restaurantRepository.FindByID(restaurantID)
	if err != nil {
		return nil, err
	}

	var businessHours []*restaurant_domain.BusinessHour
	if err := json.Unmarshal([]byte(r.BusinessHoursJSON), &businessHours); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal business hours")
	}

	images, err := s.restaurantImageRepository.FindAllByRestaurantID(restaurantID)
	if err != nil {
		return nil, err
	}

	restaurantImageURLs := make([]string, 0, len(images))
	for _, image := range images {
		restaurantImageURLs = append(restaurantImageURLs, image.ImageURL)
	}

	return &restaurant_model.Restaurant{
		RestaurantID:          r.RestaurantID,
		Name:                  r.Name,
		Address:               r.Address,
		Description:           r.Description,
		MaximumPricePerPerson: r.MaximumPricePerPerson,
		MinimumPricePerPerson: r.MinimumPricePerPerson,
		Longitude:             r.Longitude,
		Latitude:              r.Latitude,
		BusinessHours:         businessHours,
		RestaurantImageURLs:   restaurantImageURLs,
	}, nil
}

func (s *restaurantService) ListRestaurantMenus(restaurantID int64) ([]*restaurant_model.RestaurantMenu, error) {
	menus, err := s.restaurantMenuRepository.FindAllByRestaurantID(restaurantID)
	if err != nil {
		return nil, err
	}

	var modelMenus []*restaurant_model.RestaurantMenu
	for _, menu := range menus {
		modelMenus = append(modelMenus, &restaurant_model.RestaurantMenu{
			RestaurantMenuID: menu.RestaurantMenuID,
			ImageURL:         menu.ImageURL,
			Name:             menu.Name,
			Price:            menu.Price,
			Description:      nil,
		})
	}

	return modelMenus, nil
}

func (s *restaurantService) GetRestaurantReview(restaurantID int64) (*restaurant_model.RestaurantReview, error) {
	r, err := s.restaurantRepository.FindByID(restaurantID)
	if err != nil {
		return nil, err
	}

	reviews, err := s.restaurantReviewRepository.FindAllByRestaurantID(restaurantID)
	if err != nil {
		return nil, err
	}

	reviewModels := make([]*restaurant_model.RestaurantReviewItem, 0, len(reviews))
	for _, item := range reviews {
		reviewModels = append(reviewModels, &restaurant_model.RestaurantReviewItem{
			RestaurantReviewID: item.RestaurantReviewID,
			WriterName:         item.WriterName,
			Score:              item.Score,
			Content:            item.Content,
			WroteAt:            item.WroteAt,
		})
	}

	return &restaurant_model.RestaurantReview{
		Statistics: &restaurant_model.RestaurantReviewStatistics{
			Kakao: &restaurant_model.RestaurantReviewKakaoStatistics{
				AverageScore: r.AverageScoreFromKakao,
				Count:        r.TotalReviewCountFromKakao,
			},
			Naver: &restaurant_model.RestaurantReviewNaverStatistics{
				AverageScore: r.AverageScoreFromNaver,
				Count:        r.TotalReviewCountFromNaver,
			},
		},
		Reviews:    reviewModels,
		TotalCount: r.TotalReviewCount,
	}, nil
}
