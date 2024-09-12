package recommendation

import (
	"github.com/BOAZ-LKVK/LKVK-server/domain/recommendation"
	"gorm.io/gorm"
)

type RestaurantRecommendationRepository interface {
	FindAllByRestaurantRecommendationRequestID(restaurantRecommendationRequestID int64, cursorRestaurantRecommendationRequestID *int64, limit *int64) ([]*recommendation.RestaurantRecommendation, error)
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
		Order("restaurant_recommendation_request_id DESC").
		Limit(limitQuery).
		Find(&recommendations)
	if result.Error != nil {
		return nil, result.Error
	}

	return recommendations, nil
}
