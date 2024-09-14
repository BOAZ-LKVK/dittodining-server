package recommendation

import (
	recommendation_domain "github.com/BOAZ-LKVK/LKVK-server/server/domain/recommendation"
	"gorm.io/gorm"
)

type SelectedRestaurantRecommendationRepository interface {
	SaveAll(selectedRestaurantRecommendations []*recommendation_domain.SelectedRestaurantRecommendation) error
	FindAllByRestaurantRecommendationRequestID(restaurantRecommendationRequestID int64) ([]*recommendation_domain.SelectedRestaurantRecommendation, error)
}

func NewSelectedRestaurantRecommendationRepository(db *gorm.DB) SelectedRestaurantRecommendationRepository {
	return &selectedRestaurantRecommendationRepository{db: db}
}

type selectedRestaurantRecommendationRepository struct {
	db *gorm.DB
}

func (r *selectedRestaurantRecommendationRepository) SaveAll(selectedRestaurantRecommendations []*recommendation_domain.SelectedRestaurantRecommendation) error {
	result := r.db.Save(selectedRestaurantRecommendations)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *selectedRestaurantRecommendationRepository) FindAllByRestaurantRecommendationRequestID(restaurantRecommendationRequestID int64) ([]*recommendation_domain.SelectedRestaurantRecommendation, error) {
	var selectedRestaurantRecommendations []*recommendation_domain.SelectedRestaurantRecommendation
	result := r.db.
		Where(recommendation_domain.SelectedRestaurantRecommendation{
			RestaurantRecommendationRequestID: restaurantRecommendationRequestID,
		}).
		Find(&selectedRestaurantRecommendations)
	if result.Error != nil {
		return nil, result.Error
	}

	return selectedRestaurantRecommendations, nil
}
