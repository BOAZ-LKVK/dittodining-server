package gormfx

import (
	"context"
	recommendation2 "github.com/BOAZ-LKVK/LKVK-server/server/domain/recommendation"
	restaurant2 "github.com/BOAZ-LKVK/LKVK-server/server/domain/restaurant"
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
					&restaurant2.RestaurantImage{}, &restaurant2.Restaurant{}, &restaurant2.RestaurantMenu{}, &restaurant2.RestaurantReview{},
					&recommendation2.RestaurantRecommendation{}, &recommendation2.RestaurantRecommendationRequest{},
				)
			},
		},
	)

	return Result{DB: db}, nil
}

func connectDB(db *gorm.DB) {}
