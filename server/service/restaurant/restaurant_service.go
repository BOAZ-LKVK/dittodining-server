package restaurant

import (
	"github.com/BOAZ-LKVK/LKVK-server/server/repository/restaurant"
	model2 "github.com/BOAZ-LKVK/LKVK-server/server/service/restaurant/model"
)

type RestaurantService interface {
	GetRestaurant(restaurantID int64) (*model2.Restaurant, error)
	ListRestaurantMenus(restaurantID int64) ([]model2.RestaurantMenu, error)
	ListRestaurantReviews(restaurantID int64) ([]model2.RestaurantReview, error)
}

func NewRestaurantService(
	restaurantRepository restaurant.RestaurantRepository,
	restaurantMenuRepository restaurant.RestaurantMenuRepository,
	restaurantReviewRepository restaurant.RestaurantReviewRepository,
) RestaurantService {
	return &restaurantService{
		restaurantRepository:       restaurantRepository,
		restaurantMenuRepository:   restaurantMenuRepository,
		restaurantReviewRepository: restaurantReviewRepository,
	}
}

type restaurantService struct {
	restaurantRepository       restaurant.RestaurantRepository
	restaurantMenuRepository   restaurant.RestaurantMenuRepository
	restaurantReviewRepository restaurant.RestaurantReviewRepository
}

func (s *restaurantService) GetRestaurant(restaurantID int64) (*model2.Restaurant, error) {
	restaurant, err := s.restaurantRepository.FindByID(restaurantID)
	if err != nil {
		return nil, err
	}

	return &model2.Restaurant{
		RestaurantID:        restaurant.RestaurantID,
		Name:                restaurant.Name,
		Description:         restaurant.Description,
		PriceRangePerPerson: "restaurant.PriceRangePerPerson",
		Distance:            restaurant.dis,
		BusinessHours:       restaurant.BusinessHours,
		RestaurantImageURLs: restaurant.RestaurantImageUrls,
	}, nil
}

func (s *restaurantService) ListRestaurantMenus(restaurantID int64) ([]model2.RestaurantMenu, error) {
	menus, err := s.restaurantMenuRepository.FindAllByRestaurantID(restaurantID)
	if err != nil {
		return nil, err
	}

	var modelMenus []model2.RestaurantMenu
	for _, menu := range menus {
		modelMenus = append(modelMenus, model2.RestaurantMenu{
			RestaurantMenuID: menu.RestaurantMenuID,
			Name:             menu.Name,
			Price:            menu.Price,
		})
	}

	return modelMenus, nil
}

func (s *restaurantService) ListRestaurantReviews(restaurantID int64) ([]model2.RestaurantReview, error) {
	reviews, err := s.restaurantReviewRepository.FindAllByRestaurantID(restaurantID)
	if err != nil {
		return nil, err
	}

	var modelReviews []model2.RestaurantReview
	for _, review := range reviews {
		modelReviews = append(modelReviews, model2.RestaurantReview{
			RestaurantReviewID: review.RestaurantReviewID,
			UserName:           review.UserName,
			Review:             review.Review,
			Star:               review.Star,
		})
	}

	return modelReviews, nil
}
