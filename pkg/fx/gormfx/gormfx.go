package gormfx

import (
	"context"
	"github.com/BOAZ-LKVK/LKVK-server/domain/recommendation"
	"github.com/BOAZ-LKVK/LKVK-server/domain/restaurant"
	"go.uber.org/fx"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Params struct {
	fx.In

	Config *Config
}

type Result struct {
	fx.Out

	DB *gorm.DB
}

var Module = fx.Module("gorm",
	fx.Provide(parseConfig),
	fx.Provide(New),
	fx.Invoke(connectDB),
)

func New(lc fx.Lifecycle, p Params) (Result, error) {
	// TODO: mysql로 변경
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return Result{}, err
	}

	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				return db.AutoMigrate(
					&restaurant.RestaurantImage{}, &restaurant.Restaurant{}, &restaurant.RestaurantMenu{}, &restaurant.RestaurantReview{},
					&recommendation.RestaurantRecommendation{}, &recommendation.RestaurantRecommendationRequest{},
				)
			},
		},
	)

	return Result{DB: db}, nil
}

func connectDB(db *gorm.DB) {}
