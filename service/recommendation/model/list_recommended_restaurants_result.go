package model

type ListRecommendedRestaurantsResult struct {
	RecommendedRestaurants []*RecommendedRestaurant
	NextCursor             *string
}
