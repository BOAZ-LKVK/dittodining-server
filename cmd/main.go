package main

import (
	"github.com/BOAZ-LKVK/LKVK-server/pkg/fx/fiberfx"
	"github.com/BOAZ-LKVK/LKVK-server/pkg/fx/gormfx"
	"github.com/BOAZ-LKVK/LKVK-server/pkg/fx/zapfx"
	recommendation_api "github.com/BOAZ-LKVK/LKVK-server/server/api/recommendation"
	"github.com/BOAZ-LKVK/LKVK-server/server/repository/recommendation"
	"github.com/BOAZ-LKVK/LKVK-server/server/repository/restaurant"
	recommendation_service "github.com/BOAZ-LKVK/LKVK-server/server/service/recommendation"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func main() {
	fx.New(
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger { return &fxevent.ZapLogger{Logger: log} }),
		fx.Provide(
			zapfx.NewZapLogger,
		),
		fx.Provide(
			fiberfx.AsAPIController(recommendation_api.NewRecommendationAPIController),
			recommendation.NewRestaurantRecommendationRepository,
			recommendation.NewRestaurantRecommendationRequestRepository,
			recommendation.NewSelectedRestaurantRecommendationRepository,
			recommendation_service.NewRestaurantRecommendationService,
			restaurant.NewRestaurantRepository,
			restaurant.NewRestaurantMenuRepository,
			restaurant.NewRestaurantReviewRepository,
		),
		fiberfx.Module,
		gormfx.Module,
	).Run()
}
