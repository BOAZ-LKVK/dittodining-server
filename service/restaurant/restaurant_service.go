package restaurant

import (
	restaurant_repository "github.com/BOAZ-LKVK/LKVK-server/repository/restaurant"
	"github.com/BOAZ-LKVK/LKVK-server/service/restaurant/model"
)

type RestaurantService interface {
	GetRestaurant(restaurantID int64) (*model.Restaurant, error)
	ListRestaurantMenus(restaurantID int64) ([]model.RestaurantMenu, error)
	ListRestaurantReviews(restaurantID int64) ([]model.RestaurantReview, error)
}

func NewRestaurantService(
	restaurantRepository restaurant_repository.RestaurantRepository,
	restaurantMenuRepository restaurant_repository.RestaurantMenuRepository,
	restaurantReviewRepository restaurant_repository.RestaurantReviewRepository,
) RestaurantService {
	return &restaurantService{
		restaurantRepository:       restaurantRepository,
		restaurantMenuRepository:   restaurantMenuRepository,
		restaurantReviewRepository: restaurantReviewRepository,
	}
}

type restaurantService struct {
	restaurantRepository       restaurant_repository.RestaurantRepository
	restaurantMenuRepository   restaurant_repository.RestaurantMenuRepository
	restaurantReviewRepository restaurant_repository.RestaurantReviewRepository
}

func (s *restaurantService) GetRestaurant(restaurantID int64) (*model.Restaurant, error) {
	restaurant, err := s.restaurantRepository.FindByID(restaurantID)
	if err != nil {
		return nil, err
	}

	return &model.Restaurant{
		RestaurantID:        restaurant.RestaurantID,
		Name:                restaurant.Name,
		Description:         restaurant.Description,
		PriceRangePerPerson: "restaurant.PriceRangePerPerson",
		Distance:            restaurant.dis,
		BusinessHours:       restaurant.BusinessHours,
		RestaurantImageURLs: restaurant.RestaurantImageUrls,
	}, nil
}

func (s *restaurantService) ListRestaurantMenus(restaurantID int64) ([]model.RestaurantMenu, error) {
	menus, err := s.restaurantMenuRepository.FindAllByRestaurantID(restaurantID)
	if err != nil {
		return nil, err
	}

	var modelMenus []model.RestaurantMenu
	for _, menu := range menus {
		modelMenus = append(modelMenus, model.RestaurantMenu{
			RestaurantMenuID: menu.RestaurantMenuID,
			Name:             menu.Name,
			Price:            menu.Price,
		})
	}

	return modelMenus, nil
}

func (s *restaurantService) ListRestaurantReviews(restaurantID int64) ([]model.RestaurantReview, error) {
	reviews, err := s.restaurantReviewRepository.FindAllByRestaurantID(restaurantID)
	if err != nil {
		return nil, err
	}

	var modelReviews []model.RestaurantReview
	for _, review := range reviews {
		modelReviews = append(modelReviews, model.RestaurantReview{
			RestaurantReviewID: review.RestaurantReviewID,
			UserName:           review.UserName,
			Review:             review.Review,
			Star:               review.Star,
		})
	}

	return modelReviews, nil
}
