package gormfx

import (
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
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		return Result{}, err
	}

	// TODO: close 등 lifecycle 관리가 필요한지 확인하고 관련 코드 추가

	return Result{DB: db}, nil
}

func connectDB(db *gorm.DB) {}
