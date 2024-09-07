package model

type RestaurantMenu struct {
	RestaurantMenuID int64   `json:"restaurantMenuId"`
	ImageURL         string  `json:"imageUrl"`
	Name             string  `json:"name"`
	Price            int64   `json:"price"`
	Description      *string `json:"description"`
}
