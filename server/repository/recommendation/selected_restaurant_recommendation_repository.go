package recommendation

import (
	"context"
	recommendation_domain "github.com/BOAZ-LKVK/LKVK-server/server/domain/recommendation"
	"gorm.io/gorm"
)

type SelectedRestaurantRecommendationRepository interface {
	SaveAll(ctx context.Context, db *gorm.DB, selectedRestaurantRecommendations []*recommendation_domain.SelectedRestaurantRecommendation) error
	FindAllByRestaurantRecommendationRequestID(ctx context.Context, db *gorm.DB, restaurantRecommendationRequestID int64) ([]*recommendation_domain.SelectedRestaurantRecommendation, error)
}

func NewSelectedRestaurantRecommendationRepository() SelectedRestaurantRecommendationRepository {
	return &selectedRestaurantRecommendationRepository{}
}

type selectedRestaurantRecommendationRepository struct {
	db *gorm.DB
}

func (r *selectedRestaurantRecommendationRepository) SaveAll(ctx context.Context, db *gorm.DB, selectedRestaurantRecommendations []*recommendation_domain.SelectedRestaurantRecommendation) error {
	result := db.Save(selectedRestaurantRecommendations)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *selectedRestaurantRecommendationRepository) FindAllByRestaurantRecommendationRequestID(ctx context.Context, db *gorm.DB, restaurantRecommendationRequestID int64) ([]*recommendation_domain.SelectedRestaurantRecommendation, error) {
	var selectedRestaurantRecommendations []*recommendation_domain.SelectedRestaurantRecommendation
	result := db.
		Where(recommendation_domain.SelectedRestaurantRecommendation{
			RestaurantRecommendationRequestID: restaurantRecommendationRequestID,
		}).
		Find(&selectedRestaurantRecommendations)
	if result.Error != nil {
		return nil, result.Error
	}

	return selectedRestaurantRecommendations, nil
}
