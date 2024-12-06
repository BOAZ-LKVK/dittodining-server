package gormfx

import (
	"context"
	"fmt"
	"go.uber.org/fx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
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

	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				sqlDB, err := db.DB()
				if err != nil {
					return err
				}

				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()

				if err := sqlDB.PingContext(ctx); err != nil {
					return fmt.Errorf("database connection timeout: %w", err)
				}

				sqlDB.SetConnMaxLifetime(time.Duration(p.Config.ConnMaxLifetimeSeconds) * time.Second)
				sqlDB.SetMaxIdleConns(p.Config.MaxIdleConns)
				sqlDB.SetMaxOpenConns(p.Config.MaxOpenConns)

				return nil
			},
		},
	)

	return Result{DB: db}, nil
}
