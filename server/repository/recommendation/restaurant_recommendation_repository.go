package recommendation

import (
	"github.com/BOAZ-LKVK/LKVK-server/server/domain/recommendation"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

var ErrRestaurantRecommendationNotFound = errors.New("restaurant recommendation not found")

type RestaurantRecommendationRepository interface {
	FindAllByRestaurantRecommendationRequestID(restaurantRecommendationRequestID int64, cursorRestaurantRecommendationRequestID *int64, limit *int64) ([]*recommendation.RestaurantRecommendation, error)
	FindAllByIDs(restaurantRecommendationIDs []int64) ([]*recommendation.RestaurantRecommendation, error)
	FindByID(restaurantRecommendationID int64) (*recommendation.RestaurantRecommendation, error)
	SaveAll(recommendations []*recommendation.RestaurantRecommendation) error
}

func NewRestaurantRecommendationRepository(db *gorm.DB) RestaurantRecommendationRepository {
	return &restaurantRecommendationRepository{db: db}
}

type restaurantRecommendationRepository struct {
	db *gorm.DB
}

func (r *restaurantRecommendationRepository) FindAllByRestaurantRecommendationRequestID(
	restaurantRecommendationRequestID int64,
	cursorRestaurantRecommendationID *int64,
	limit *int64,
) ([]*recommendation.RestaurantRecommendation, error) {
	var recommendations []*recommendation.RestaurantRecommendation

	whereConditions := make([]interface{}, 0)
	if cursorRestaurantRecommendationID != nil {
		whereConditions = append(whereConditions, "restaurant_recommendation_request_id < ?", *cursorRestaurantRecommendationID)
	}

	limitQuery := 10
	if limit != nil {
		limitQuery = int(*limit)
	}

	result := r.db.
		Where(
			recommendation.RestaurantRecommendation{
				RestaurantRecommendationRequestID: restaurantRecommendationRequestID,
			},
			whereConditions...).
		Order("restaurant_recommendation_id DESC").
		Limit(limitQuery).
		Find(&recommendations)
	if result.Error != nil {
		return nil, result.Error
	}

	return recommendations, nil
}

func (r *restaurantRecommendationRepository) FindAllByIDs(restaurantRecommendationIDs []int64) ([]*recommendation.RestaurantRecommendation, error) {
	var recommendations []*recommendation.RestaurantRecommendation

	result := r.db.
		Where("restaurant_recommendation_id IN ?", restaurantRecommendationIDs).
		Find(&recommendations)
	if result.Error != nil {
		return nil, result.Error
	}

	return recommendations, nil
}

func (r *restaurantRecommendationRepository) FindByID(restaurantRecommendationID int64) (*recommendation.RestaurantRecommendation, error) {
	var existingRecommendation recommendation.RestaurantRecommendation

	result := r.db.
		Where(&recommendation.RestaurantRecommendation{
			RestaurantRecommendationID: restaurantRecommendationID,
		}).
		First(&existingRecommendation)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrRestaurantRecommendationNotFound
		}

		return nil, result.Error
	}

	return &existingRecommendation, nil
}

func (r *restaurantRecommendationRepository) SaveAll(recommendations []*recommendation.RestaurantRecommendation) error {
	result := r.db.Save(recommendations)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
