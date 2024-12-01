package recommendation

import (
	"context"
	"github.com/BOAZ-LKVK/LKVK-server/server/domain/recommendation"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

var ErrRestaurantRecommendationNotFound = errors.New("restaurant recommendation not found")

type RestaurantRecommendationRepository interface {
	FindAllByRestaurantRecommendationRequestID(ctx context.Context, db *gorm.DB, restaurantRecommendationRequestID int64, cursorRestaurantRecommendationRequestID *int64, limit *int64) ([]*recommendation.RestaurantRecommendation, error)
	FindLastOneByRestaurantRecommendationRequestID(ctx context.Context, db *gorm.DB, restaurantRecommendationRequestID int64) (*recommendation.RestaurantRecommendation, error)
	FindAllByIDs(ctx context.Context, db *gorm.DB, restaurantRecommendationIDs []int64) ([]*recommendation.RestaurantRecommendation, error)
	FindByID(ctx context.Context, db *gorm.DB, restaurantRecommendationID int64) (*recommendation.RestaurantRecommendation, error)
	SaveAll(ctx context.Context, db *gorm.DB, recommendations []*recommendation.RestaurantRecommendation) error
}

func NewRestaurantRecommendationRepository() RestaurantRecommendationRepository {
	return &restaurantRecommendationRepository{}
}

type restaurantRecommendationRepository struct{}

func (r *restaurantRecommendationRepository) FindAllByRestaurantRecommendationRequestID(
	ctx context.Context,
	db *gorm.DB,
	restaurantRecommendationRequestID int64,
	cursorRestaurantRecommendationID *int64,
	limit *int64,
) ([]*recommendation.RestaurantRecommendation, error) {
	var recommendations []*recommendation.RestaurantRecommendation

	whereConditions := make([]interface{}, 0)
	if cursorRestaurantRecommendationID != nil {
		whereConditions = append(whereConditions, "restaurant_recommendation_id > ?", *cursorRestaurantRecommendationID)
	}

	limitQuery := 10
	if limit != nil {
		limitQuery = int(*limit)
	}

	result := db.
		Where(
			recommendation.RestaurantRecommendation{
				RestaurantRecommendationRequestID: restaurantRecommendationRequestID,
			},
			whereConditions...).
		Order("restaurant_recommendation_id ASC").
		Limit(limitQuery).
		Find(&recommendations)
	if result.Error != nil {
		return nil, result.Error
	}

	return recommendations, nil
}

func (r *restaurantRecommendationRepository) FindAllByIDs(
	ctx context.Context,
	db *gorm.DB,
	restaurantRecommendationIDs []int64,
) ([]*recommendation.RestaurantRecommendation, error) {
	var recommendations []*recommendation.RestaurantRecommendation

	result := db.
		Where("restaurant_recommendation_id IN ?", restaurantRecommendationIDs).
		Find(&recommendations)
	if result.Error != nil {
		return nil, result.Error
	}

	return recommendations, nil
}

func (r *restaurantRecommendationRepository) FindByID(
	ctx context.Context,
	db *gorm.DB,
	restaurantRecommendationID int64,
) (*recommendation.RestaurantRecommendation, error) {
	var existingRecommendation recommendation.RestaurantRecommendation

	result := db.
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

func (r *restaurantRecommendationRepository) SaveAll(
	ctx context.Context,
	db *gorm.DB,
	recommendations []*recommendation.RestaurantRecommendation,
) error {
	result := db.Save(recommendations)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *restaurantRecommendationRepository) FindLastOneByRestaurantRecommendationRequestID(ctx context.Context, db *gorm.DB, restaurantRecommendationRequestID int64) (*recommendation.RestaurantRecommendation, error) {
	var existingRecommendation recommendation.RestaurantRecommendation

	result := db.
		Where(&recommendation.RestaurantRecommendation{
			RestaurantRecommendationRequestID: restaurantRecommendationRequestID,
		}).
		Order("restaurant_recommendation_id DESC").
		First(&existingRecommendation)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrRestaurantRecommendationNotFound
		}

		return nil, result.Error
	}

	return &existingRecommendation, nil

}
