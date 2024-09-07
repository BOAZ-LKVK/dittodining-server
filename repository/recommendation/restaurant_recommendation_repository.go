package recommendation

import "gorm.io/gorm"

type RestaurantRecommendationRepository struct {
	db *gorm.DB
}

func NewRestaurantRecommendationRepository(db *gorm.DB) *RestaurantRecommendationRepository {
	return &RestaurantRecommendationRepository{db: db}
}
