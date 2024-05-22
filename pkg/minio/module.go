package minio

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

const moduleName = "minio"

func NewModule() fx.Option {
	return fx.Module(
		moduleName,

		fx.Provide(NewClient),

		fx.Invoke(func(
			lc fx.Lifecycle,
			log *zap.Logger,
			client *Client,
		) {
			lc.Append(
				fx.Hook{
					OnStart: func(_ context.Context) error {
						client.log.Info("minio client created")
						return nil
					},
					OnStop: func(ctx context.Context) error {
						client.log.Info("minio client stopped")
						return nil
					},
				})
		}),

		fx.Decorate(func(log *zap.Logger) *zap.Logger {
			return log.Named(moduleName)
		}),
	)
}
