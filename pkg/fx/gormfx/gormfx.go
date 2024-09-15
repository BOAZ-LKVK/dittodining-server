package gormfx

import (
	"fmt"
	"go.uber.org/fx"
	"gorm.io/driver/mysql"
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
	var dialector gorm.Dialector
	switch p.Config.DBDriver {
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=UTC",
			p.Config.DBUser, p.Config.DBPassword, p.Config.DBHost, p.Config.DBPort, p.Config.DBName,
		)

		dialector = mysql.Open(dsn)
	default:
		return Result{}, fmt.Errorf("unsupported db driver: %s", p.Config.DBDriver)
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return Result{}, err
	}

	return Result{DB: db}, nil
}

func connectDB(db *gorm.DB) {}
