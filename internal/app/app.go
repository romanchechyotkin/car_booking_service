package app

import (
	"context"

	"github.com/romanchechyotkin/car_booking_service/internal/httpsrv"
	"github.com/romanchechyotkin/car_booking_service/pkg/config"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

type HTTPServer interface {
	RegisterRoutes()
}

func NewService() *fx.App {
	return fx.New(
		fx.Provide(
			zap.NewProduction,
			config.New,
		),

		fx.Options(
			httpsrv.NewModule(),
		),

		fx.Invoke(
			func(server *httpsrv.Server, lc fx.Lifecycle) {
				lc.Append(HttpServerOnStart(server))
			},
		),

		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{
				Logger: log,
			}
		}),
	)
}

func HttpServerOnStart(server HTTPServer) fx.Hook {
	return fx.Hook{
		OnStart: func(ctx context.Context) error {
			server.RegisterRoutes()
			return nil
		},
	}
}
