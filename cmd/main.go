package main

import (
	recommendation_api "github.com/BOAZ-LKVK/LKVK-server/api/recommendation"
	sample_api "github.com/BOAZ-LKVK/LKVK-server/api/sample"
	"github.com/BOAZ-LKVK/LKVK-server/pkg/fx/fiberfx"
	"github.com/BOAZ-LKVK/LKVK-server/pkg/fx/gormfx"
	"github.com/BOAZ-LKVK/LKVK-server/pkg/fx/zapfx"
	recommendation_repository "github.com/BOAZ-LKVK/LKVK-server/repository/recommendation"
	restaurant_repository "github.com/BOAZ-LKVK/LKVK-server/repository/restaurant"
	sample_repository "github.com/BOAZ-LKVK/LKVK-server/repository/sample"
	recommendation_service "github.com/BOAZ-LKVK/LKVK-server/service/recommendation"
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
			sample_repository.NewSampleRepository,
			fiberfx.AsAPIController(sample_api.NewSampleAPIHandler),
			fiberfx.AsAPIController(recommendation_api.NewRecommendationAPIController),
			recommendation_repository.NewRestaurantRecommendationRepository,
			recommendation_repository.NewRestaurantRecommendationRequestRepository,
			recommendation_service.NewRestaurantRecommendationService,
			restaurant_repository.NewRestaurantRepository,
			restaurant_repository.NewRestaurantMenuRepository,
			restaurant_repository.NewRestaurantReviewRepository,
		),
		fiberfx.Module,
		gormfx.Module,
	).Run()
}
