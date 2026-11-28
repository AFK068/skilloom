package auth

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

//go:generate go tool mockgen -source=$GOFILE -destination=mock/redis_mock.go -package=mock

type RedisStore interface {
	Get(ctx context.Context, key string) (string, error)
}

type Deps struct {
	fx.In

	Redis RedisStore
	Log   *zap.Logger
}

type Service struct {
	Deps
	log *zap.Logger
}

func New(d Deps) *Service {
	return &Service{
		Deps: d,
		log:  d.Log.Named("auth service"),
	}
}
