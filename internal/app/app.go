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
				lc.Append(fillData(server))
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

func fillData(p *httpsrv.Server) fx.Hook {
	return fx.Hook{
		OnStart: func(_ context.Context) error {
			return p.FillData()
		},
		// OnStop: func(_ context.Context) error {
		// 	return p.Cleanup()
		// },
	}
}
